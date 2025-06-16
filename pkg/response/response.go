package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Response 统一响应结构体
type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Details interface{} `json:"details,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

var defaultSuccessMsg = "success"

// Success 返回成功响应
func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code:    http.StatusOK,
		Message: defaultSuccessMsg,
		Data:    data,
	})
}

// Error 返回错误响应
func Error(c *gin.Context, code int, message string, details interface{}) {
	c.JSON(code, Response{
		Code:    code,
		Message: message,
		Details: details,
	})
}

// AbortWithError 终止请求并返回错误响应
func AbortWithError(c *gin.Context, code int, message string, details interface{}) {
	c.AbortWithStatusJSON(code, Response{
		Code:    code,
		Message: message,
		Details: details,
	})
}
