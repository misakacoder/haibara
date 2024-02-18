package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/misakacoder/logger"
	"haibara/config"
	"haibara/middleware"
	"haibara/model"
	"haibara/response"
	"haibara/router/auth"
	"haibara/router/log"
	"haibara/router/user"
	"haibara/util"
	"strings"
	"time"
)

var startTime = time.Now()

func init() {
	gin.SetMode(gin.ReleaseMode)
	initLogger()
}

func main() {
	router := gin.New()
	initTable()
	initRouter(router)
	start(router)
}

func initLogger() {
	logConfig := config.Configuration.Log
	logger.SetLogger(logger.NewSimpleLogger(logConfig.Filename))
	level, _ := logger.Parse(logConfig.Level)
	logger.SetLevel(level)
}

func initTable() {
	model.FirstOrCreate()
}

func initRouter(router *gin.Engine) {
	router.Use(middleware.NetWork)
	router.Use(middleware.Recovery)
	router.Use(middleware.CSRF)
	auth.LoadRouter(router)
	log.LoadRouter(router)
	user.LoadRouter(router)
	router.NoRoute(response.NotFound)
}

func start(router *gin.Engine) {
	port := config.Configuration.Server.Port
	banner := strings.Builder{}
	startUpTime := time.Since(startTime)
	banner.WriteString(fmt.Sprintf("Started haibara in %.2f seconds...", startUpTime.Seconds()))
	addresses := util.GetLocalAddr()
	for _, address := range addresses {
		banner.WriteString(fmt.Sprintf("\n - Listen on: http://%s:%d", address, port))
	}
	logger.Info(banner.String())
	err := router.Run(fmt.Sprintf(":%d", port))
	if err != nil {
		logger.Error(err.Error())
	}
}
