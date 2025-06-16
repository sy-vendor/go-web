package i18n

import (
	"strings"

	"go-web/pkg/response"

	"github.com/gin-gonic/gin"
)

var messages = map[string]map[string]string{
	"zh": {
		"success":        "成功",
		"invalid_param":  "参数无效",
		"internal_error": "服务器内部错误",
	},
	"en": {
		"success":        "Success",
		"invalid_param":  "Invalid parameter",
		"internal_error": "Internal server error",
	},
}

// GetLang 获取请求语言
func GetLang(c *gin.Context) string {
	lang := c.GetHeader("Accept-Language")
	if strings.HasPrefix(lang, "zh") {
		return "zh"
	}
	return "en"
}

// T 返回指定 key 的多语言消息
func T(c *gin.Context, key string) string {
	lang := GetLang(c)
	if msg, ok := messages[lang][key]; ok {
		return msg
	}
	// fallback
	if msg, ok := messages["en"][key]; ok {
		return msg
	}
	return key
}

// ErrorResponse 返回国际化错误响应
func ErrorResponse(c *gin.Context, code int, key string, details interface{}) {
	response.Error(c, code, T(c, key), details)
}

// SuccessResponse 返回国际化成功响应
func SuccessResponse(c *gin.Context, data interface{}) {
	response.Success(c, data)
}
