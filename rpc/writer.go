package rpc

import (
	"bufio"
	"fmt"
	"gomicro-tools/common"
	"path/filepath"
	"strings"
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
	w.WriteString(fmt.Sprintf("message %s {\n", strings.Title(messageType)))
	for i, parameter := range parameters {
		w.WriteString(common.Tab)
		if parameter.IsSlice {
			w.WriteString("repeated ")
		}
		index := i + 1
		if parameter.Name == "err" { // TODO: 重构一下
			index = 60
		}
		w.WriteString(fmt.Sprintf("%s %s = %d;\n", MapGoTypeToProtoType(parameter.Type), MapGoNameToProtoName(parameter.Name), index))
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

func genFullProto(protoFilePath string, parsedMethods []*Method, parsedStructs []*Struct) {
	f := common.CreateFile(protoFilePath)
	defer f.Close()

	w := bufio.NewWriter(f)
	_, err := w.WriteString("syntax = \"proto3\";\n\npackage proto;\n\n")
	common.Check(err)

	serviceName := filepath.Base(protoFilePath)
	serviceName = strings.TrimSuffix(serviceName, filepath.Ext(serviceName))
	serviceName = strings.Title(serviceName)
	writeService(w, serviceName, parsedMethods)

	for _, parsedStruct := range parsedStructs {
		w.WriteString("\n")
		writeMessage(w, parsedStruct.Name, parsedStruct.Members)
	}

	for _, method := range parsedMethods {
		w.WriteString("\n")
		writeMessage(w, RequestMessageType(method.Name), method.Parameters)
		w.WriteString("\n")
		writeMessage(w, ResponseMessageType(method.Name), method.Returns)
	}

	err = w.Flush()
	common.Check(err)
}
