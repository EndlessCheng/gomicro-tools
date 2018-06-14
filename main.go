package main

import "gomicro-tools/rpc/handler"

func genRepository() {

}

func main() {
	// TODO: interface + structs

	//srcFilePaths := [...]string{"src/input.struct"}
	//dstProtoFilePaths := [...]string{"out/output.proto"}
	//for i, srcFilePath := range srcFilePaths {
	//	dstProtoFilePath := dstProtoFilePaths[i]
	//
	//	//	rpc.GenProto(srcFilePath, dstProtoFilePath)
	//	rpc.GenMessages(srcFilePath, dstProtoFilePath)
	//}

	//repository.GenRepository("src/input.interface", "out/storage_repository.go", "storage")
	
	handler.GenHandler("src/input.interface", "out/handler.go", "storage")
}
