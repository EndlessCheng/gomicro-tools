package main

import "gomicro-tools/rpc"

func main() {
	// TODO: interface + structs

	srcFilePaths := [...]string{"src/input.struct"}
	dstProtoFilePaths := [...]string{"out/output.proto"}
	for i, srcFilePath := range srcFilePaths {
		dstProtoFilePath := dstProtoFilePaths[i]

		//	rpc.GenProto(srcFilePath, dstProtoFilePath)
		rpc.GenMessages(srcFilePath, dstProtoFilePath)
	}
}
