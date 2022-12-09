package config

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	dbClient "github.com/lugondev/mpc-tss-lib/db/client"
	dbGateway "github.com/lugondev/mpc-tss-lib/db/gateway"
	"github.com/rs/zerolog"
	"time"
)

func NewClientDB(postgresConfig PostgresConfig, logger *zerolog.Logger) (*dbClient.SQLStore, error) {
	openDB, err := newDB(postgresConfig)
	if err != nil {
		return nil, err
	}

	return dbClient.NewStore(openDB, logger), nil
}

func NewGatewayDB(postgresConfig PostgresConfig, logger *zerolog.Logger) (*dbGateway.SQLStore, error) {
	openDB, err := newDB(postgresConfig)
	if err != nil {
		return nil, err
	}

	return dbGateway.NewStore(openDB, logger), nil
}

func newDB(postgresConfig PostgresConfig) (*sql.DB, error) {
	dsn := getDSN(postgresConfig)
	openDB, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	openDB.SetConnMaxIdleTime(time.Duration(postgresConfig.ConnectionMaxIdleTime) * time.Minute)
	openDB.SetConnMaxLifetime(time.Duration(postgresConfig.ConnectionMaxLifetime) * time.Minute)
	openDB.SetMaxIdleConns(postgresConfig.MaxIdleConnections)
	openDB.SetMaxOpenConns(postgresConfig.MaxOpenConnections)

	return openDB, nil
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
