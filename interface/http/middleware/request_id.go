package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

const (
	// RequestIDKey 是请求 ID 在 gin.Context 中的键名
	RequestIDKey = "X-Request-ID"
)

// RequestID 返回一个中间件，为每个请求生成唯一的请求 ID
func RequestID(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 优先使用请求头中的 X-Request-ID
		requestID := c.GetHeader(RequestIDKey)
		if requestID == "" {
			// 如果没有，则生成一个新的 UUID
			requestID = uuid.New().String()
		}

		// 将请求 ID 设置到请求头中
		c.Header(RequestIDKey, requestID)

		// 将请求 ID 设置到上下文中，方便后续使用
		c.Set(RequestIDKey, requestID)

		// 将请求 ID 添加到日志中
		logger = logger.With(zap.String("request_id", requestID))
		c.Set("logger", logger)

		c.Next()
	}
}
