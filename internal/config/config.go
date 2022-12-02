package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

const (
	defaultConfigPath string = "configuration"
	EnvPrefix         string = "MPC"
)

func LoadConfig(configPath *string) (*Config, error) {
	log.Debug().Msg("Loading config...")

	v, err := readConfigYML(configPath)
	if err != nil {
		return nil, err
	}
	serverCfg, err := LoadServerConfig()
	if err != nil {
		return nil, err
	}
	if serverCfg == nil {
		return nil, fmt.Errorf("cannot load server config")
	}

	var cfg Config
	if err = v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("unmarshal cfg: %w", err)
	}
	var postgresConfig PostgresConfig
	if cfg.DB.Profile != "" {
		log.Debug().Msgf("Using DB profile: %s", cfg.DB.Profile)
		if err := v.UnmarshalKey(fmt.Sprintf("db.%s", cfg.DB.Profile), &postgresConfig); err != nil {
			panic(err)
		}
		cfg.DB.Postgresql = postgresConfig
	}

	cfg.Server = *serverCfg
	return &cfg, nil
}

func LoadServerConfig() (*ServerConfig, error) {
	log.Debug().Msg("Loading server config...")

	v, err := readConfigJSON()
	if err != nil {
		return nil, err
	}

	var cfg ServerConfig
	if err = v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("unmarshal cfg: %w", err)
	}

	return &cfg, nil
}

func readConfigYML(configPath *string) (*viper.Viper, error) {
	v := viper.New()
	if configPath != nil {
		v.SetConfigFile(*configPath)
	} else {
		v.SetConfigFile(defaultConfigPath)
	}
	v.AddConfigPath("./")

	v.SetEnvPrefix(EnvPrefix)
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("read cfg: %w", err)
	}
	return v, nil
}

func readConfigJSON() (*viper.Viper, error) {
	v := viper.New()
	configFileName := os.Getenv("CONFIG_FILENAME")
	if configFileName == "" {
		configFileName = "local"
	}
	v.SetConfigName(configFileName)
	v.AddConfigPath("config")
	v.SetConfigType("json")

	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("read cfg: %w", err)
	}
	return v, nil
}
