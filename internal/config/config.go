package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

const defaultConfigPath string = "configuration"

func LoadConfig(configPath *string) (*Config, *viper.Viper, error) {
	log.Debug().Msg("Loading config...")

	v, err := readConfigYML(configPath)
	if err != nil {
		return nil, nil, err
	}
	serverCfg, err := LoadServerConfig()
	if err != nil {
		return nil, nil, err
	}
	if serverCfg == nil {
		return nil, nil, fmt.Errorf("cannot load server config")
	}

	var cfg Config
	if err = v.Unmarshal(&cfg); err != nil {
		return nil, nil, fmt.Errorf("unmarshal cfg: %w", err)
	}

	cfg.Server = *serverCfg
	return &cfg, v, nil
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

	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

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
