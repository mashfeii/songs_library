package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	ServingPort int    `mapstructure:"SERVING_PORT"`
	APIEndpoint string `mapstructure:"API_ENDPOINT"`
	Host        string `mapstructure:"DB_HOST"`
	Port        int    `mapstructure:"DB_PORT"`
	Username    string `mapstructure:"DB_USER"`
	Password    string `mapstructure:"DB_PASSWORD"`
	Name        string `mapstructure:"DB_NAME"`
}

func (d *Config) ToDSN() string {
	return fmt.Sprintf("postgresql://%s:%s@%s:%d/%s?sslmode=disable",
		d.Username,
		d.Password,
		d.Host,
		d.Port,
		d.Name,
	)
}

func NewConfigFromFile(name string) (cfg *Config, err error) {
	v := viper.New()

	v.SetConfigType("env")
	v.AddConfigPath(".")
	v.SetConfigName(name)

	v.SetDefault("SERVING_PORT", 8080)
	v.SetDefault("API_ENDPOINT", "http://localhost:8081/info")
	v.SetDefault("DB_HOST", "localhost")
	v.SetDefault("DB_PORT", 5432)
	v.SetDefault("DB_USER", "postgres")
	v.SetDefault("DB_PASSWORD", "password")
	v.SetDefault("DB_NAME", "mydb")

	v.AutomaticEnv()

	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("reading config: %w", err)
		}
	}

	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("unmarshalling config: %w", err)
	}

	return cfg, nil
}
