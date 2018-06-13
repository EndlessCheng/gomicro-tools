package rpc

import (
	"bufio"
	"fmt"
	"strings"
	"path/filepath"
	"gomicro-tools/common"
)

func writeService(w *bufio.Writer, serviceName string, methods []*Method) {
	w.WriteString(fmt.Sprintf("service %s {\n", serviceName))
	for _, method := range methods {
		reqMsgType := RequestMessageType(method.Name)
		respMsgType := ResponseMessageType(method.Name)
		w.WriteString(fmt.Sprintf("%srpc %s (%s) returns (%s);\n", common.Tab, method.Name, reqMsgType, respMsgType))
	}
	w.WriteString("}\n")
}

func writeMessage(w *bufio.Writer, messageType string, parameters []*Var) {
	w.WriteString(fmt.Sprintf("message %s {\n", messageType))
	for i, parameter := range parameters {
		w.WriteString(common.Tab)
		if parameter.IsSlice {
			w.WriteString("repeated ")
		}
		w.WriteString(fmt.Sprintf("%s %s = %d;\n", MapGoTypeToProtoType(parameter.Type), MapGoNameToProtoName(parameter.Name), i+1))
	}
	w.WriteString("}\n")
}

func genProto(protoFilePath string, parsedMethods []*Method) {
	f := common.CreateFile(protoFilePath)
	defer f.Close()

	w := bufio.NewWriter(f)
	_, err := w.WriteString("syntax = \"proto3\";\n\npackage proto;\n\n")
	common.Check(err)

	serviceName := filepath.Base(protoFilePath)
	serviceName = strings.TrimSuffix(serviceName, filepath.Ext(serviceName))
	serviceName = strings.Title(serviceName)
	writeService(w, serviceName, parsedMethods)

	for _, method := range parsedMethods {
		w.WriteString("\n")
		writeMessage(w, RequestMessageType(method.Name), method.Parameters)
		w.WriteString("\n")
		writeMessage(w, ResponseMessageType(method.Name), method.Returns)
	}

	err = w.Flush()
	common.Check(err)
}

func genMessages(protoFilePath string, parsedStructs []*Struct) {
	f := common.CreateFile(protoFilePath)
	defer f.Close()

	w := bufio.NewWriter(f)
	for _, parsedStruct := range parsedStructs {
		writeMessage(w, parsedStruct.Name, parsedStruct.Members)
		w.WriteString("\n")
	}

	err := w.Flush()
	common.Check(err)
}
