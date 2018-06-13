package rpc

import "strings"

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
	switch goName {
	case "err":
		return "errCode"
	default:
		return goName
	}
}
