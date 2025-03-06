package config

import (
	"bytes"
	_ "embed"
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	App      App
	Postgres Postgres
}

type App struct {
	Port string `yaml:"port"`
}

type Postgres struct {
	ConnStr string `env:"connstr"`
}

//go:embed config.yaml
var config []byte

func New() (*Config, error) {
	const configPath = "internal/config/config.yaml"
	const configType = "yaml"
	viper.SetConfigFile(configPath)
	viper.SetConfigType(configType)

	if err := viper.ReadConfig(bytes.NewReader(config)); err != nil {
		return nil, fmt.Errorf("config.New: %w", err)
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("config.New: %w", err)
	}
	return &config, nil
}
