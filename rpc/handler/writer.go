package handler

import (
	"bufio"
	"fmt"
	"gomicro-tools/common"
	"gomicro-tools/rpc"
	"strings"
)

const (
	idVarName     = "id"
	hashIDVarName = "hashID"
)

var (
	serviceName      string
	serviceNameUpper string
)

func writeMethodReturns(w *bufio.Writer, returns []*rpc.Var) {
	for i, v := range returns {
		vName := v.Name
		if vName == hashIDVarName {
			vName = idVarName
		}

		w.WriteString(vName)
		if i < len(returns)-1 {
			w.WriteString(", ")
		}
	}
}

func writeUseCaseParams(w *bufio.Writer, args []*rpc.Var) {
	for i, v := range args {
		var vName string

		if v.Name == hashIDVarName {
			vName = idVarName
		} else {
			protoType := rpc.MapGoTypeToProtoType(v.Type)
			if protoType != v.Type {
				vName = fmt.Sprintf("%s(req.%s)", v.Type, strings.Title(v.Name))
			} else {
				vName = "req." + strings.Title(v.Name)
			}
		}

		w.WriteString(vName)
		if i < len(args)-1 {
			w.WriteString(", ")
		}
	}
}

func writeResponseAssign(w *bufio.Writer, returns []*rpc.Var) {
	for _, v := range returns {
		if v.Name == hashIDVarName {
			w.WriteString(common.Tab + common.Tab + fmt.Sprintf("%s, _ := utils.EncodeID(%s)\n", hashIDVarName, idVarName))
		}

		w.WriteString(common.Tab + common.Tab + fmt.Sprintf("resp.%s = ", strings.Title(v.Name)))

		vName := v.Name
		protoType := rpc.MapGoTypeToProtoType(v.Type)
		if protoType != v.Type {
			vName = fmt.Sprintf("%s(%s)", protoType, v.Name)
		}

		w.WriteString(vName + "\n")
	}
}

func writeMethod(w *bufio.Writer, structName string, method *rpc.Method) {
	reqType := rpc.RequestMessageType(method.Name)
	respType := rpc.ResponseMessageType(method.Name)
	w.WriteString(fmt.Sprintf("func (h *%s) %s(ctx context.Context, req *proto.%s) (*proto.%s, error) {\n",
		structName, method.Name, reqType, respType))

	w.WriteString(common.Tab + fmt.Sprintf("resp := proto.%s{}\n\n", respType))

	// hashID 转成 id
	if len(method.Parameters) > 0 && method.Parameters[0].Name == hashIDVarName {
		w.WriteString(common.Tab + idVarName + ", err := utils.DecodeHashID(req.HashID)\n")
		w.WriteString(common.Tab + "if err != nil {\n")
		w.WriteString(common.Tab + common.Tab + "resp.ErrCode = model.GetErrorCode(model.InvalidParameterError)\n")
		w.WriteString(common.Tab + common.Tab + "return &resp, nil\n")
		w.WriteString(common.Tab + "}\n\n")
	}

	w.WriteString(common.Tab)
	writeMethodReturns(w, method.Returns)
	w.WriteString(fmt.Sprintf(" := h.ucase.%s(", method.Name))
	writeUseCaseParams(w, method.Parameters)
	w.WriteString(")\n")

	w.WriteString(common.Tab)
	w.WriteString(fmt.Sprintf("utils.LogIfInnerError(err, \"%s\", ", method.Name))
	writeUseCaseParams(w, method.Parameters)
	w.WriteString(")\n\n")

	w.WriteString(common.Tab + "resp.ErrCode = model.GetErrorCode(err)\n")

	if len(method.Returns) > 1 {
		w.WriteString(common.Tab + "if err == nil {\n")
		writeResponseAssign(w, method.Returns[:len(method.Returns)-1])
		w.WriteString(common.Tab + "}\n")
	}

	w.WriteString("\n")
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

	"%[1]s/model"
	"%[1]s/model/usecase"
	"%[1]s/proto"

)

func New%[2]sHandler(ucase usecase.%[2]sUseCase) proto.%[2]sServer {
	return &%[3]s{ucase}
}

type %[3]s struct {
	ucase usecase.%[2]sUseCase
}
`, common.ProjectImportPrefix, serviceNameUpper, implStructName))

	for _, method := range parsedInterface.Methods {
		w.WriteString("\n")
		writeMethod(w, implStructName, method)
	}

	w.Flush()
}
