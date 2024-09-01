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
}

func MakeMigrationService(dir string, migrationLogRepository *MigrationLogRepository) *MigrationService {
	migrationService := &MigrationService{
		migrationDir:           dir,
		migrationLogRepository: migrationLogRepository,
		migrationsByName:       make(map[string]*Migration),
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
		m.migrationsByName[migration.name] = migration
	}
}

func (m *MigrationService) getMigrationByName(name string) *Migration {
	return m.migrationsByName[name]
}

func (m *MigrationService) Up() {
	var migrationsToRun []*Migration

	for _, migration := range m.migrations {

		isExecuted, err := m.migrationLogRepository.IsMigrationExecuted(migration.name)

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
			tempDir := copyMigrationToTemp(migration)

			stdout, _, err := execute(path.Join(tempDir, "wrapper.sql"), path.Join(tempDir, "up.sql"))

			if err != nil {
				fmt.Println(stdout)
				panic(err)
			}

			m.migrationLogRepository.AddMigrationLog(migration.name, batch)
		}
	}
}

func (m *MigrationService) Down() {
	migrationNames, err := m.migrationLogRepository.GetMigrationLogsToRollback()
	check(err)

	for _, name := range migrationNames {
		migration := m.getMigrationByName(name)

		if migration.HasDown {
			tempDir := copyMigrationToTemp(migration)

			stdout, _, err := execute(path.Join(tempDir, "wrapper.sql"), path.Join(tempDir, "down.sql"))

			if err != nil {
				fmt.Println(stdout)
				panic(err)
			}

			m.migrationLogRepository.DeleteMigrationLogByName(name)
		}
	}

	return
}
