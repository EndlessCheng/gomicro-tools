package repository

import (
	"bufio"

	"fmt"
	"gomicro-tools/common"
	"gomicro-tools/rpc"
	"strings"
)

var (
	serviceName      string
	serviceNameUpper string
)

func writeVar(w *bufio.Writer, v *rpc.Var) {
	vName := v.Name
	vType := v.Type
	if vType == "error" {
		vName = "errCode"
		vType = "int32"
	}

	w.WriteString(vName + " ")

	if v.IsSlice {
		w.WriteString("[]")
	}

	c0 := vType[:1]
	if c0 == strings.ToUpper(c0) {
		vType = "*proto." + vType
	}
	w.WriteString(vType)
}

func writeVars(w *bufio.Writer, vs []*rpc.Var, isProto bool) {
	w.WriteString("(")

	for i, v := range vs {
		writeVar(w, v)

		if i < len(vs)-1 {
			w.WriteString(", ")
		}
	}

	if isProto {
		w.WriteString(", err error")
	}

	w.WriteString(")")
}

func writeStructInitParams(w *bufio.Writer, vs []*rpc.Var) {
	for i, v := range vs {
		w.WriteString(strings.Title(v.Name) + ": ")

		protoType := rpc.MapGoTypeToProtoType(v.Type)
		if protoType != v.Type {
			w.WriteString(fmt.Sprintf("%s(%s)", protoType, v.Name))
		} else {
			w.WriteString(v.Name)
		}

		if i < len(vs)-1 {
			w.WriteString(", ")
		}
	}
}

func writeMethodReturns(w *bufio.Writer, vs []*rpc.Var) {
	for _, v := range vs {
		vName := v.Name
		vName = rpc.MapGoNameToProtoName(vName)
		respGetStr := fmt.Sprintf("resp.Get%s()", strings.Title(vName))

		protoType := rpc.MapGoTypeToProtoType(v.Type)
		if protoType != v.Type && v.Type != "error" {
			w.WriteString(fmt.Sprintf("%s(%s)", v.Type, respGetStr))
		} else {
			w.WriteString(respGetStr)
		}

		w.WriteString(", ")
	}
	w.WriteString("err")
}

func writeServiceInterface(w *bufio.Writer, parsedInterface *rpc.InterFace) {
	w.WriteString(fmt.Sprintf("type %sSvcRepository interface {\n", serviceNameUpper))

	for _, method := range parsedInterface.Methods {
		w.WriteString(fmt.Sprintf("%s%s", common.Tab, method.Name))
		writeVars(w, method.Parameters, false)
		w.WriteString(" ")
		writeVars(w, method.Returns, true)
		w.WriteString("\n")
	}

	w.WriteString("}\n")
}

func writeMethod(w *bufio.Writer, structName string, method *rpc.Method) {
	w.WriteString(fmt.Sprintf("func (r *%s) %s", structName, method.Name))
	writeVars(w, method.Parameters, false)
	w.WriteString(" ")
	writeVars(w, method.Returns, true)
	w.WriteString(" {\n")

	w.WriteString(fmt.Sprintf("%sreq := proto.%s{", common.Tab, rpc.RequestMessageType(method.Name)))
	writeStructInitParams(w, method.Parameters)
	w.WriteString("}\n")

	w.WriteString(fmt.Sprintf("%sresp, err := r.client.%s(context.TODO(), &req)\n", common.Tab, method.Name))

	w.WriteString(common.Tab + "if err != nil {\n")
	w.WriteString(common.Tab + common.Tab + fmt.Sprintf("log.WithError(err).Errorln(\"[grpc.%s] error with args:\"", method.Name))
	for _, arg := range method.Parameters {
		w.WriteString(", " + arg.Name)
	}
	w.WriteString(")\n")
	w.WriteString(common.Tab + "}\n")

	w.WriteString(common.Tab + "return ")
	writeMethodReturns(w, method.Returns)
	w.WriteString("\n}\n")
}

func genRepository(dstFilePath string, parsedInterface *rpc.InterFace, dstName string) {
	serviceName = dstName
	serviceNameUpper = strings.Title(serviceName)

	f := common.CreateFile(dstFilePath)
	defer f.Close()

	w := bufio.NewWriter(f)

	_, err := w.WriteString(`package repository

import (
	"context"

	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"


)

`)
	common.Check(err)

	writeServiceInterface(w, parsedInterface)

	implStructName := serviceName + "SvcRepository"
	w.WriteString(fmt.Sprintf(`
func New%[2]sSvcRepository() %[2]sSvcRepository {
	// TODO: fill the addr
	%[1]sServiceAddr := utils.GetEnvWithDefault("", "")
	conn, err := grpc.Dial(%[1]sServiceAddr, grpc.WithInsecure())
	if err != nil {
		log.WithError(err).Fatalln("连接 %[1]s 微服务失败")
	}

	client := proto.New%[2]sClient(conn)
	return &%[3]s{client}
}

type %[3]s struct {
	client proto.%[2]sClient
}
`, serviceName, serviceNameUpper, implStructName)) // （可能是环境变量配置错误）

	for _, method := range parsedInterface.Methods {
		w.WriteString("\n")
		writeMethod(w, implStructName, method)
	}

	w.Flush()
}
