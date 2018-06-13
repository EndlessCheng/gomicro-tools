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
		reqMsgType := requestMessageType(method.Name)
		respMsgType := responseMessageType(method.Name)
		w.WriteString(fmt.Sprintf("%srpc %s (%s) returns (%s);\n", tab, method.Name, reqMsgType, respMsgType))
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

func createFile(filePath string) *os.File {
	err := os.MkdirAll(path.Dir(filePath), os.ModePerm)
	check(err)

	f, err := os.Create(filePath)
	check(err)

	return f
}

func genProto(protoFilePath string, parsedMethods []*Method) {
	f := createFile(protoFilePath)
	defer f.Close()

	w := bufio.NewWriter(f)
	_, err := w.WriteString("syntax = \"proto3\";\n\npackage proto;\n\n")
	check(err)

	serviceName := filepath.Base(protoFilePath)
	serviceName = strings.TrimSuffix(serviceName, filepath.Ext(serviceName))
	serviceName = strings.Title(serviceName)
	writeService(w, serviceName, parsedMethods)

	for _, method := range parsedMethods {
		w.WriteString("\n")
		writeMessage(w, requestMessageType(method.Name), method.Parameters)
		w.WriteString("\n")
		writeMessage(w, responseMessageType(method.Name), method.Returns)
	}

	err = w.Flush()
	check(err)
}

func genMessages(protoFilePath string, parsedStructs []*Struct) {
	f := createFile(protoFilePath)
	defer f.Close()

	w := bufio.NewWriter(f)
	for _, parsedStruct := range parsedStructs {
		writeMessage(w, parsedStruct.Name, parsedStruct.Members)
		w.WriteString("\n")
	}

	err := w.Flush()
	check(err)
}
