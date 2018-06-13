package rpc

import (
	"strings"
	"regexp"
)

var tab = strings.Repeat(" ", 4)

func requestMessageType(methodName string) string {
	return "Req" + methodName
}

func responseMessageType(methodName string) string {
	return "Resp" + methodName
}

func mapGoTypeToProtoType(goTypeName string) string {
	switch goTypeName {
	case "int":
		return "int64"
	case "uint":
		return "uint64"
	case "error":
		return "int32"
	default:
		return goTypeName
	}
}

func mapGoNameToProtoName(goName string) string {
	goName = lowerHead(goName)
	switch goName {
	case "err":
		return "errCode"
	default:
		return goName
	}
}

func lowerHead(str string) string {
	r, _ := regexp.Compile("^[A-Z]+")

	loc := r.FindStringIndex(str)
	if loc == nil {
		return str
	}

	pos := loc[1]
	return strings.ToLower(str[:pos]) + str[pos:]
}
