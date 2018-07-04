package rpc

import "gomicro-tools/common"

func RequestMessageType(methodName string) string {
	return "Req" + methodName
}

func ResponseMessageType(methodName string) string {
	return "Resp" + methodName
}

func MapGoTypeToProtoType(goTypeName string) string {
	switch goTypeName {
	case "int":
		return "int64"
	case "uint":
		return "uint64"
	case "float32":
		return "float"
	case "float64":
		return "double"
	case "error":
		return "int32"
	default:
		return goTypeName
	}
}

func MapGoNameToProtoName(goName string) string {
	goName = common.LowerHead(goName)
	switch goName {
	case "err":
		return "errCode"
	default:
		return goName
	}
}
