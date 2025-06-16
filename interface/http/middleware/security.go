package middleware

import (
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

// SecurityConfig 是安全中间件的配置选项
type SecurityConfig struct {
	// XSSProtection 是否启用 XSS 防护
	XSSProtection bool
	// ContentTypeNosniff 是否启用内容类型嗅探防护
	ContentTypeNosniff bool
	// FrameGuard 是否启用点击劫持防护
	FrameGuard bool
	// HSTS 是否启用 HSTS
	HSTS bool
	// HSTSDuration HSTS 的持续时间（秒）
	HSTSDuration int
}

// DefaultSecurityConfig 返回默认的安全配置
func DefaultSecurityConfig() *SecurityConfig {
	return &SecurityConfig{
		XSSProtection:      true,
		ContentTypeNosniff: true,
		FrameGuard:         true,
		HSTS:               true,
		HSTSDuration:       31536000, // 1 year
	}
}

// Security 返回一个安全中间件
func Security(config *SecurityConfig) gin.HandlerFunc {
	if config == nil {
		config = DefaultSecurityConfig()
	}

	return func(c *gin.Context) {
		// 设置安全相关的响应头
		if config.XSSProtection {
			c.Header("X-XSS-Protection", "1; mode=block")
		}

		if config.ContentTypeNosniff {
			c.Header("X-Content-Type-Options", "nosniff")
		}

		if config.FrameGuard {
			c.Header("X-Frame-Options", "SAMEORIGIN")
		}

		if config.HSTS {
			c.Header("Strict-Transport-Security",
				strings.Join([]string{
					"max-age=" + strconv.Itoa(config.HSTSDuration),
					"includeSubDomains",
					"preload",
				}, "; "))
		}

		// 设置其他安全相关的响应头
		c.Header("X-Download-Options", "noopen")
		c.Header("X-Permitted-Cross-Domain-Policies", "none")
		c.Header("Referrer-Policy", "strict-origin-when-cross-origin")
		c.Header("Permissions-Policy", "geolocation=(), microphone=(), camera=()")

		// 移除敏感响应头
		c.Header("Server", "")
		c.Header("X-Powered-By", "")

		c.Next()
	}
}
