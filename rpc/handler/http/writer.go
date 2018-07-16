package handler

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

func writeMethod(w *bufio.Writer, method *rpc.Method) {
	w.WriteString(fmt.Sprintf(`func (h *handler) %s(c echo.Context) error {
	// TODO

}
`, common.LowerHead(method.Name)))
}

func genHandler(dstFilePath string, parsedInterface *rpc.InterFace, dstName string) {
	serviceName = dstName
	serviceNameUpper = strings.Title(serviceName)
	//implStructName := serviceName + "Handler"

	f := common.CreateFile(dstFilePath)
	defer f.Close()

	w := bufio.NewWriter(f)

	w.WriteString(fmt.Sprintf(`package http

import (
	"encoding/json"
	"net/http"

	"github.com/labstack/echo"

	"%[1]s/model/usecase"

)

type handler struct {
	ucase usecase.%[2]s
}

func Set%[3]sHTTPHandler(e *echo.Echo, ucase usecase.%[2]s) {
	h := handler{ucase}

	// TODO
}

func createAPIResultMap(err error) map[string]interface{} {
	msg := "OK"
	if err != nil {
		msg = err.Error()
	}

	return map[string]interface{}{
		"errcode": model.GetErrorCode(err),
		"msg":     msg,
	}
}
`, common.ProjectImportPrefix, strings.Title(parsedInterface.Name), serviceNameUpper))

	for _, method := range parsedInterface.Methods {
		w.WriteString("\n")
		writeMethod(w, method)
	}

	w.Flush()
}
