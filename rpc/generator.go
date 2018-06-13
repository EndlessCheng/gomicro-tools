package rpc

import "io/ioutil"

func readText(srcFilePath string) string {
	data, err := ioutil.ReadFile(srcFilePath)
	check(err)

	return string(data)
}

func GenProto(srcFilePath string, dstProtoFilePath string) {
	sourceCode := readText(srcFilePath)
	parsedInterface := parseInterface(sourceCode)
	if parsedInterface != nil {
		genProto(dstProtoFilePath, parsedInterface.Methods)
	}
}

func GenMessages(srcFilePath string, dstProtoFilePath string) {
	sourceCode := readText(srcFilePath)
	parsedStructs := parseStructs(sourceCode)
	if parsedStructs != nil {
		genMessages(dstProtoFilePath, parsedStructs)
	}
}