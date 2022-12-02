package config

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/lugondev/mpc-tss-lib/db"
	"time"
)

func NewDB(postgresConfig PostgresConfig) (*db.SQLStore, error) {
	dsn := getDSN(postgresConfig)
	openDB, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	openDB.SetConnMaxIdleTime(time.Duration(postgresConfig.ConnectionMaxIdleTime) * time.Minute)
	openDB.SetConnMaxLifetime(time.Duration(postgresConfig.ConnectionMaxLifetime) * time.Minute)
	openDB.SetMaxIdleConns(postgresConfig.MaxIdleConnections)
	openDB.SetMaxOpenConns(postgresConfig.MaxOpenConnections)

	return db.NewStore(openDB), nil
}

func getDSN(configuration PostgresConfig) string {
	return fmt.Sprintf(
		"host=%s port=%d dbname=%s user=%s password=%s sslmode=%s TimeZone=%s",
		configuration.Host,
		configuration.Port,
		configuration.DBName,
		configuration.Username,
		configuration.Password,
		configuration.SSLMode,
		configuration.TimeZone,
	)
}
