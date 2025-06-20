package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"go-web/pkg/i18n"

	"go-web/pkg/errors"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

// ValidatorConfig 验证器配置
type ValidatorConfig struct {
	// 是否在验证失败时返回详细错误信息
	ShowDetailedErrors bool
	// 自定义错误消息
	CustomErrorMessages map[string]string
}

// DefaultValidatorConfig 返回默认的验证器配置
func DefaultValidatorConfig() *ValidatorConfig {
	return &ValidatorConfig{
		ShowDetailedErrors: true,
		CustomErrorMessages: map[string]string{
			"required": "字段 %s 是必填的",
			"email":    "字段 %s 必须是有效的邮箱地址",
			"min":      "字段 %s 的最小长度是 %s",
			"max":      "字段 %s 的最大长度是 %s",
		},
	}
}

// Validator 请求验证中间件
func Validator(logger *zap.Logger, config *ValidatorConfig) gin.HandlerFunc {
	validate := validator.New()

	// 注册自定义验证器
	registerCustomValidators(validate)

	return func(c *gin.Context) {
		// 获取请求体
		var body interface{}
		if err := c.ShouldBindJSON(&body); err != nil {
			if err != http.ErrNotSupported {
				logger.Warn("failed to bind request body",
					zap.Error(err),
					zap.String("path", c.Request.URL.Path),
				)
				i18n.ErrorResponse(c, errors.New(errors.ErrInvalidParam, "invalid_param", fmt.Sprintf("%v", err)))
				c.Abort()
				return
			}
		}

		// 验证请求体
		if body != nil {
			if err := validate.Struct(body); err != nil {
				handleValidationError(c, err, config, logger)
				return
			}
		}

		// 验证查询参数
		if err := validateQueryParams(c, validate, config); err != nil {
			handleValidationError(c, err, config, logger)
			return
		}

		c.Next()
	}
}

// registerCustomValidators 注册自定义验证器
func registerCustomValidators(validate *validator.Validate) {
	// 注册手机号验证器
	validate.RegisterValidation("phone", func(fl validator.FieldLevel) bool {
		value := fl.Field().String()
		// 简单的手机号验证，可以根据需要修改
		return len(value) == 11 && strings.HasPrefix(value, "1")
	})

	// 注册密码强度验证器
	validate.RegisterValidation("password", func(fl validator.FieldLevel) bool {
		value := fl.Field().String()
		// 密码至少包含8个字符，且包含大小写字母和数字
		return len(value) >= 8 &&
			strings.ContainsAny(value, "ABCDEFGHIJKLMNOPQRSTUVWXYZ") &&
			strings.ContainsAny(value, "abcdefghijklmnopqrstuvwxyz") &&
			strings.ContainsAny(value, "0123456789")
	})
}

// validateQueryParams 验证查询参数
func validateQueryParams(c *gin.Context, validate *validator.Validate, config *ValidatorConfig) error {
	query := c.Request.URL.Query()
	if len(query) == 0 {
		return nil
	}

	// 将查询参数转换为结构体
	queryStruct := make(map[string]interface{})
	for key, values := range query {
		if len(values) > 0 {
			queryStruct[key] = values[0]
		}
	}

	// 验证查询参数
	return validate.Struct(queryStruct)
}

// handleValidationError 处理验证错误
func handleValidationError(c *gin.Context, err error, config *ValidatorConfig, logger *zap.Logger) {
	if config.ShowDetailedErrors {
		// 获取详细的验证错误信息
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			var errorMessages []string
			for _, e := range validationErrors {
				errorMessages = append(errorMessages, formatValidationError(e))
			}
			i18n.ErrorResponse(c, errors.New(errors.ErrInvalidParam, "invalid_param", strings.Join(errorMessages, "; ")))
		} else {
			i18n.ErrorResponse(c, errors.New(errors.ErrInvalidParam, "invalid_param", err.Error()))
		}
		c.Abort()
	} else {
		// 只返回简单错误信息
		logger.Warn("validation failed",
			zap.Error(err),
			zap.String("path", c.Request.URL.Path),
		)

		i18n.ErrorResponse(c, errors.New(errors.ErrInvalidParam, "invalid_param", "请求参数验证失败"))
		c.Abort()
	}
}

// formatValidationError 格式化验证错误信息
func formatValidationError(err validator.FieldError) string {
	field := err.Field()
	switch err.Tag() {
	case "required":
		return fmt.Sprintf("%s is required", field)
	case "email":
		return fmt.Sprintf("%s must be a valid email address", field)
	case "min":
		return fmt.Sprintf("%s must be at least %s characters", field, err.Param())
	case "max":
		return fmt.Sprintf("%s must be at most %s characters", field, err.Param())
	default:
		return fmt.Sprintf("%s failed validation: %s", field, err.Tag())
	}
}
