package migration

import (
	_ "embed"
	"fmt"
)

//go:embed templates/wrapper.sql
var wrapperScript string

func runSqlplus(script string) {
	fmt.Println(wrapperScript)
}
