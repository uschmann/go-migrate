package migration

import (
	"os"
	"strconv"
)

type DatabaseConfig struct {
	User     string
	Password string
	Host     string
	Port     int
	Service  string
}

type Config struct {
	Database DatabaseConfig
}

func MakeConfig() *Config {
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	host := os.Getenv("DB_HOST")
	port, _ := strconv.Atoi(os.Getenv("DB_PORT"))
	service := os.Getenv("DB_SERVICE")

	return &Config{
		Database: DatabaseConfig{
			User:     user,
			Password: password,
			Host:     host,
			Port:     port,
			Service:  service,
		},
	}
}
