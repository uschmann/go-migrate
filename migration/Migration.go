package migration

import (
	"errors"
	"os"
	"path"
	"path/filepath"
)

type Migration struct {
	name      string
	directory string
	HasUp     bool
	HasDown   bool
}

func MakeMigration(dir string) *Migration {
	migration := &Migration{
		name:      filepath.Base(dir),
		directory: dir,
		HasUp:     fileExists(path.Join(dir, "up.sql")),
		HasDown:   fileExists(path.Join(dir, "down.sql")),
	}

	return migration
}

func fileExists(filename string) bool {
	if _, err := os.Stat(filename); errors.Is(err, os.ErrNotExist) {
		return false
	} else {
		return true
	}
}
