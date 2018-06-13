package main

import "gomicro-tools/rpc"

func main() {
	srcFilePath := "input.txt"
	dstProtoFilePath := "out/output.proto"
	rpc.GenProto(srcFilePath, dstProtoFilePath)
}
