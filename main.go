package main

import (
	"github.com/uschmann/dbmigrate/migration"
)

func main() {
	service := migration.MakeMigrationService("./sql")
	service.GenerateMigration("foo")
}
