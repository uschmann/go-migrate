package migration

import (
	"log"
	"os"
	"path"
	"path/filepath"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

type MigrationService struct {
	migrationDir string
	migrations   []*Migration
}

func MakeMigrationService(dir string) *MigrationService {
	migrationService := &MigrationService{
		migrationDir: dir,
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
