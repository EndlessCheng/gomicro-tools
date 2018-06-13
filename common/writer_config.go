package common

import (
	"strings"
	"regexp"
)

var Tab = strings.Repeat(" ", 4)

func LowerHead(str string) string {
	r, _ := regexp.Compile("^[A-Z]+")

	loc := r.FindStringIndex(str)
	if loc == nil {
		return str
	}

	pos := loc[1]
	return strings.ToLower(str[:pos]) + str[pos:]
}
