package utils

// read the configurations

import (
	"time"

	"github.com/spf13/viper"
)

// contains the configs for the database in a yaml file:
type DatabaseConfig struct {
	Type     string `mapstructure:"type"`
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Database string `mapstructure:"database"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
}

// the app config contains related configurations for the app
type AppConfig struct {
	Port                   int           `mapstructure:"port"`
	JwtSecret              string        `mapstructure:"jwtSecret"`
	TokenAccessExpiration  time.Duration `mapstructure:"tokenAccessExpiration"`
	TokenRefreshExpiration time.Duration `mapstructure:"tokenRefreshExpiration"`
}

// a helper function that takes the name, path and struct
func readConfigFile(path, name string) error {
	viper.SetConfigName(name)
	viper.AddConfigPath(path)
	return viper.ReadInConfig()
}

func unmarshalConfig(config interface{}) error {
	return viper.Unmarshal(config)
}

// read and parse the configurations from the yaml file from the path given
func ConfigStore(path, database, app string) (DatabaseConfig, AppConfig, error) {
	var dbConfig DatabaseConfig
	var appConfig AppConfig

	if err := readConfigFile(path, database); err != nil {
		return dbConfig, appConfig, err
	}
	if err := unmarshalConfig(&dbConfig); err != nil {
		return dbConfig, appConfig, err
	}

	if err := readConfigFile(path, app); err != nil {
		return dbConfig, appConfig, err
	}
	if err := unmarshalConfig(&appConfig); err != nil {
		return dbConfig, appConfig, err
	}

	return dbConfig, appConfig, nil
}
