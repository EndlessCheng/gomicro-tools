package main

import (
	"fmt"

	"os"
	"github.com/urfave/cli"
	"log"
	"gomicro-tools/rpc/handler"
	httpHandler "gomicro-tools/rpc/handler/http"

	"gomicro-tools/rpc"
	"gomicro-tools/rpc/repository"
	gmain "gomicro-tools/rpc/gmain"
	"gomicro-tools/common"
	"bufio"
)

const (
	flagForce = "force"
	flagAll   = "all"

	flagInit        = "init"
	flagMain        = "main"
	flagProto       = "proto"
	flagProtoRepo   = "proto_repository"
	flagRPCHandler  = "rpc_handler"
	flagHTTPHandler = "http_handler"
	flagDeploy      = "deploy"
)

const (
	ucaseFilePath  = "proto/interface/usecase.go"
	modelsFilePath = "proto/interface/models.go"
)

var (
	projectName string

	//force bool
	genInitFiles     bool
	genMain          bool
	genProto         bool
	genSvcRepository bool
	genRPCHandler    bool
	genHTTPHandler   bool

	genDeploy bool
)

func initGoFile(filePath string, packageName string) {
	f, err := os.Create(filePath)
	common.Check(err)

	w := bufio.NewWriter(f)
	w.WriteString("package " + packageName + "\n\n\n")
	w.Flush()
}

func gen() {
	if genInitFiles {
		err := os.MkdirAll("proto/interface", os.ModePerm)
		common.Check(err)
		initGoFile(ucaseFilePath, "_interface")
		initGoFile(modelsFilePath, "_interface")
		return
	}

	if genMain {
		dstPath := "main.go"
		fmt.Print("generating ", dstPath, " ...")
		gmain.GenMain(ucaseFilePath, dstPath, projectName)
		fmt.Println("Done")
	}

	if genProto {
		dstPath := fmt.Sprintf("proto/%s.proto", projectName)
		fmt.Print("generating ", dstPath, " ...")
		rpc.GenFullProto(ucaseFilePath, dstPath)
		fmt.Println("Done")
	}

	if genSvcRepository {
		dstPath := fmt.Sprintf("proto/repository/%s_repository.go", projectName)
		fmt.Print("generating ", dstPath, " ...")
		repository.GenRepository(ucaseFilePath, dstPath, projectName)
		fmt.Println("Done")
	}

	if genRPCHandler {
		// TODO: 完善 newXXXModel
		dstPath := fmt.Sprintf("handler/rpc/%s_handler.go", projectName)
		fmt.Print("generating ", dstPath, " ...")
		handler.GenHandler(ucaseFilePath, dstPath, projectName)
		fmt.Println("Done")
	}

	if genHTTPHandler {
		dstPath := fmt.Sprintf("handler/http/%s_handler.go", projectName)
		fmt.Print("generating ", dstPath, " ...")
		httpHandler.GenHandler(ucaseFilePath, dstPath, projectName)
		fmt.Println("Done")
	}
}

func main() {
	// 1. 修改 projectName
	// 2. 修改 input.interface，移除不需要的方法声明，id uint 改成 hashID string
	// 3. 修改 input.struct，ID uint 改成 HashID string，其他字段看情况
	// 4. 生成的 .proto service 名字可能会与 message 的名字重名，此情况需要修改 service 名字

	app := cli.NewApp()
	app.Usage = "微服务代码生成"
	app.Version = "0.0.1"

	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:  flagForce,
			Usage: "generate and replace old files",
		},

		cli.BoolFlag{
			Name:  flagInit,
			Usage: "create interface folder",
		},
		cli.BoolFlag{
			Name:  flagAll,
			Usage: "generate ALL files",
		},
		cli.BoolFlag{
			Name:  flagMain,
			Usage: "generate main.go",
		},
		cli.BoolFlag{
			Name:  flagProto,
			Usage: "generate .proto",
		},
		cli.BoolFlag{
			Name:  flagProtoRepo,
			Usage: "generate proto service repository",
		},
		cli.BoolFlag{
			Name:  flagRPCHandler,
			Usage: "generate RPC handler",
		},
		cli.BoolFlag{
			Name:  flagHTTPHandler,
			Usage: "generate HTTP handler",
		},
		//cli.BoolFlag{
		//	Name:  flagDeploy,
		//	Usage: "generate deploy files",
		//},
	}

	app.Action = func(c *cli.Context) error {
		projectName = c.Args().Get(0)
		if projectName == "" {
			fmt.Println(`Please enter the project name(e.g. switch, video, etc.).
Usage example: $ gomicro-tools -all <projectName>
More infomation: $ gomicro-tools help`)
			os.Exit(1)
		}

		common.Force = c.Bool(flagForce)

		genInitFiles = c.Bool(flagInit)
		all := c.Bool(flagAll)
		genMain = all || c.Bool(flagMain)
		genProto = all || c.Bool(flagProto)
		genSvcRepository = all || c.Bool(flagProtoRepo)
		genRPCHandler = all || c.Bool(flagRPCHandler)
		genHTTPHandler = all || c.Bool(flagHTTPHandler)
		genDeploy = all || c.Bool(flagDeploy)

		gen()

		return nil
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
