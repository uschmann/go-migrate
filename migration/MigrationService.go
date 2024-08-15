package migration

import (
	"log"
	"os"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

type MigrationService struct {
	migrationDir string
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

	for _, entry := range entries {
		log.Println(entry)
	}
}
