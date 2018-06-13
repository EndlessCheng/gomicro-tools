package rpc

import "io/ioutil"

func GenProto(srcFilePath string, rpcFilePath string) {
	data, err := ioutil.ReadFile(srcFilePath)
	check(err)

	parsedInterface := parseInterface(string(data))
	if parsedInterface != nil {
		genProto(rpcFilePath, parsedInterface.Methods)
	}
}
