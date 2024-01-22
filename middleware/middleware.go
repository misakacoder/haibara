package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/misakacoder/logger"
	"haibara/consts"
	"haibara/definition"
	"haibara/response"
	"haibara/ro"
	"haibara/router/auth"
	"haibara/util"
	"net/http"
	"runtime"
	"strings"
	"time"
)

func NetWork(context *gin.Context) {
	start := time.Now()
	context.Next()
	duration := time.Since(start)
	request := context.Request
	ip := context.ClientIP()
	ip = util.ConditionalExpression(ip == "::1", consts.Localhost, ip)
	logger.Info("%s %s %s %d %dms", ip, request.Method, request.RequestURI, context.Writer.Status(), duration.Milliseconds())
}

func Recovery(context *gin.Context) {
	defer func() {
		if err := recover(); err != nil {
			logger.Error("%s", getStackTrace(err))
			var message string
			switch tp := err.(type) {
			case error:
				message = tp.Error()
			case string:
				message = tp
			default:
				message = fmt.Sprintf("%v", tp)
			}
			response.AbortWithMessageJSON(context, message)
		}
	}()
	context.Next()
}

func CSRF(context *gin.Context) {
	method := context.Request.Method
	context.Header("Access-Control-Allow-Origin", "*")
	context.Header("Access-Control-Allow-Methods", " GET, POST, PUT, DELETE, HEAD, OPTIONS")
	context.Header("Access-Control-Allow-Headers", "Content-Type, AccessToken, X-CSRF-Token, Authorization, Token")
	context.Header("Access-Control-Expose-Headers", "Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type, Content-Length")
	context.Header("Access-Control-Allow-Credentials", "true")
	if method == "OPTIONS" {
		context.AbortWithStatus(http.StatusMethodNotAllowed)
	}
	context.Next()
}

func Auth(context *gin.Context) {
	hasRole := HasRole()
	hasRole(context)
}

func HasRole(roles ...definition.Role) gin.HandlerFunc {
	return func(context *gin.Context) {
		token := util.RequireNonNullElse(context.Query("token"), context.GetHeader("token"))
		if token != "" {
			claims, err := auth.ParseTokenString(token)
			if err != nil {
				response.ErrorWithMessageJSON(context, err.Error())
				return
			}
			if claims != nil {
				user, err := ro.ValidateUserRO(&ro.UserRO{Username: claims.Username, Password: claims.Password})
				if err != nil {
					response.ErrorWithMessageJSON(context, err.Error())
					return
				}
				hashRole(context, user.Roles(), roles)
			}
			context.Next()
			return
		}
		if value, _ := context.Cookie(consts.LoginCookieKey); value != "" {
			userRO, err := util.ParseObject[ro.UserRO](util.DecAES(value))
			if err != nil {
				response.ErrorWithMessageJSON(context, err.Error())
				return
			}
			user, err := ro.ValidateUserRO(&userRO)
			if err != nil {
				response.ErrorWithMessageJSON(context, err.Error())
				return
			}
			if user != nil {
				hashRole(context, user.Roles(), roles)
			}
			context.Next()
			return
		}
		response.AbortJSON(context, http.StatusUnauthorized, "Unauthorized", nil)
	}
}

func hashRole(context *gin.Context, source []definition.Role, target []definition.Role) {
	if len(target) == 0 {
		return
	}
	var roleNames []string
	for _, targetRole := range target {
		roleNames = append(roleNames, targetRole.String())
		for _, sourceRole := range source {
			if sourceRole == targetRole {
				return
			}
		}
	}
	response.AbortJSON(context, http.StatusUnauthorized, fmt.Sprintf("访问失败，仅支持`%s`角色访问", strings.Join(roleNames, "、")), nil)
}

func getStackTrace(err any) string {
	stackTrace := strings.Builder{}
	stackTrace.WriteString(fmt.Sprintf("%v", err))
	for i := 1; ; i++ {
		pc, file, line, ok := runtime.Caller(i)
		if !ok {
			break
		}
		stackTrace.WriteString(fmt.Sprintf("\n - %s:%d (0x%x)", file, line, pc))
	}
	return stackTrace.String()
}
