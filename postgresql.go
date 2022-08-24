package postgresql

import (
	"database/sql"
	"fmt"
	"github.com/Selly-Modules/logger"
	"time"

	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

type Config struct {
	Host               string
	Port               int
	User               string
	Password           string
	DBName             string
	SSLMode            string
	IsDebug            bool
	MaxOpenConnections int
	MaxIdleConnections int
	ConnectionLifetime time.Duration
}

// Connect ...
func Connect(cfg Config, server string) *sql.DB {
	uri := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DBName, cfg.SSLMode)

	// connect
	db, err := sql.Open("pgx", uri)
	if err != nil {
		panic(err)
	}

	// ping
	if err = db.Ping(); err != nil {
		logger.Error("pgx ping", logger.LogData{
			Source:  "Connect",
			Message: err.Error(),
			Data:    cfg,
		})
		panic(err)
	}

	// config
	if cfg.MaxOpenConnections == 0 {
		cfg.MaxOpenConnections = 25
	}
	if cfg.MaxIdleConnections == 0 {
		cfg.MaxIdleConnections = 25
	}
	if cfg.ConnectionLifetime == 0 {
		cfg.ConnectionLifetime = 5 * time.Minute
	}
	db.SetMaxOpenConns(cfg.MaxOpenConnections)
	db.SetMaxIdleConns(cfg.MaxIdleConnections)
	db.SetConnMaxLifetime(cfg.ConnectionLifetime)

	// run migration
	runMigration(db, server)

	// debug mode
	boil.DebugMode = cfg.IsDebug

	fmt.Printf("⚡️[postgres]: connected to %s:%d \n", cfg.Host, cfg.Port)

	return db
}
