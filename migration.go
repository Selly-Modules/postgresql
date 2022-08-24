package postgresql

import (
	"database/sql"
	"fmt"
	"github.com/Selly-Modules/logger"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"os"
)

func runMigration(db *sql.DB, server string) {
	// init migrate data
	driver, _ := postgres.WithInstance(db, &postgres.Config{})
	m, err := migrate.NewWithDatabaseInstance(
		fmt.Sprintf("file:///%s", getMigrationDirectoryPath(server)),
		"postgres",
		driver,
	)

	// up
	if err == nil {
		if err = m.Up(); err != nil && err.Error() != "no change" {
			logger.Error("run migration", logger.LogData{
				Source:  "runMigration",
				Message: err.Error(),
				Data:    nil,
			})
			fmt.Println("run schema migration error:", err.Error())
		}
		fmt.Printf("⚡️[postgres]: done migration \n")
	}
}

func getMigrationDirectoryPath(server string) string {
	migrationDir := fmt.Sprintf("/external/postgresql/%s/migrations/sql", server)

	dirname, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	// Check path existed
	path := dirname + migrationDir
	if _, err = os.Stat(path); os.IsNotExist(err) {
		// Create if not existed
		err = os.Mkdir(path, 0755)
		if err != nil {
			panic(err)
		}
	}

	return path
}
