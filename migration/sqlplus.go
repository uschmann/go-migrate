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

func execute(wrapper string, script string) (string, string, error) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	cmd := exec.Command("sqlplus", "auschmann/secret@localhost:1522/FREE", "@"+wrapper, script)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	return stdout.String(), stderr.String(), err
}

func copyMigrationToTemp(migration *Migration) string {
	tempDir, err := os.MkdirTemp("", "db-migrate")
	check(err)
	//defer os.RemoveAll(tempDir)

	err = gorecurcopy.CopyDirectory(migration.directory, tempDir)
	check(err)

	// Write wrapper to temp directory
	tempWrapper := filepath.Join(tempDir, "wrapper.sql")
	err = os.WriteFile(tempWrapper, []byte(wrapperScript), 0644)
	check(err)

	return tempDir
}

func runSqlplus(wrapper string, script string) {
	// Run sqlplus
	stdout, _, err := execute(wrapper, script)

	fmt.Println(stdout)
	check(err)
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
