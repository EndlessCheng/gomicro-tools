package rpc

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"path/filepath"
	"path"
)

func writeService(w *bufio.Writer, serviceName string, methods []*Method) {
	w.WriteString(fmt.Sprintf("service %s {\n", serviceName))
	for _, method := range methods {
		reqMsgType := requestMessageType(method.MethodName)
		respMsgType := responseMessageType(method.MethodName)
		w.WriteString(fmt.Sprintf("%srpc %s (%s) returns (%s);\n", tab, method.MethodName, reqMsgType, respMsgType))
	}
	w.WriteString("}\n")
}

func writeMessage(w *bufio.Writer, messageType string, parameters []*Var) {
	w.WriteString(fmt.Sprintf("message %s {\n", messageType))
	for i, parameter := range parameters {
		w.WriteString(tab)
		if parameter.IsSlice {
			w.WriteString("repeated ")
		}
		w.WriteString(fmt.Sprintf("%s %s = %d;\n", mapGoTypeToProtoType(parameter.Type), mapGoNameToProtoName(parameter.Name), i+1))
	}
	w.WriteString("}\n")
}

func genRpc(rpcFilePath string, parsedMethods []*Method) {
	err := os.MkdirAll(path.Dir(rpcFilePath), os.ModePerm)
	check(err)

	f, err := os.Create(rpcFilePath)
	check(err)
	defer f.Close()

	w := bufio.NewWriter(f)
	_, err = w.WriteString("syntax = \"proto3\";\n\npackage proto;\n\n")
	check(err)

	serviceName := filepath.Base(rpcFilePath)
	serviceName = strings.TrimSuffix(serviceName, filepath.Ext(serviceName))
	serviceName = strings.Title(serviceName)
	writeService(w, serviceName, parsedMethods)

	for _, method := range parsedMethods {
		w.WriteString("\n")
		writeMessage(w, requestMessageType(method.MethodName), method.Parameters)
		w.WriteString("\n")
		writeMessage(w, responseMessageType(method.MethodName), method.Returns)
	}

	w.Flush()
}
