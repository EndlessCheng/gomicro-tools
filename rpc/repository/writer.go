package repository

import (
	"bufio"

	"gomicro-tools/common"
	"gomicro-tools/rpc"
	"fmt"
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
	"os"

	k8s "github.com/micro/kubernetes/go/micro"
	"github.com/micro/go-micro"
)

`)
	common.Check(err)

	writeServiceInterface(w, parsedInterface)

	implStructName := serviceName + "SvcRepository"
	w.WriteString(fmt.Sprintf(`
func New%[1]sSvcRepository() %[1]sSvcRepository {
	// TODO
}

type %s struct {
	client proto.%[1]sService
}

`, serviceNameUpper, implStructName))

	for _, method := range parsedInterface.Methods {
		writeMethod(w, implStructName, method)
		w.WriteString("\n")
	}

	w.Flush()
}
