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

func genMain(dstFilePath string, parsedInterface *rpc.InterFace, dstName string) {
	serviceName = dstName
	serviceNameUpper = strings.Title(serviceName)
	//implStructName := serviceName + "Handler"

	f := common.CreateFile(dstFilePath)
	defer f.Close()

	w := bufio.NewWriter(f)

	w.WriteString(fmt.Sprintf(`package main

import (
	"net"
	"os"
	"strconv"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo"
	echoMiddleware "github.com/labstack/echo/middleware"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"

	httpHandler "%[1]s/handler/http"
	"%[1]s/handler/rpc"
	"%[1]s/model/repository"
	"%[1]s/model/usecase"
	"%[1]s/proto"

)

const driverName = "mysql"

var (
	devMode = os.Getenv("DEV") == "true"

	dataPath = utils.GetEnvWithDefault("DATA_PATH", "")

	httpPort = 80
	gRPCPort = 8001

	db *sqlx.DB
)

func init() {
	if devMode {
		httpPort = 10080
		gRPCPort = 18001
	}

	log.Info("连接数据库 ... ")
	mysqlHost := utils.GetEnvWithDefault("MYSQL_HOST", "192.168.2.198")
	mysqlPort := utils.GetEnvWithDefault("MYSQL_PORT", "3306")
	mysqlUser := utils.GetEnvWithDefault("MYSQL_USER", "") // TODO
	mysqlPass := utils.GetEnvWithDefault("MYSQL_PASS", "") // TODO
	mysqlDBName := utils.GetEnvWithDefault("MYSQL_DBNAME", "") // TODO
	initDB(mysqlHost, mysqlPort, mysqlUser, mysqlPass, mysqlDBName)
	log.Infoln("OK")

	setUpLogger()
}

func initDB(host, port, user, password, dbName string) {
	var err error
	connStr := user + ":" + password + "@tcp(" + host + ":" + port + ")/" + dbName + "?parseTime=true&loc=Local"
	db, err = sqlx.Open(driverName, connStr) // don't shadow outside db
	if err != nil {
		log.WithError(err).Fatalf("could not get a connection (%%s)", connStr)
	}
	if err := db.Ping(); err != nil {
		db.Close()
		log.WithError(err).Fatalf("could not establish a good connection (%%s)", connStr)
	}
	db.SetConnMaxLifetime(10 * time.Minute) // close expired connections
}

func setUpLogger() {
	file, err := os.OpenFile(dataPath+"main.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if err == nil {
		log.SetOutput(file)
	} else {
		log.WithError(err).Fatalln("Failed to log to file, using default stderr")
	}

	log.SetFormatter(&log.JSONFormatter{})
}

func newRouter(ucase usecase.%[2]s) *echo.Echo {
	e := echo.New()
	e.Use(echoMiddleware.Recover())

	accessLogFile, err := os.OpenFile(dataPath+"access.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		log.WithError(err).Fatalln("Failed create access log file")
	}
	e.Use(echoMiddleware.LoggerWithConfig(echoMiddleware.LoggerConfig{Output: accessLogFile}))

	httpHandler.Set%[3]sHTTPHandler(e, ucase)
	return e
}

func runWebService(ucase usecase.%[2]s) {
	r := newRouter(ucase)
	go func() {
		if err := r.Start(":" + strconv.Itoa(httpPort)); err != nil {
			log.WithError(err).Fatal("Web 服务运行失败！")
		}
	}()
}

func runGRPCService(ucase usecase.%[2]s) {
	gRPCServer := grpc.NewServer(grpc.KeepaliveParams(keepalive.ServerParameters{
		Time: 10 * time.Second,
	}))

	proto.Register%[3]sServer(gRPCServer, rpc.New%[3]sHandler(ucase))

	listener, err := net.Listen("tcp", ":"+strconv.Itoa(gRPCPort))
	if err != nil {
		log.WithError(err).Fatalln("gRPC failed to listen", "tcp", ":"+strconv.Itoa(gRPCPort))
	}

	if err := gRPCServer.Serve(listener); err != nil {
		log.WithError(err).Fatalln("gRPC 服务运行失败！")
	}
}

func main() {
	// TODO: new repositories and usecase

	runWebService(ucase)
	runGRPCService(ucase)
}
`, common.ProjectImportPrefix, strings.Title(parsedInterface.Name), serviceNameUpper))

	w.Flush()
}
