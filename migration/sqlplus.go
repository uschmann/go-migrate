package migration

import (
	"bytes"
	_ "embed"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"

	gorecurcopy "github.com/uschmann/go-migrate/utils"
)

//go:embed templates/wrapper.sql
var wrapperScript string

func execute(config *Config, wrapper string, script string) (string, string, error) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	connectionString := config.Database.BuildSqlplusConnectionString()
	cmd := exec.Command(config.Sqlplus, connectionString, "@"+wrapper, script)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	// Todo: Unset http_proxy
	// https://stackoverflow.com/questions/41133115/pass-env-var-to-exec-command

	err := cmd.Run()
	return stdout.String(), stderr.String(), err
}

func copyMigrationToTemp(migration *Migration) string {
	tempDir, err := os.MkdirTemp("", "db-migrate")
	check(err)

	err = gorecurcopy.CopyDirectory(migration.directory, tempDir)
	check(err)

	// Write wrapper to temp directory
	tempWrapper := filepath.Join(tempDir, "wrapper.sql")
	err = os.WriteFile(tempWrapper, []byte(wrapperScript), 0644)
	check(err)

	return tempDir
}

func copyFile(src, dst string) (int64, error) {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return 0, err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return 0, fmt.Errorf("%s is not a regular file", src)
	}

	source, err := os.Open(src)
	if err != nil {
		return 0, err
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
		return 0, err
	}
	defer destination.Close()
	nBytes, err := io.Copy(destination, source)
	return nBytes, err
}
