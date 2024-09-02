package migration

import (
	"fmt"
	"log"
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
	Database DatabaseConfig
}

func readDotEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func MakeConfig() *Config {
	readDotEnv()

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

func (d *DatabaseConfig) BuildUrl() string {
	return go_ora.BuildUrl(d.Host, d.Port, d.Service, d.User, d.Password, nil)
}

func (d *DatabaseConfig) BuildSqlplusConnectionString() string {
	return fmt.Sprintf("%s/%s@%s:%d/%s", d.User, d.Password, d.Host, d.Port, d.Service)
	//return "auschmann/secret@localhost:1522/FREE"
}
