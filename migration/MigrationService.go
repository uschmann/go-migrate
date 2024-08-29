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
	migrationLogRepository *MigrationLogRepository
}

func MakeMigrationService(dir string, migrationLogRepository *MigrationLogRepository) *MigrationService {
	migrationService := &MigrationService{
		migrationDir:           dir,
		migrationLogRepository: migrationLogRepository,
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
	}

}

func (m *MigrationService) Up() {
	var migrationsToRun []*Migration

	for _, migration := range m.migrations {
		fmt.Println(migration.name)

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

	for _, migration := range migrationsToRun {
		if migration.HasUp {
			runSqlplus(migration.GetUpFilename())
		}
	}
}
