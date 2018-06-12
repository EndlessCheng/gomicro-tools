package rpc

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"path/filepath"
)

var tab = strings.Repeat(" ", 4)

func mapGoTypeToProtoType(goTypeName string) string {
	switch goTypeName {
	case "int":
		return "int64"
	case "uint":
		return "uint64"
	case "error":
		return "int32"
	default:
		return goTypeName
	}
}

func mapGoNameToProtoName(goName string) string {
	switch goName {
	case "err":
		return "errCode"
	default:
		return goName
	}
}

func writeService(w *bufio.Writer, serviceName string, methods []*Method) {
	w.WriteString(fmt.Sprintf("service %s {\n", serviceName))
	for _, method := range methods {
		w.WriteString(fmt.Sprintf("%srpc %[2]s (Req%[2]s) returns (Resp%[2]s);\n", tab, method.MethodName))
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
		writeMessage(w, "Req"+method.MethodName, method.Parameters)
		w.WriteString("\n")
		writeMessage(w, "Resp"+method.MethodName, method.Returns)
	}

	w.Flush()
}
