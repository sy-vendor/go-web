package middleware

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
)

// CORSConfig 是 CORS 中间件的配置选项
type CORSConfig struct {
	// AllowedOrigins 允许的源列表
	AllowedOrigins []string
	// AllowedMethods 允许的 HTTP 方法列表
	AllowedMethods []string
	// AllowedHeaders 允许的请求头列表
	AllowedHeaders []string
	// ExposedHeaders 允许客户端访问的响应头列表
	ExposedHeaders []string
	// AllowCredentials 是否允许发送认证信息（cookies等）
	AllowCredentials bool
	// MaxAge 预检请求的缓存时间（秒）
	MaxAge time.Duration
}

// DefaultCORSConfig 返回默认的 CORS 配置
func DefaultCORSConfig() *CORSConfig {
	return &CORSConfig{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowedHeaders: []string{
			"Origin",
			"Content-Type",
			"Accept",
			"Authorization",
			"X-Requested-With",
			"X-Request-ID",
		},
		ExposedHeaders: []string{
			"Content-Length",
			"Content-Type",
			"X-Request-ID",
		},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}
}

// Cors 返回一个配置了 CORS 的中间件
func Cors(config *CORSConfig) gin.HandlerFunc {
	if config == nil {
		config = DefaultCORSConfig()
	}

	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")
		if origin != "" {
			// 检查源是否在允许列表中
			allowed := false
			for _, allowedOrigin := range config.AllowedOrigins {
				if allowedOrigin == "*" || allowedOrigin == origin {
					allowed = true
					break
				}
			}

			if allowed {
				c.Header("Access-Control-Allow-Origin", origin)
				c.Header("Access-Control-Allow-Methods", joinStrings(config.AllowedMethods, ", "))
				c.Header("Access-Control-Allow-Headers", joinStrings(config.AllowedHeaders, ", "))
				c.Header("Access-Control-Expose-Headers", joinStrings(config.ExposedHeaders, ", "))
				c.Header("Access-Control-Max-Age", fmt.Sprintf("%d", int(config.MaxAge.Seconds())))

				if config.AllowCredentials {
					c.Header("Access-Control-Allow-Credentials", "true")
				}

				// 处理预检请求
				if c.Request.Method == "OPTIONS" {
					c.AbortWithStatus(204)
					return
				}
			}
		}

		c.Next()
	}
}

// joinStrings 将字符串切片连接成单个字符串
func joinStrings(strs []string, sep string) string {
	if len(strs) == 0 {
		return ""
	}
	if len(strs) == 1 {
		return strs[0]
	}
	result := strs[0]
	for _, s := range strs[1:] {
		result += sep + s
	}
	return result
}
