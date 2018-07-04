package main

import (
	"fmt"

	"gomicro-tools/rpc/handler"
	"gomicro-tools/rpc/repository"
	"gomicro-tools/rpc"
)

func main() {
	projectName := "switch"

	// TODO: 修改 id uint => hashID string
	srcFilePath := "src/input.interface"
	dstProtoFilePath := fmt.Sprintf("out/%s.proto", projectName)
	rpc.GenProto(srcFilePath, dstProtoFilePath)

	// TODO: 修改 id uint => hashID string
	// TODO: 合并到 .proto 中
	srcFilePath = "src/input.struct"
	dstProtoFilePath = fmt.Sprintf("out/%s.message.proto", projectName)
	rpc.GenMessages(srcFilePath, dstProtoFilePath)

	// TODO: 注释也复制过来
	repository.GenRepository("src/input.interface", fmt.Sprintf("out/%s_repository.go", projectName), projectName)

	// TODO: 完善 newXXXModel
	// TODO: 处理各种 hashID, <prefix>HashID
	handler.GenHandler("src/input.interface", fmt.Sprintf("out/%s_handler.go", projectName), projectName)
}
