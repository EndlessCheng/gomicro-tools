package main

import (
	"fmt"

	"gomicro-tools/rpc/handler"
	"gomicro-tools/rpc/repository"
	"gomicro-tools/rpc"
)

func main() {
	// 1. 修改 projectName
	// 2. 修改 input.interface，移除不需要的方法声明，id uint 改成 hashID string
	// 3. 修改 input.struct，ID 类型改成 string，其他字段看情况
	// 4. 生成的 .proto service 名字可能会与 message 的名字重名，此情况需要修改 service 名字
	projectName := "video"

	srcFilePath := "src/input.interface"
	dstProtoFilePath := fmt.Sprintf("out/%s.proto", projectName)
	rpc.GenProto(srcFilePath, dstProtoFilePath)

	// TODO: 合并到 .proto 中
	srcFilePath = "src/input.struct"
	dstProtoFilePath = fmt.Sprintf("out/%s.message.proto", projectName)
	rpc.GenMessages(srcFilePath, dstProtoFilePath)

	repository.GenRepository("src/input.interface", fmt.Sprintf("out/%s_repository.go", projectName), projectName)

	// TODO: 完善 newXXXModel
	handler.GenHandler("src/input.interface", fmt.Sprintf("out/%s_handler.go", projectName), projectName)
}
