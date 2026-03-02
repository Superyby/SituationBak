package validator

import (
	"fmt"
	"regexp"
	"strings"
	"unicode/utf8"

	"SituationBak/shared/errors"

	"github.com/gofiber/fiber/v3"
)

// Validatable 可验证接口
// 实现此接口的结构体可以使用 BindAndValidate 进行自动验证
type Validatable interface {
	Validate() error
}

// BindAndValidate 绑定请求体并验证
func BindAndValidate(c fiber.Ctx, req Validatable) error {
	if err := c.Bind().Body(req); err != nil {
		return errors.New(errors.CodeInvalidParams, "请求参数格式错误")
	}

	if err := req.Validate(); err != nil {
		if appErr, ok := err.(*errors.AppError); ok {
			return appErr
		}
		return errors.New(errors.CodeInvalidParams, err.Error())
	}

	return nil
}

// BindQueryAndValidate 绑定查询参数并验证
func BindQueryAndValidate(c fiber.Ctx, req Validatable) error {
	if err := c.Bind().Query(req); err != nil {
		return errors.New(errors.CodeInvalidParams, "查询参数格式错误")
	}

	if err := req.Validate(); err != nil {
		if appErr, ok := err.(*errors.AppError); ok {
			return appErr
		}
		return errors.New(errors.CodeInvalidParams, err.Error())
	}

	return nil
}

// ==================== 验证规则 ====================

// ValidationError 验证错误
type ValidationError struct {
	Field   string
	Message string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("%s: %s", e.Field, e.Message)
}

// NewValidationError 创建验证错误
func NewValidationError(field, message string) *ValidationError {
	return &ValidationError{Field: field, Message: message}
}

// ==================== 字符串验证 ====================

// Required 必填验证
func Required(value, field string) error {
	if strings.TrimSpace(value) == "" {
		return NewValidationError(field, "不能为空")
	}
	return nil
}

// MinLength 最小长度验证
func MinLength(value, field string, min int) error {
	if utf8.RuneCountInString(value) < min {
		return NewValidationError(field, fmt.Sprintf("长度不能少于%d个字符", min))
	}
	return nil
}

// MaxLength 最大长度验证
func MaxLength(value, field string, max int) error {
	if utf8.RuneCountInString(value) > max {
		return NewValidationError(field, fmt.Sprintf("长度不能超过%d个字符", max))
	}
	return nil
}

// Length 长度范围验证
func Length(value, field string, min, max int) error {
	length := utf8.RuneCountInString(value)
	if length < min || length > max {
		return NewValidationError(field, fmt.Sprintf("长度应在%d-%d个字符之间", min, max))
	}
	return nil
}

// ==================== 格式验证 ====================

var (
	emailRegex    = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	usernameRegex = regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)
	phoneRegex    = regexp.MustCompile(`^1[3-9]\d{9}$`)
)

// Email 邮箱格式验证
func Email(value, field string) error {
	if value != "" && !emailRegex.MatchString(value) {
		return NewValidationError(field, "邮箱格式不正确")
	}
	return nil
}

// Username 用户名格式验证（只允许字母、数字、下划线、横线）
func Username(value, field string) error {
	if value != "" && !usernameRegex.MatchString(value) {
		return NewValidationError(field, "只能包含字母、数字、下划线和横线")
	}
	return nil
}

// Phone 手机号格式验证
func Phone(value, field string) error {
	if value != "" && !phoneRegex.MatchString(value) {
		return NewValidationError(field, "手机号格式不正确")
	}
	return nil
}

// ==================== 数值验证 ====================

// Min 最小值验证
func Min(value int, field string, min int) error {
	if value < min {
		return NewValidationError(field, fmt.Sprintf("不能小于%d", min))
	}
	return nil
}

// Max 最大值验证
func Max(value int, field string, max int) error {
	if value > max {
		return NewValidationError(field, fmt.Sprintf("不能大于%d", max))
	}
	return nil
}

// Range 数值范围验证
func Range(value int, field string, min, max int) error {
	if value < min || value > max {
		return NewValidationError(field, fmt.Sprintf("应在%d-%d之间", min, max))
	}
	return nil
}

// Positive 正数验证
func Positive(value int, field string) error {
	if value <= 0 {
		return NewValidationError(field, "必须为正数")
	}
	return nil
}

// ==================== 切片验证 ====================

// NotEmpty 非空切片验证
func NotEmpty[T any](values []T, field string) error {
	if len(values) == 0 {
		return NewValidationError(field, "不能为空")
	}
	return nil
}

// MaxItems 最大元素数量验证
func MaxItems[T any](values []T, field string, max int) error {
	if len(values) > max {
		return NewValidationError(field, fmt.Sprintf("最多只能有%d个元素", max))
	}
	return nil
}

// ==================== 复合验证 ====================

// Validator 验证器
type Validator struct {
	errors []error
}

// NewValidator 创建验证器
func NewValidator() *Validator {
	return &Validator{errors: make([]error, 0)}
}

// Check 检查并收集错误
func (v *Validator) Check(err error) *Validator {
	if err != nil {
		v.errors = append(v.errors, err)
	}
	return v
}

// HasErrors 是否有错误
func (v *Validator) HasErrors() bool {
	return len(v.errors) > 0
}

// Error 返回第一个错误
func (v *Validator) Error() error {
	if len(v.errors) == 0 {
		return nil
	}
	return v.errors[0]
}

// Errors 返回所有错误
func (v *Validator) Errors() []error {
	return v.errors
}

// FirstErrorMessage 返回第一个错误消息
func (v *Validator) FirstErrorMessage() string {
	if len(v.errors) == 0 {
		return ""
	}
	return v.errors[0].Error()
}

// AllErrorMessages 返回所有错误消息
func (v *Validator) AllErrorMessages() []string {
	messages := make([]string, len(v.errors))
	for i, err := range v.errors {
		messages[i] = err.Error()
	}
	return messages
}
