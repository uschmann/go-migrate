package migration

import (
	"database/sql"

	go_ora "github.com/sijms/go-ora/v2"
)

func ConnectToDatabase() (*sql.DB, error) {
	connStr := go_ora.BuildUrl("localhost", 1522, "FREE", "auschmann", "secret", nil)
	return sql.Open("oracle", connStr)
}
