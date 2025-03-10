package config

import (
	"bytes"
	_ "embed"
	"strings"

	"github.com/joho/godotenv"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

type Config struct {
	App      App
	Postgres Postgres
}

type App struct {
	Port string `mapstructure:"port"`
}

type Postgres struct {
	ConnStr string `mapstructure:"connstr"`
}

//go:embed config.yaml
var config []byte

func New(configPath string) (*Config, error) {
	//loading env variables to ovveride
	if err := godotenv.Load(); err != nil {
		return nil, errors.Wrap(err, "config.New")
	}

	const configType = "yaml"
	viper.SetConfigFile(configPath)
	viper.SetConfigType(configType)

	if err := viper.ReadConfig(bytes.NewReader(config)); err != nil {
		return nil, errors.Wrap(err, "config.New")
	}

	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, errors.Wrap(err, "config.New")
	}

	return &config, nil
}
