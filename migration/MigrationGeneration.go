package migration

import (
	"bytes"
	_ "embed"
	"os"
	"path"
	"text/template"
	"time"
)

//go:embed templates/up.sql
var upTemplate string

//go:embed templates/down.sql
var downTemplate string

func formatSqlTemplate(text string, name string) string {
	var doc bytes.Buffer
	templ := template.Must(template.New("myname").Parse(text))
	err := templ.Execute(&doc, map[string]interface{}{
		"name": name,
	})

	check(err)

	return doc.String()
}

func generatePrefix() string {
	return time.Now().Format("2006_01_02_150405_")
}

func GenerateMigration(migrationDir string, name string) string {
	timestamp := generatePrefix()

	folderPath := path.Join(migrationDir, timestamp+name)
	upPath := path.Join(folderPath, "up.sql")
	downPath := path.Join(folderPath, "down.sql")

	err := os.Mkdir(folderPath, os.ModePerm)
	check(err)

	f, err := os.Create(upPath)
	check(err)
	defer f.Close()

	f, err = os.Create(downPath)
	check(err)
	defer f.Close()

	return folderPath
}

func GenerateMigrationWithTemplate(migrationDir string, name string, table string) string {
	timestamp := generatePrefix()

	folderPath := path.Join(migrationDir, timestamp+name)
	upPath := path.Join(folderPath, "up.sql")
	downPath := path.Join(folderPath, "down.sql")

	err := os.Mkdir(folderPath, os.ModePerm)
	check(err)

	err = os.WriteFile(upPath, []byte(formatSqlTemplate(upTemplate, table)), os.ModePerm)
	check(err)

	err = os.WriteFile(downPath, []byte(formatSqlTemplate(downTemplate, table)), os.ModePerm)
	check(err)

	return folderPath
}
