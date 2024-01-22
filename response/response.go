package response

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

const (
	success = "ok"
	failure = "error"
)

type Response struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}

func Ok(context *gin.Context) {
	JSON(context, http.StatusOK, success, nil)
}

func OkWithMessageJSON(context *gin.Context, message string) {
	JSON(context, http.StatusOK, message, nil)
}

func OkWithDataJSON(context *gin.Context, data any) {
	JSON(context, http.StatusOK, success, data)
}

func ErrorWithMessageJSON(context *gin.Context, message string) {
	JSON(context, http.StatusInternalServerError, message, nil)
}

func ErrorWithDataJSON(context *gin.Context, data any) {
	JSON(context, http.StatusInternalServerError, failure, data)
}

func AbortWithMessageJSON(context *gin.Context, message string) {
	AbortJSON(context, http.StatusInternalServerError, message, nil)
}

func NotFound(context *gin.Context) {
	JSON(context, http.StatusNotFound, "not found 404", nil)
}

func JSON(context *gin.Context, code int, message string, data any) {
	context.JSON(code, Response{Code: code, Message: message, Data: data})
}

func AbortJSON(context *gin.Context, code int, message string, data any) {
	JSON(context, code, message, data)
	context.Abort()
}
