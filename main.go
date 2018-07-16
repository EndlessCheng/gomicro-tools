package main

import (
	"fmt"

	"gomicro-tools/rpc"
	"os"
	"github.com/urfave/cli"
	"log"
	"gomicro-tools/rpc/handler"
	httpHandler "gomicro-tools/rpc/handler/http"

	"gomicro-tools/rpc/repository"
)

func main() {
	// 1. 修改 projectName
	// 2. 修改 input.interface，移除不需要的方法声明，id uint 改成 hashID string
	// 3. 修改 input.struct，ID 类型改成 string，其他字段看情况
	// 4. 生成的 .proto service 名字可能会与 message 的名字重名，此情况需要修改 service 名字
	var projectName string

	app := cli.NewApp()

	app.Action = func(c *cli.Context) error {
		projectName = c.Args().Get(0)
		return nil
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}

	ucaseFilePath := "proto/interface/usecase.go"
	dstProtoFilePath := fmt.Sprintf("proto/%s.proto", projectName)
	rpc.GenFullProto(ucaseFilePath, dstProtoFilePath)

	repository.GenRepository(ucaseFilePath, fmt.Sprintf("proto/repository/%s_repository.go", projectName), projectName)

	// TODO: 完善 newXXXModel
	handler.GenHandler(ucaseFilePath, fmt.Sprintf("handler/rpc/%s_handler.go", projectName), projectName)

	httpHandler.GenHandler(ucaseFilePath, fmt.Sprintf("handler/http/%s_handler.go", projectName), projectName)
}
