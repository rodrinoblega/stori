package config

import (
	"fmt"
	"github.com/spf13/viper"
)

type (
	Config struct {
		Env string
		EmailConfig
		PgConfig
	}

	EmailConfig struct {
		EmailHost     string
		EmailPort     string
		EmailUsername string
		EmailPassword string
	}

	PgConfig struct {
		PgUser       string
		PgPassword   string
		PgDatabase   string
		PostgresPort string
		PostgresHost string
	}
)

func Load(env string) *Config {
	if env == "" {
		env = "local"
	}

	viper.AutomaticEnv()
	viper.SetConfigName(fmt.Sprintf("config_%s", env))
	viper.SetConfigType("json")
	viper.AddConfigPath("./../config") // local
	viper.AddConfigPath("./config")    // local
	viper.AddConfigPath("/app/config") // inside container

	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}

	return &Config{
		Env: env,
		EmailConfig: EmailConfig{
			EmailHost:     viper.GetString("EMAIL_HOST"),
			EmailPort:     viper.GetString("EMAIL_PORT"),
			EmailUsername: viper.GetString("EMAIL_USERNAME"),
			EmailPassword: viper.GetString("EMAIL_PASSWORD"),
		},
		PgConfig: PgConfig{
			PgUser:       viper.GetString("POSTGRES_USER"),
			PgPassword:   viper.GetString("POSTGRES_PASSWORD"),
			PgDatabase:   viper.GetString("POSTGRES_DB"),
			PostgresPort: viper.GetString("POSTGRES_PORT"),
			PostgresHost: viper.GetString("POSTGRES_HOST"),
		},
	}
}
