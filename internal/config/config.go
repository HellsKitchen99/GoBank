package config

import (
	"time"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

func LoadEnv() error {
	if err := godotenv.Load(); err != nil {
		return err
	}
	return nil
}

type DataBaseConfig struct {
	User     string `env:"DB_USER,required"`
	Password string `env:"DB_PASSWORD,required"`
	Host     string `env:"DB_HOST,required"`
	Port     string `env:"DB_PORT,required"`
	DbName   string `env:"DB_NAME,required"`
}

type JwtConfig struct {
	Key      string        `env:"JWT_KEY,required"`
	Duration time.Duration `env:"JWT_KEY_LIVING_PERIOD,required"`
}

type Config struct {
	DB  DataBaseConfig
	Jwt JwtConfig
}

func Load() (*Config, error) {
	var config Config
	_ = LoadEnv()
	if err := env.Parse(&config.DB); err != nil {
		return &config, err
	}
	if err := env.Parse(&config.Jwt); err != nil {
		return &config, err
	}
	return &config, nil
}
