package migration

import (
	"bytes"
	_ "embed"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
)

//go:embed templates/wrapper.sql
var wrapperScript string

func execute(wrapper string, script string) (string, string, error) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	fmt.Println()

	cmd := exec.Command("sqlplus", "auschmann/secret@localhost:1522/FREE", "@"+wrapper, script)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	return stdout.String(), stderr.String(), err
}

func runSqlplus(script string) {
	tempDir, err := os.MkdirTemp("", "db-migrate")
	check(err)
	//defer os.RemoveAll(tempDir)

	// Write wrapper to temp directory
	tempWrapper := filepath.Join(tempDir, "wrapper.sql")
	err = os.WriteFile(tempWrapper, []byte(wrapperScript), 0644)
	check(err)

	// Copy the script to the temp directory
	tempScript := filepath.Join(tempDir, "script.sql")
	_, err = copyFile(script, tempScript)
	check(err)

	// Run sqlplus
	stdout, _, err := execute(tempWrapper, tempScript)

	check(err)

	fmt.Println(stdout)

	fmt.Println(tempDir)
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
