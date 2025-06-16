package i18n

import (
	"net/http"
	"strings"

	"go-web/pkg/errors"
	"go-web/pkg/response"

	"github.com/gin-gonic/gin"
)

var messages = map[string]map[string]string{
	"zh": {
		// 系统级错误
		"success":       "成功",
		"system_error":  "系统错误",
		"unknown_error": "未知错误",

		// 参数验证错误
		"invalid_param":  "参数无效",
		"missing_param":  "缺少必要参数",
		"invalid_format": "参数格式错误",

		// 认证授权错误
		"unauthorized":  "未授权访问",
		"forbidden":     "禁止访问",
		"token_expired": "登录已过期",

		// 业务逻辑错误
		"not_found":     "资源不存在",
		"already_exist": "资源已存在",
		"invalid_state": "状态无效",
	},
	"en": {
		// System errors
		"success":       "Success",
		"system_error":  "System Error",
		"unknown_error": "Unknown Error",

		// Parameter validation errors
		"invalid_param":  "Invalid Parameter",
		"missing_param":  "Missing Required Parameter",
		"invalid_format": "Invalid Format",

		// Authentication errors
		"unauthorized":  "Unauthorized",
		"forbidden":     "Forbidden",
		"token_expired": "Token Expired",

		// Business logic errors
		"not_found":     "Resource Not Found",
		"already_exist": "Resource Already Exists",
		"invalid_state": "Invalid State",
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

// TByLang 通过 lang 获取多语言消息
func TByLang(lang, key string) string {
	if msg, ok := messages[lang][key]; ok {
		return msg
	}
	if msg, ok := messages["en"][key]; ok {
		return msg
	}
	return key
}

// ErrorResponse 返回国际化错误响应
func ErrorResponse(c *gin.Context, err error) {
	if e, ok := err.(*errors.Error); ok {
		response.Error(c, e.HTTPStatus(), T(c, e.Message), e.Details)
		return
	}
	response.Error(c, http.StatusInternalServerError, T(c, "system_error"), err.Error())
}

// SuccessResponse 返回国际化成功响应
func SuccessResponse(c *gin.Context, data interface{}) {
	response.Success(c, data)
}
