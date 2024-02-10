package utils

// read the configurations

import (
	"time"

	"github.com/spf13/viper"
)

// contains the configs for the database in a yaml file:
type DatabaseConfig struct {
	Type     string `mapstructure:"type"`     // the type of the database
	Host     string `mapstructure:"host"`     // the host for the database
	Port     int    `mapstructure:"port"`     // the port for the database
	Database string `mapstructure:"database"` // the name of the database
	Username string `mapstructure:"username"` // the username for the database
	Password string `mapstructure:"password"` // the password for the database
}

// the app config contains related configurations for the app
type AppConfig struct {
	Port                   int           `mapstructure:"aport" env:"APORT"`
	Env                    string        `mapstructure:"env" env:"ENV"`
	JwtSecret              string        `mapstructure:"jwtSecret"`
	TokenAccessExpiration  time.Duration `mapstructure:"tokenAccessExpiration"`
	TokenRefreshExpiration time.Duration `mapstructure:"tokenRefreshExpiration"`
}

// a helper function that takes the name, path and struct
func readConfigFile(path, name string) error {
	viper.SetConfigName(name)
	viper.AddConfigPath(path)

	// override from environment variables
	viper.AutomaticEnv()

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
