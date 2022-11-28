package config

import (
	"fmt"
	"strings"

	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

const defaultConfigPath string = "configuration.yml"

func LoadConfig(configPath *string) (*Config, *viper.Viper, error) {
	log.Debug().Msg("Loading config...")

	v, err := readConfig(configPath)
	if err != nil {
		return nil, nil, err
	}

	var cfg Config
	if err = v.Unmarshal(&cfg); err != nil {
		return nil, nil, fmt.Errorf("unmarshal cfg: %w", err)
	}

	return &cfg, v, nil
}

func readConfig(configPath *string) (*viper.Viper, error) {
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
