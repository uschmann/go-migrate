package main

import (
	"github.com/uschmann/go-migrate/migration"
)

func main() {
	service := migration.MakeMigrationService("./sql")
	service.GenerateMigration("foo")
}
