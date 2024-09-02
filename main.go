package main

import (
	"fmt"
	"log"
	"os"

	"github.com/urfave/cli/v2"
	"github.com/uschmann/go-migrate/migration"
)

const DEFAULT_MIGRATION_DIR string = "./sql"

func main() {
	var directory string

	app := &cli.App{
		Name:  "dbmigrate",
		Usage: "Create and execute migrations for oracle db",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "dir",
				Value:       DEFAULT_MIGRATION_DIR,
				Aliases:     []string{"d"},
				Usage:       "Path to the folder that contains the migrations",
				Destination: &directory,
			},
		},
		Commands: []*cli.Command{
			{
				Name:  "status",
				Usage: "List pending and executed migrations",
				Action: func(ctx *cli.Context) error {
					config := migration.MakeConfig()
					connection, err := migration.ConnectToDatabase(config)
					if err != nil {
						panic(err)
					}

					migrationLogRepository := migration.NewMigrationLogRepository(connection, config)
					migrationService := migration.MakeMigrationService(directory, config, migrationLogRepository)

					migrationLogRepository.CreateMigrationLogsTable()

					migrationStatus := migrationService.GetMigrationStatus()

					fmt.Println("\tName\t\t\t\t\t\tExecuted?\tBatch")
					for _, status := range migrationStatus {
						isExecuted := "No"
						if status.IsExecuted {
							isExecuted = "Yes"
						}

						fmt.Println(status.Index, "\t"+status.Migration.Name+"\t\t", isExecuted, "\t\t", status.Batch)
					}
					return nil
				},
			},
			{
				Name:  "migrate",
				Usage: "Execute all pending migrations",
				Action: func(ctx *cli.Context) error {
					config := migration.MakeConfig()
					connection, err := migration.ConnectToDatabase(config)
					if err != nil {
						panic(err)
					}

					migrationLogRepository := migration.NewMigrationLogRepository(connection, config)
					migrationService := migration.MakeMigrationService(directory, config, migrationLogRepository)

					migrationLogRepository.CreateMigrationLogsTable()
					migrationService.Up()

					return nil
				},
			},
			{
				Name:  "rollback",
				Usage: "Rollback the last batch of migrations",
				Action: func(ctx *cli.Context) error {
					config := migration.MakeConfig()
					connection, err := migration.ConnectToDatabase(config)
					if err != nil {
						panic(err)
					}

					migrationLogRepository := migration.NewMigrationLogRepository(connection, config)
					migrationService := migration.MakeMigrationService(directory, config, migrationLogRepository)

					migrationLogRepository.CreateMigrationLogsTable()
					migrationService.Down()

					return nil
				},
			},
			{
				Name:  "make",
				Usage: "Create a new migration",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "name",
						Aliases:  []string{"n"},
						Usage:    "The name of the new migration",
						Required: true,
					},
					&cli.StringFlag{
						Name:    "table",
						Value:   "",
						Aliases: []string{"t"},
						Usage:   "Generate boilerplate code for the given table name",
					},
				},
				Action: func(ctx *cli.Context) error {
					var folder string
					name := ctx.String("name")
					table := ctx.String("table")

					if table == "" {
						folder = migration.GenerateMigration(directory, name)
					} else {
						folder = migration.GenerateMigrationWithTemplate(directory, name, table)
					}

					fmt.Println("Created new migration in", folder)

					return nil
				},
			},
			{
				Name:  "test",
				Usage: "Test configuration and print results",
				Action: func(ctx *cli.Context) error {
					config := migration.MakeConfig()
					fmt.Println(config.Database.Host)
					return nil
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
