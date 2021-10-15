package postgresql

import (
	"fmt"
	"log"
	"time"


	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // For postgres dialect
	"github.com/logrusorgru/aurora"

	"go.elastic.co/apm/module/apmsql"
	_ "go.elastic.co/apm/module/apmsql/pq" // For apm dialect
)

var (
	sqlxClient *sqlx.DB
)

// Connect to postgresql database
func Connect(host, user, password, dbname, port, sslmode string) error {
	// Connect string
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=UTC",
		host, user, password, dbname, port, sslmode,
	)

	// TODO: write case for SSL mode

	// db, err := sqlx.Connect("postgres", dsn)
	// if err != nil {
	// 	log.Fatalln(err)
	// }

	apmDB, err := apmsql.Open("postgres", dsn)
	if err != nil {
		log.Fatalln(err)
	}

	db := sqlx.NewDb(apmDB, "postgres")

	fmt.Println(aurora.Green("*** CONNECTED TO POSTGRESQL - SQLX: " + dsn))

	// Config connection pool
	sqlDB := db.DB
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetMaxIdleConns(20)
	sqlDB.SetConnMaxLifetime(time.Minute * 5)

	// Assign client
	sqlxClient = db

	return nil
}

// GetSqlxInstance ...
func GetSqlxInstance() *sqlx.DB {
	return sqlxClient
}

