package rpc

import "io/ioutil"

func GenRpc(srcFilePath string, rpcFilePath string) {
	data, err := ioutil.ReadFile(srcFilePath)
	check(err)

	parsedInterface := parseInterface(string(data))
	if parsedInterface != nil {
		genRpc(rpcFilePath, parsedInterface.Methods)
	}
}
