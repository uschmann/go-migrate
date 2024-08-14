package migration

import (
	"log"
	"os"
	"path"
	"time"
)

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

func (m *MigrationService) GenerateMigration(name string) {
	timestamp := time.Now().Format("2006_01_02_150405_")

	folderPath := path.Join(m.migrationDir, timestamp+name)
	os.Mkdir(folderPath, os.ModePerm)
}
