package handler

import (
	"gomicro-tools/common"
	"gomicro-tools/rpc"
)

func GenHandler(srcFilePath string, dstFilePath string, dstName string) {
	sourceCode := common.ReadText(srcFilePath)
	parsedInterface := rpc.ParseInterface(sourceCode)
	if parsedInterface != nil {
		genHandler(dstFilePath, parsedInterface, dstName)
	}
}
