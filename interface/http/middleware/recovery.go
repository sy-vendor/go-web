package middleware

import (
	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// Recovery returns a middleware that recovers from any panics and writes a 500 error if any was caught.
func Recovery(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// 获取请求信息
				path := c.Request.URL.Path
				method := c.Request.Method
				clientIP := c.ClientIP()
				requestID := c.GetString(RequestIDKey)

				// 获取带有请求 ID 的 logger
				reqLogger := logger
				if l, exists := c.Get("logger"); exists {
					if log, ok := l.(*zap.Logger); ok {
						reqLogger = log
					}
				}

				// 记录详细的错误信息
				reqLogger.Error("panic recovered",
					zap.Any("error", err),
					zap.String("path", path),
					zap.String("method", method),
					zap.String("client_ip", clientIP),
					zap.String("request_id", requestID),
					zap.String("stack", string(debug.Stack())),
				)

				// 返回 500 错误
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
					"code":       http.StatusInternalServerError,
					"message":    fmt.Sprintf("Internal Server Error: %v", err),
					"request_id": requestID,
				})
			}
		}()

		c.Next()
	}
}
