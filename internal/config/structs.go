package config

import "time"

type PostgresConfig struct {
	Host     string
	Port     int
	DBName   string
	Username string
	Password string
	SSLMode  string
	TimeZone string

	ConnectionMaxIdleTime int
	ConnectionMaxLifetime int
	MaxIdleConnections    int
	MaxOpenConnections    int

	Logger struct {
		LogLevel string
	}
}

type OIDCConfig struct {
	Issuer     string
	CacheTTL   time.Duration
	Audience   []string
	ClaimsPath string
}

type AuthConfig struct {
	Domain       string
	ClientId     string
	ClientSecret string
	OIDC         *OIDCConfig
}

type DBConfig struct {
	Profile    string
	Postgresql PostgresConfig
}

type ServerConfig struct {
	Port    int64    `json:"port"`
	Clients []string `json:"clients"`
}

type GrpcConfig struct {
	Port int64
}

type Config struct {
	DB         DBConfig     `yml:"db"`
	AuthConfig AuthConfig   `yml:"auth"`
	Server     ServerConfig `yml:"server"`
	Grpc       GrpcConfig   `yml:"grpc"`
}
