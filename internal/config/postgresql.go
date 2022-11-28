package config

import (
	"database/sql"
	"fmt"
	"github.com/ambrosus/ambrosus-bridge/relay/db"
	_ "github.com/lib/pq"
	"github.com/spf13/viper"
	"time"
)

func NewDB(v *viper.Viper) (*db.SQLStore, error) {
	var config DBConfig
	var postgresConfig PostgresConfig
	if err := v.UnmarshalKey("db", &config); err != nil {
		return nil, err
	}
	var dsn string
	fmt.Println("Using DB profile:", config.Profile)
	if config.Profile != "" {
		if err := v.UnmarshalKey(fmt.Sprintf("db.%s", config.Profile), &postgresConfig); err != nil {
			panic(err)
		}
		dsn = getDSN(postgresConfig)
	} else {
		dsn = getDSN(config.Postgresql)
	}

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
