package rpc

import (
	"gomicro-tools/common"
)

func GenProto(srcFilePath string, dstProtoFilePath string) {
	sourceCode := common.ReadText(srcFilePath)
	parsedInterface := ParseInterface(sourceCode)
	if parsedInterface != nil {
		genProto(dstProtoFilePath, parsedInterface.Methods)
	}
}

func GenMessages(srcFilePath string, dstProtoFilePath string) {
	sourceCode := common.ReadText(srcFilePath)
	parsedStructs := parseStructs(sourceCode)
	if parsedStructs != nil {
		genMessages(dstProtoFilePath, parsedStructs)
	}
}
