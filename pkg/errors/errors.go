package errors

import (
	"fmt"
	"net/http"
)

// ErrorCode 错误码类型
type ErrorCode int

const (
	// 系统级错误码 (1-999)
	ErrSuccess ErrorCode = 0
	ErrSystem  ErrorCode = 100
	ErrUnknown ErrorCode = 999

	// 参数验证错误码 (1000-1999)
	ErrInvalidParam  ErrorCode = 1000
	ErrMissingParam  ErrorCode = 1001
	ErrInvalidFormat ErrorCode = 1002

	// 认证授权错误码 (2000-2999)
	ErrUnauthorized ErrorCode = 2000
	ErrForbidden    ErrorCode = 2001
	ErrTokenExpired ErrorCode = 2002

	// 业务逻辑错误码 (3000-3999)
	ErrNotFound     ErrorCode = 3000
	ErrAlreadyExist ErrorCode = 3001
	ErrInvalidState ErrorCode = 3002
)

// Error 自定义错误结构
type Error struct {
	Code    ErrorCode `json:"code"`
	Message string    `json:"message"`
	Details string    `json:"details,omitempty"`
}

func (e *Error) Error() string {
	if e.Details != "" {
		return fmt.Sprintf("[%d] %s: %s", e.Code, e.Message, e.Details)
	}
	return fmt.Sprintf("[%d] %s", e.Code, e.Message)
}

// New 创建新的错误
func New(code ErrorCode, message string, details ...string) *Error {
	err := &Error{
		Code:    code,
		Message: message,
	}
	if len(details) > 0 {
		err.Details = details[0]
	}
	return err
}

// HTTPStatus 获取错误对应的 HTTP 状态码
func (e *Error) HTTPStatus() int {
	switch {
	case e.Code >= 1000 && e.Code < 2000:
		return http.StatusBadRequest
	case e.Code >= 2000 && e.Code < 3000:
		return http.StatusUnauthorized
	case e.Code >= 3000 && e.Code < 4000:
		return http.StatusConflict
	default:
		return http.StatusInternalServerError
	}
}

// IsNotFound 判断是否为资源不存在错误
func (e *Error) IsNotFound() bool {
	return e.Code == ErrNotFound
}

// IsInvalidParam 判断是否为参数错误
func (e *Error) IsInvalidParam() bool {
	return e.Code >= 1000 && e.Code < 2000
}

// IsAuthError 判断是否为认证错误
func (e *Error) IsAuthError() bool {
	return e.Code >= 2000 && e.Code < 3000
}

// IsBusinessError 判断是否为业务错误
func (e *Error) IsBusinessError() bool {
	return e.Code >= 3000 && e.Code < 4000
}
