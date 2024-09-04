package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/pterm/pterm"
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

					tableData := pterm.TableData{
						{"#", "Name", "Executed", "Batch"},
					}

					for _, status := range migrationStatus {
						isExecuted := pterm.Red("No")
						if status.IsExecuted {
							isExecuted = pterm.Green("Yes")
						}

						batch := strconv.Itoa(status.Batch)
						if batch == "0" {
							batch = "-"
						}

						tableData = append(tableData, []string{
							strconv.Itoa(status.Index),
							status.Migration.Name,
							isExecuted,
							batch,
						})
					}
					pterm.DefaultTable.WithHasHeader().WithData(tableData).Render()

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

					// Print configuration
					pterm.DefaultBasicText.Print(pterm.LightCyan("Printing configuration:\n"))
					databaseTableData := pterm.TableData{
						{"Name", "Value"},
						{pterm.LightMagenta("Database"), ""},
						{"DB_USER", config.Database.User},
						{"DB_PASSWORD", "*****"},
						{"DB_HOST", config.Database.Host},
						{"DB_PORT", strconv.Itoa(config.Database.Port)},
						{"DB_SERVICE", config.Database.Service},
						{pterm.LightMagenta("Misc"), ""},
						{"DB_MIGRATION_LOG_TABLE", config.MigrationLogTable},
						{"SQLPLUS_BIN", config.Sqlplus},
					}
					pterm.DefaultTable.WithHasHeader().WithData(databaseTableData).Render()

					pterm.DefaultBasicText.Print(pterm.LightCyan("Check sqlplus:\n") + "Running " + pterm.LightMagenta(config.Sqlplus+" -v") + "\n")

					sqlPLusVersion, stderr, err := migration.GetSqlplusVersion(config)
					if err != nil {
						pterm.DefaultBasicText.Println(stderr)
						panic(err)
					}

					pterm.DefaultBasicText.Println(sqlPLusVersion)

					pterm.DefaultBasicText.Println(pterm.White("sqlplus was executed " + pterm.Green("successfully\n")))

					// Connect to database

					pterm.DefaultBasicText.Print(pterm.LightCyan("Checking database connection:\n"))
					connection, err := migration.ConnectToDatabase(config)
					if err != nil {
						panic(err)
					}
					migrationLogRepository := migration.NewMigrationLogRepository(connection, config)

					err = migrationLogRepository.TestConnection()
					if err != nil {
						panic(err)
					}
					pterm.DefaultBasicText.Print(pterm.White("Database conenction with go-ora was " + pterm.Green("successfull\n")))

					if migrationLogRepository.MigrationLogsTableExists() {
						pterm.DefaultBasicText.Println(config.MigrationLogTable + " table " + pterm.Green("exists"))
					} else {
						pterm.DefaultBasicText.Println(config.MigrationLogTable + " table " + pterm.Red("does not exist"))
					}

					return nil
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
