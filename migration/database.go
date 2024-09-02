package migration

import (
	"database/sql"
)

func ConnectToDatabase(config *Config) (*sql.DB, error) {
	return sql.Open("oracle", config.Database.BuildUrl())
}
