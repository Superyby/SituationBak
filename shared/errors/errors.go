package errors

import (
	"fmt"
	"net/http"
)

// 错误码常量
const (
	CodeOK                 = 0
	CodeSuccess            = 0
	CodeBadRequest         = 400
	CodeUnauthorized       = 401
	CodeForbidden          = 403
	CodeNotFound           = 404
	CodeConflict           = 409
	CodeInternalError      = 500
	CodeServiceUnavailable = 503

	// 业务错误码
	CodeInvalidParams    = 1000
	CodeUsernameExist    = 1001
	CodeEmailExists      = 1002
	CodeLoginFailed      = 1003
	CodePasswordWrong    = 1004
	CodeTokenInvalid     = 1005
	CodeTokenExpired     = 1006
	CodeUserNotFound     = 1007
	CodePermissionDenied = 1008
	CodeAlreadyExists    = 1009
)

// AppError 应用错误
type AppError struct {
	Code     int    `json:"code"`
	Message  string `json:"message"`
	HTTPCode int    `json:"-"`
	Err      error  `json:"-"`
}

// Error 实现error接口
func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("[%d] %s: %v", e.Code, e.Message, e.Err)
	}
	return fmt.Sprintf("[%d] %s", e.Code, e.Message)
}

// Unwrap 返回底层错误
func (e *AppError) Unwrap() error {
	return e.Err
}

// New 创建新错误
func New(code int, message string) *AppError {
	return &AppError{
		Code:     code,
		Message:  message,
		HTTPCode: codeToHTTP(code),
	}
}

// Wrap 包装错误
func Wrap(code int, message string, err error) *AppError {
	return &AppError{
		Code:     code,
		Message:  message,
		HTTPCode: codeToHTTP(code),
		Err:      err,
	}
}

// WithCode 根据错误码创建错误
func WithCode(code int) *AppError {
	return New(code, codeToMessage(code))
}

// ErrInternal 内部错误
func ErrInternal(err error) *AppError {
	return Wrap(CodeInternalError, "内部错误", err)
}

// ErrBadRequest 错误请求
func ErrBadRequest(message string) *AppError {
	return New(CodeBadRequest, message)
}

// ErrUnauthorized 未授权
func ErrUnauthorized() *AppError {
	return New(CodeUnauthorized, "未授权")
}

// ErrForbidden 禁止访问
func ErrForbidden() *AppError {
	return New(CodeForbidden, "禁止访问")
}

// ErrNotFound 未找到
func ErrNotFound(resource string) *AppError {
	return New(CodeNotFound, fmt.Sprintf("%s未找到", resource))
}

// ErrUserNotFound 用户未找到
func ErrUserNotFound() *AppError {
	return New(CodeUserNotFound, "用户不存在")
}

// ErrTokenInvalid Token无效
func ErrTokenInvalid() *AppError {
	return New(CodeTokenInvalid, "Token无效")
}

// ErrTokenExpired Token过期
func ErrTokenExpired() *AppError {
	return New(CodeTokenExpired, "Token已过期")
}

// codeToMessage 错误码转消息
func codeToMessage(code int) string {
	messages := map[int]string{
		CodeOK:                 "成功",
		CodeBadRequest:         "请求参数错误",
		CodeUnauthorized:       "未授权",
		CodeForbidden:          "禁止访问",
		CodeNotFound:           "资源未找到",
		CodeConflict:           "资源冲突",
		CodeInternalError:      "内部错误",
		CodeServiceUnavailable: "服务不可用",
		CodeInvalidParams:      "无效的参数",
		CodeUsernameExist:      "用户名已存在",
		CodeEmailExists:        "邮箱已存在",
		CodeLoginFailed:        "用户名或密码错误",
		CodePasswordWrong:      "密码错误",
		CodeTokenInvalid:       "Token无效",
		CodeTokenExpired:       "Token已过期",
		CodeUserNotFound:       "用户不存在",
		CodePermissionDenied:   "权限不足",
		CodeAlreadyExists:      "资源已存在",
	}

	if msg, ok := messages[code]; ok {
		return msg
	}
	return "未知错误"
}

// GetMessage 获取错误消息
func GetMessage(code int) string {
	return codeToMessage(code)
}

// GetHTTPStatus 获取HTTP状态码
func GetHTTPStatus(code int) int {
	return codeToHTTP(code)
}

// AsAppError 将error转换为AppError
func AsAppError(err error) (*AppError, bool) {
	if appErr, ok := err.(*AppError); ok {
		return appErr, true
	}
	return nil, false
}

// codeToHTTP 错误码转HTTP状态码
func codeToHTTP(code int) int {
	switch {
	case code >= 1000 && code < 2000:
		return http.StatusBadRequest
	case code == CodeUnauthorized || code == CodeTokenInvalid || code == CodeTokenExpired:
		return http.StatusUnauthorized
	case code == CodeForbidden || code == CodePermissionDenied:
		return http.StatusForbidden
	case code == CodeNotFound || code == CodeUserNotFound:
		return http.StatusNotFound
	case code == CodeConflict:
		return http.StatusConflict
	case code >= 500:
		return http.StatusInternalServerError
	default:
		return code
	}
}
