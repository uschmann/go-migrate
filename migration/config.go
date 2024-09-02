package migration

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	go_ora "github.com/sijms/go-ora/v2"
)

type DatabaseConfig struct {
	User     string
	Password string
	Host     string
	Port     int
	Service  string
}

type Config struct {
	MigrationLogTable string
	Database          DatabaseConfig
}

func MakeConfig() *Config {
	godotenv.Load()

	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	host := os.Getenv("DB_HOST")
	port, _ := strconv.Atoi(os.Getenv("DB_PORT"))
	service := os.Getenv("DB_SERVICE")

	migrationLogTable := os.Getenv("DB_MIGRATION_LOG_TABLE")
	if migrationLogTable == "" {
		migrationLogTable = "MIGRATION_LOGS"
	}

	return &Config{
		MigrationLogTable: migrationLogTable,
		Database: DatabaseConfig{
			User:     user,
			Password: password,
			Host:     host,
			Port:     port,
			Service:  service,
		},
	}
}

func (d *DatabaseConfig) BuildUrl() string {
	return go_ora.BuildUrl(d.Host, d.Port, d.Service, d.User, d.Password, nil)
}

func (d *DatabaseConfig) BuildSqlplusConnectionString() string {
	return fmt.Sprintf("%s/%s@%s:%d/%s", d.User, d.Password, d.Host, d.Port, d.Service)
}
