package common

import (
	"regexp"
	"strings"
	"path/filepath"
	"os"
)

var Tab = strings.Repeat(" ", 4)

var Force bool

var ProjectImportPrefix string

func init() {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	Check(err)

	i := strings.LastIndex(dir, "/go/src/")
	ProjectImportPrefix = dir[i+len("/go/src/"):]

	// or GOPATH or HOME + "/go"
}

func LowerHead(str string) string {
	r, _ := regexp.Compile("^[A-Z]+")

	loc := r.FindStringIndex(str)
	if loc == nil {
		return str
	}

	pos := loc[1]
	return strings.ToLower(str[:pos]) + str[pos:]
}
