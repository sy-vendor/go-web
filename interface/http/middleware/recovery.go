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

				// 记录详细的错误信息
				logger.Error("panic recovered",
					zap.Any("error", err),
					zap.String("path", path),
					zap.String("method", method),
					zap.String("client_ip", clientIP),
					zap.String("stack", string(debug.Stack())),
				)

				// 返回 500 错误
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
					"code":    http.StatusInternalServerError,
					"message": fmt.Sprintf("Internal Server Error: %v", err),
				})
			}
		}()

		c.Next()
	}
}
