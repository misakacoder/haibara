package log

import (
	"github.com/gin-gonic/gin"
	"github.com/misakacoder/logger"
	"haibara/middleware"
	"haibara/response"
	"haibara/ro"
)

func LoadRouter(router *gin.Engine) {
	user := router.Group("/api/v1")
	{
		user.GET("/log", middleware.Auth, setLogLevel)
	}
}

func setLogLevel(context *gin.Context) {
	loggerRO := ro.ValidateRO(context, &ro.LoggerRO{})
	level, ok := logger.Parse(loggerRO.Level)
	if ok {
		logger.SetLevel(level)
		response.OkWithMessageJSON(context, "日志等级修改成功")
	} else {
		response.ErrorWithMessageJSON(context, "不支持此日志等级")
	}
}
