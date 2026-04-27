// Package errors 提供统一的错误类型定义
//
// 所有 dong 家族 CLI 工具共享这些错误类型
package errors

import "fmt"

// ErrorCode 错误代码常量
type ErrorCode string

const (
	ErrValidation    ErrorCode = "VALIDATION_ERROR"
	ErrNotFound      ErrorCode = "NOT_FOUND"
	ErrConflict      ErrorCode = "CONFLICT"
	ErrPermission    ErrorCode = "PERMISSION_DENIED"
	ErrInternal      ErrorCode = "INTERNAL_ERROR"
	ErrValue         ErrorCode = "VALUE_ERROR"
	ErrType          ErrorCode = "TYPE_ERROR"
	ErrKey           ErrorCode = "KEY_ERROR"
	ErrAttribute     ErrorCode = "ATTRIBUTE_ERROR"
	ErrFileNotFound  ErrorCode = "NOT_FOUND"
	ErrIO            ErrorCode = "IO_ERROR"
	ErrUnknown       ErrorCode = "UNKNOWN_ERROR"
)

// DongError dong 家族 CLI 基础错误类型
type DongError struct {
	Code    ErrorCode
	Message string
	Details map[string]interface{}
}

// Error 实现 error 接口
func (e *DongError) Error() string {
	return e.Message
}

// ToDict 转换为字典格式
func (e *DongError) ToDict() map[string]interface{} {
	result := map[string]interface{}{
		"code":    string(e.Code),
		"message": e.Message,
	}
	if len(e.Details) > 0 {
		result["details"] = e.Details
	}
	return result
}

// NewDongError 创建新的 DongError
func NewDongError(code ErrorCode, message string, details map[string]interface{}) *DongError {
	return &DongError{
		Code:    code,
		Message: message,
		Details: details,
	}
}

// ValidationError 验证错误
type ValidationError struct {
	Code    ErrorCode
	Field   string
	Message string
	Details map[string]interface{}
}

// Error 实现 error 接口
func (e *ValidationError) Error() string {
	return fmt.Sprintf("%s (field: %s)", e.Message, e.Field)
}

// ToDict 转换为字典格式
func (e *ValidationError) ToDict() map[string]interface{} {
	details := map[string]interface{}{
		"field": e.Field,
	}
	for k, v := range e.Details {
		details[k] = v
	}
	return map[string]interface{}{
		"code":    string(e.Code),
		"message": e.Message,
		"details": details,
	}
}

// NewValidationError 创建验证错误
func NewValidationError(field, message string) *ValidationError {
	return &ValidationError{
		Code:    ErrValidation,
		Field:   field,
		Message: message,
		Details: make(map[string]interface{}),
	}
}

// NotFoundError 资源未找到错误
type NotFoundError struct {
	ResourceType string
	ResourceID   interface{}
	Message      string
	Details      map[string]interface{}
}

// Error 实现 error 接口
func (e *NotFoundError) Error() string {
	if e.Message != "" {
		return e.Message
	}
	if e.ResourceType != "" && e.ResourceID != nil {
		return fmt.Sprintf("%s %v 不存在", e.ResourceType, e.ResourceID)
	}
	if e.ResourceType != "" {
		return e.ResourceType + " 不存在"
	}
	return "资源未找到"
}

// ToDict 转换为字典格式
func (e *NotFoundError) ToDict() map[string]interface{} {
	details := make(map[string]interface{})
	if e.ResourceType != "" {
		details["resource_type"] = e.ResourceType
	}
	if e.ResourceID != nil {
		details["resource_id"] = fmt.Sprintf("%v", e.ResourceID)
	}
	for k, v := range e.Details {
		details[k] = v
	}
	return map[string]interface{}{
		"code":    string(ErrNotFound),
		"message": e.Error(),
		"details": details,
	}
}

// NewNotFoundError 创建未找到错误
func NewNotFoundError(resourceType string, resourceID interface{}) *NotFoundError {
	return &NotFoundError{
		ResourceType: resourceType,
		ResourceID:   resourceID,
		Details:      make(map[string]interface{}),
	}
}

// NewNotFoundErrorWithMessage 创建带自定义消息的未找到错误
func NewNotFoundErrorWithMessage(message string) *NotFoundError {
	return &NotFoundError{
		Message: message,
		Details: make(map[string]interface{}),
	}
}

// ConflictError 冲突错误
type ConflictError struct {
	ResourceType      string
	ConflictingField  string
	ConflictingValue  interface{}
	Message           string
	Details           map[string]interface{}
}

// Error 实现 error 接口
func (e *ConflictError) Error() string {
	if e.Message == "" {
		return fmt.Sprintf("%s '%v' 已被使用", e.ConflictingField, e.ConflictingValue)
	}
	return e.Message
}

// ToDict 转换为字典格式
func (e *ConflictError) ToDict() map[string]interface{} {
	details := map[string]interface{}{
		"resource_type":      e.ResourceType,
		"conflicting_field":  e.ConflictingField,
		"conflicting_value":  fmt.Sprintf("%v", e.ConflictingValue),
	}
	for k, v := range e.Details {
		details[k] = v
	}
	return map[string]interface{}{
		"code":    string(ErrConflict),
		"message": e.Error(),
		"details": details,
	}
}

// NewConflictError 创建冲突错误
func NewConflictError(resourceType, conflictingField string, conflictingValue interface{}) *ConflictError {
	return &ConflictError{
		ResourceType:     resourceType,
		ConflictingField: conflictingField,
		ConflictingValue: conflictingValue,
		Details:          make(map[string]interface{}),
	}
}

// ExtractErrorInfo 从 error 中提取错误信息
func ExtractErrorInfo(err error) map[string]string {
	if err == nil {
		return nil
	}

	switch e := err.(type) {
	case *DongError:
		return map[string]string{
			"code":    string(e.Code),
			"message": e.Message,
		}
	case *ValidationError:
		return map[string]string{
			"code":    string(e.Code),
			"message": e.Message,
		}
	case *NotFoundError:
		return map[string]string{
			"code":    string(ErrNotFound),
			"message": e.Error(),
		}
	case *ConflictError:
		return map[string]string{
			"code":    string(ErrConflict),
			"message": e.Error(),
		}
	default:
		// 内置异常类型映射
		return map[string]string{
			"code":    string(ErrUnknown),
			"message": err.Error(),
		}
	}
}
