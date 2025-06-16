package middleware

import (
	"bytes"
	"io"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// bodyLogWriter 是一个自定义的响应写入器，用于捕获响应体
type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

// Write 重写 Write 方法，同时写入响应体和缓冲区
func (w bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

// Logger 返回一个日志中间件，记录请求和响应的详细信息
func Logger(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 开始时间
		start := time.Now()

		// 获取请求体
		var requestBody []byte
		if c.Request.Body != nil {
			requestBody, _ = io.ReadAll(c.Request.Body)
			// 恢复请求体，因为读取后需要重新设置
			c.Request.Body = io.NopCloser(bytes.NewBuffer(requestBody))
		}

		// 创建自定义响应写入器
		blw := &bodyLogWriter{
			ResponseWriter: c.Writer,
			body:           bytes.NewBufferString(""),
		}
		c.Writer = blw

		// 处理请求
		c.Next()

		// 结束时间
		end := time.Now()
		latency := end.Sub(start)

		// 获取请求 ID
		requestID := c.GetString(RequestIDKey)

		// 获取带有请求 ID 的 logger
		reqLogger := logger
		if l, exists := c.Get("logger"); exists {
			if log, ok := l.(*zap.Logger); ok {
				reqLogger = log
			}
		}

		// 记录请求和响应信息
		reqLogger.Info("request completed",
			zap.String("request_id", requestID),
			zap.String("method", c.Request.Method),
			zap.String("path", c.Request.URL.Path),
			zap.String("query", c.Request.URL.RawQuery),
			zap.String("ip", c.ClientIP()),
			zap.String("user_agent", c.Request.UserAgent()),
			zap.Int("status", c.Writer.Status()),
			zap.Int("size", c.Writer.Size()),
			zap.Duration("latency", latency),
			zap.ByteString("request_body", requestBody),
			zap.String("response_body", blw.body.String()),
			zap.Strings("errors", c.Errors.Errors()),
		)
	}
}
