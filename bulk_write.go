package postgresql

import (
	"context"
	"fmt"
	"reflect"
	"strings"

	"github.com/volatiletech/sqlboiler/v4/boil"

	"github.com/Selly-Modules/logger"
)

const (
	indexSign = "$%d"
	decodeTag = "boil"
)

type bulkInsertItem map[string]interface{}

type BulkInsertPayload struct {
	TableName string
	Data      interface{}
	Columns   []string
}

func BulkInsert(ctx context.Context, db boil.ContextExecutor, payload BulkInsertPayload) error {
	// convert data to map
	insertRows, err := toMap(payload.Data)
	if err != nil {
		logger.Error("convert data to map", logger.LogData{
			Source:  "module.postgresql.bulk_insert.BulkInsert",
			Message: err.Error(),
			Data:    payload.Data,
		})
		return err
	}

	var (
		numOfColumns  = len(payload.Columns)
		incSigns      = make([]string, 0)
		insertValues  = make([]interface{}, 0)
		listOfSigns   = make([]string, numOfColumns)
		listOfColumns = make([]string, numOfColumns)
	)

	// prepare array of dollar signs
	for i := range listOfSigns {
		listOfSigns[i] = indexSign
		listOfColumns[i] = payload.Columns[i]
	}

	// prepare sign for every column
	signs := strings.Join(listOfSigns, ",")
	insertColumns := strings.Join(listOfColumns, ",")

	for index, r := range insertRows {
		currentIncSigns := getIncSignValues(index, numOfColumns)
		incSigns = append(incSigns, fmt.Sprintf("("+signs+")",
			currentIncSigns...,
		))

		currentInsertValues := getInsertValues(r, payload.Columns, numOfColumns)
		insertValues = append(insertValues, currentInsertValues...)
	}

	// exec
	stm := getSQLStatement(payload.TableName, insertColumns, incSigns)
	_, err = db.ExecContext(ctx, stm, insertValues...)
	if err != nil {
		logger.Error("insert to db", logger.LogData{
			Source:  "module.postgresql.bulk_insert.BulkInsert",
			Message: err.Error(),
			Data:    logger.Map{"statement": stm},
		})
		return err
	}

	return nil
}

func getSQLStatement(tableName, insertColumns string, incSigns []string) string {
	return fmt.Sprintf(`
			INSERT INTO %s (
				%s
			)
			VALUES %s
	`, tableName, insertColumns, strings.Join(incSigns, ","))
}

func getIncSignValues(index, numOfColumns int) (result []interface{}) {
	for i := 1; i <= numOfColumns; i++ {
		result = append(result, index*numOfColumns+i)
	}
	return
}

func getInsertValues(row bulkInsertItem, columns []string, numOfColumns int) (result []interface{}) {
	for i := 0; i < numOfColumns; i++ {
		result = append(result, row[columns[i]])
	}
	return
}

func toMap(input interface{}) ([]bulkInsertItem, error) {
	result := make([]bulkInsertItem, 0)

	v := reflect.ValueOf(input)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	// only accept slices
	if v.Kind() != reflect.Slice {
		err := fmt.Errorf("toMap only accepts slices; got %T", v)
		logger.Error("invalid type", logger.LogData{
			Source:  "external.postgresql.bulk_insert.toMap",
			Message: err.Error(),
			Data:    nil,
		})
		return nil, err
	}

	// loop and assign data to result
	for x := 0; x < v.Len(); x++ {
		item := bulkInsertItem{}
		typ := v.Index(x).Type()

		// loop each field
		for i := 0; i < v.Index(x).NumField(); i++ {
			fi := typ.Field(i)
			t := fi.Tag.Get(decodeTag)
			if t != "" {
				// set key of map to value in struct field
				item[t] = v.Index(x).Field(i).Interface()
			}
		}
		result = append(result, item)
	}

	return result, nil
}
