package rpc

import (
	"gomicro-tools/common"
	"path"
	"io/ioutil"
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

func GenFullProto(srcFilePath string, dstProtoFilePath string) {
	sourceCode := common.ReadText(srcFilePath)
	parsedInterface := ParseInterface(sourceCode)

	dirName := path.Dir(srcFilePath)
	fis, err := ioutil.ReadDir(dirName)
	common.Check(err)

	sourceCodes := make([]string, len(fis))
	for i, fi := range fis {
		modelPath := dirName + "/" + fi.Name()
		sourceCodes[i] = common.ReadText(modelPath)
	}
	parsedStructs := parseStructsFromCodes(sourceCodes)

	if parsedInterface != nil {
		genFullProto(dstProtoFilePath, parsedInterface.Methods, parsedStructs)
	}
}
