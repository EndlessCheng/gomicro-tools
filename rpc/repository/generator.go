package repository

import (
	"gomicro-tools/common"
	"gomicro-tools/rpc"
)

func GenRepository(srcFilePath string, dstFilePath string, dstName string) {
	sourceCode := common.ReadText(srcFilePath)
	parsedInterface := rpc.ParseInterface(sourceCode)
	if parsedInterface != nil {
		genRepository(dstFilePath, parsedInterface, dstName)
	}
}
