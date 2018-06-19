package handler

import (
	"gomicro-tools/rpc"
	"strings"
	"gomicro-tools/common"
	"bufio"
	"fmt"
)

var (
	serviceName      string
	serviceNameUpper string
)

func writeMethodReturns(w *bufio.Writer, returns []*rpc.Var) {
	for i, v := range returns {
		w.WriteString(v.Name)
		if i < len(returns)-1 {
			w.WriteString(", ")
		}
	}
}

func writeUseCaseParams(w *bufio.Writer, args []*rpc.Var) {
	for i, v := range args {
		vName := strings.Title(v.Name)

		protoType := rpc.MapGoTypeToProtoType(v.Type)
		if protoType != v.Type {
			w.WriteString(fmt.Sprintf("%s(req.%s)", v.Type, vName))
		} else {
			w.WriteString("req." + vName)
		}

		if i < len(args)-1 {
			w.WriteString(", ")
		}
	}
}

func writeResponseAssign(w *bufio.Writer, returns []*rpc.Var) {
	for _, v := range returns {
		w.WriteString(fmt.Sprintf("%[1]s%[1]sresp.%s = ", common.Tab, strings.Title(v.Name)))
		protoType := rpc.MapGoTypeToProtoType(v.Type)
		if protoType != v.Type {
			w.WriteString(fmt.Sprintf("%s(%s)\n", protoType, v.Name))
		} else {
			w.WriteString(v.Name + "\n")
		}
	}
}

func writeMethod(w *bufio.Writer, structName string, method *rpc.Method) {
	reqType := rpc.RequestMessageType(method.Name)
	respType := rpc.ResponseMessageType(method.Name)
	w.WriteString(fmt.Sprintf("func (h *%s) %s(ctx context.Context, req *proto.%s) (*proto.%s, error) {\n",
		structName, method.Name, reqType, respType))

	w.WriteString(common.Tab)
	writeMethodReturns(w, method.Returns)
	w.WriteString(fmt.Sprintf(" := h.ucase.%s(", method.Name))
	writeUseCaseParams(w, method.Parameters)
	w.WriteString(")\n")

	w.WriteString(common.Tab + fmt.Sprintf("resp := proto.%s{}\n", respType))

	w.WriteString(common.Tab + "resp.ErrCode = model.GetErrorCode(err)\n")

	if len(method.Returns) > 1 {
		w.WriteString(common.Tab + "if err == nil {\n")
		writeResponseAssign(w, method.Returns[:len(method.Returns)-1])
		w.WriteString(common.Tab + "}\n")
	}

	w.WriteString(common.Tab + "return &resp, nil\n")

	w.WriteString("}\n")
}

func genHandler(dstFilePath string, parsedInterface *rpc.InterFace, dstName string) {
	serviceName = dstName
	serviceNameUpper = strings.Title(serviceName)
	implStructName := serviceName + "Handler"

	f := common.CreateFile(dstFilePath)
	defer f.Close()

	w := bufio.NewWriter(f)

	w.WriteString(fmt.Sprintf(`package rpc

import (
	"context"

)

func New%[1]sHandler(ucase usecase.%[1]sUseCase) proto.%[1]sServer {
	return &%[2]s{ucase}
}

type %[2]s struct {
	ucase usecase.%[1]sUseCase
}

`, serviceNameUpper, implStructName))

	for _, method := range parsedInterface.Methods {
		writeMethod(w, implStructName, method)
		w.WriteString("\n")
	}

	w.Flush()
}
