package middleware

import (
	"crypto/rand"
	"encoding/base64"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// CSRFConfig 是 CSRF 中间件的配置选项
type CSRFConfig struct {
	// CookieName CSRF token 的 cookie 名称
	CookieName string
	// HeaderName CSRF token 的请求头名称
	HeaderName string
	// CookiePath CSRF token cookie 的路径
	CookiePath string
	// Secure 是否只在 HTTPS 连接中发送 cookie
	Secure bool
	// SameSite cookie 的 SameSite 属性
	SameSite http.SameSite
	// SkipPaths 不需要 CSRF 防护的路径列表
	SkipPaths []string
}

// DefaultCSRFConfig 返回默认的 CSRF 配置
func DefaultCSRFConfig() *CSRFConfig {
	return &CSRFConfig{
		CookieName: "csrf_token",
		HeaderName: "X-CSRF-Token",
		CookiePath: "/",
		Secure:     true,
		SameSite:   http.SameSiteStrictMode,
		SkipPaths: []string{
			"/health",
			"/metrics",
		},
	}
}

// generateToken 生成一个随机的 CSRF token
func generateToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(b), nil
}

// shouldSkip 检查路径是否需要跳过 CSRF 防护
func shouldSkip(path string, skipPaths []string) bool {
	for _, skipPath := range skipPaths {
		if path == skipPath {
			return true
		}
	}
	return false
}

// CSRF 返回一个 CSRF 防护中间件
func CSRF(logger *zap.Logger, config *CSRFConfig) gin.HandlerFunc {
	if config == nil {
		config = DefaultCSRFConfig()
	}

	return func(c *gin.Context) {
		// 检查是否需要跳过 CSRF 防护
		if shouldSkip(c.Request.URL.Path, config.SkipPaths) {
			c.Next()
			return
		}

		// 只对非 GET 请求进行 CSRF 检查
		if c.Request.Method != "GET" {
			// 获取请求头中的 token
			token := c.GetHeader(config.HeaderName)
			if token == "" {
				// 尝试从 cookie 中获取 token
				cookie, err := c.Cookie(config.CookieName)
				if err != nil || cookie == "" {
					logger.Warn("CSRF token missing",
						zap.String("path", c.Request.URL.Path),
						zap.String("method", c.Request.Method),
						zap.String("ip", c.ClientIP()),
					)
					c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
						"code":    http.StatusForbidden,
						"message": "CSRF token missing",
					})
					return
				}
				token = cookie
			}

			// 验证 token
			cookie, err := c.Cookie(config.CookieName)
			if err != nil || cookie != token {
				logger.Warn("CSRF token mismatch",
					zap.String("path", c.Request.URL.Path),
					zap.String("method", c.Request.Method),
					zap.String("ip", c.ClientIP()),
				)
				c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
					"code":    http.StatusForbidden,
					"message": "CSRF token mismatch",
				})
				return
			}
		}

		// 生成新的 token
		token, err := generateToken()
		if err != nil {
			logger.Error("failed to generate CSRF token",
				zap.Error(err),
				zap.String("path", c.Request.URL.Path),
				zap.String("method", c.Request.Method),
				zap.String("ip", c.ClientIP()),
			)
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"code":    http.StatusInternalServerError,
				"message": "Internal Server Error",
			})
			return
		}

		// 设置 cookie
		c.SetCookie(
			config.CookieName,
			token,
			0, // 会话 cookie
			config.CookiePath,
			"",
			config.Secure,
			true, // HttpOnly
		)

		// 将 token 添加到响应头中
		c.Header(config.HeaderName, token)

		c.Next()
	}
}
