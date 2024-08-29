package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/uschmann/go-migrate/migration"
)

// https://github.com/lucasjellema/go-oracle-database/blob/main/with-oracleinstant-client.go

func main() {
	connection, err := migration.ConnectToDatabase()

	if err != nil {
		panic(err)
	}

	migrationLogRepository := migration.NewMigrationLogRepository(connection)
	migrationService := migration.MakeMigrationService("./sql", migrationLogRepository)

	migrationService.Up()
}

func main2() {
	if len(os.Args) < 2 {
		fmt.Println("Expected at least one command")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "make":
		makeCmd := flag.NewFlagSet("make", flag.ExitOnError)
		makeName := makeCmd.String("n", "", "The migration name")
		makeTable := makeCmd.String("t", "", "Specify an optional table name")

		makeCmd.Parse(os.Args[2:])

		if *makeName == "" {
			printMakeHelp()
			os.Exit(1)
		}

		make(*makeName, *makeTable)
	case "status":
		fmt.Println("TBD: status")
	}

}

func make(name string, table string) {
	dir := "./sql"

	if table == "" {
		migration.GenerateMigration(dir, name)
	} else {
		migration.GenerateMigrationWithTemplate(dir, name, table)
	}

}

func printMakeHelp() {
	fmt.Println("Usage: make -n migrationName -t tableName")
}
