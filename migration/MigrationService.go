package migration

import (
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
)

type MigrationService struct {
	migrationDir           string
	migrations             []*Migration
	migrationsByName       map[string]*Migration
	migrationLogRepository *MigrationLogRepository
	config                 *Config
}

func MakeMigrationService(dir string, config *Config, migrationLogRepository *MigrationLogRepository) *MigrationService {
	migrationService := &MigrationService{
		migrationDir:           dir,
		migrationLogRepository: migrationLogRepository,
		migrationsByName:       make(map[string]*Migration),
		config:                 config,
	}

	migrationService.readDir()

	return migrationService
}

func (m *MigrationService) readDir() {
	entries, err := os.ReadDir(m.migrationDir)

	if err != nil {
		log.Fatal(err)
	}

	absPath, _ := filepath.Abs(m.migrationDir)

	for _, dir := range entries {
		migration := MakeMigration(path.Join(absPath, dir.Name()))
		m.migrations = append(m.migrations, migration)
		m.migrationsByName[migration.Name] = migration
	}
}

func (m *MigrationService) getMigrationByName(name string) *Migration {
	return m.migrationsByName[name]
}

func (m *MigrationService) Up() {
	var migrationsToRun []*Migration

	for _, migration := range m.migrations {

		isExecuted, err := m.migrationLogRepository.IsMigrationExecuted(migration.Name)

		if err != nil {
			panic(err)
		}

		if !isExecuted {
			migrationsToRun = append(migrationsToRun, migration)
		}
	}

	if len(migrationsToRun) == 0 {
		fmt.Println("Nothing to migrate")
	}

	batch, err := m.migrationLogRepository.GetHighestBatch()
	check(err)
	batch++

	for _, migration := range migrationsToRun {
		if migration.HasUp {
			fmt.Println("Migrating " + migration.Name)

			tempDir := copyMigrationToTemp(migration)
			//defer os.RemoveAll(tempDir)

			stdout, stderr, err := execute(m.config, path.Join(tempDir, "wrapper.sql"), path.Join(tempDir, "up.sql"))

			if err != nil {
				fmt.Println(stdout)
				fmt.Println(stderr)
				fmt.Println(tempDir)
				panic(err)
			}

			m.migrationLogRepository.AddMigrationLog(migration.Name, batch)
		}
	}
}

func (m *MigrationService) Down() {
	migrationNames, err := m.migrationLogRepository.GetMigrationLogsToRollback()
	check(err)

	if len(migrationNames) == 0 {
		fmt.Println("Nothing to rollback")
	}

	for _, name := range migrationNames {
		migration := m.getMigrationByName(name)

		if migration.HasDown {
			fmt.Println("Rolling back " + migration.Name)

			tempDir := copyMigrationToTemp(migration)

			stdout, _, err := execute(m.config, path.Join(tempDir, "wrapper.sql"), path.Join(tempDir, "down.sql"))

			if err != nil {
				fmt.Println(stdout)
				panic(err)
			}

			m.migrationLogRepository.DeleteMigrationLogByName(name)
		}
	}

	return
}

type MigrationStatus struct {
	Migration  *Migration
	IsExecuted bool
	Index      int
	Batch      int
}

func (m *MigrationService) GetMigrationStatus() []MigrationStatus {
	migrationStatus := make([]MigrationStatus, 0)
	migrationLogs, err := m.migrationLogRepository.GetAllMigrationLogs()
	check(err)

	mappedMigrationLogs := make(map[string]MigrationLog)
	for _, migrationLog := range migrationLogs {
		mappedMigrationLogs[migrationLog.Name] = migrationLog
	}

	for index, migration := range m.migrations {
		migrationLog, ok := mappedMigrationLogs[migration.Name]
		batch := 0
		isExecuted := false

		if ok {
			batch = migrationLog.Batch
			isExecuted = true
		}

		migrationStatus = append(migrationStatus, MigrationStatus{
			Migration:  migration,
			IsExecuted: isExecuted,
			Index:      index + 1,
			Batch:      batch,
		})
	}

	return migrationStatus
}
