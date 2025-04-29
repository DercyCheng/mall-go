package errors

import (
	std_errors "errors"
	"fmt"
	"net/http"
)

// 错误类型
const (
	TypeValidation   = "validation"
	TypeBusiness     = "business"
	TypeDatabase     = "database"
	TypeInternal     = "internal"
	TypeUnauthorized = "unauthorized"
	TypeForbidden    = "forbidden"
	TypeNotFound     = "not_found"
	TypeConflict     = "conflict"
	TypeThirdParty   = "third_party"
)

// 错误代码
const (
	// 通用错误
	CodeUnknown         = "unknown"
	CodeInternalError   = "internal_error"
	CodeValidationError = "validation_error"
	CodeDatabaseError   = "database_error"

	// 认证和授权错误
	CodeInvalidCredentials = "invalid_credentials"
	CodeInvalidToken       = "invalid_token"
	CodeTokenExpired       = "token_expired"
	CodePermissionDenied   = "permission_denied"

	// 用户相关错误
	CodeUserNotFound      = "user_not_found"
	CodeUserAlreadyExists = "user_already_exists"
	CodePasswordMismatch  = "password_mismatch"

	// 角色相关错误
	CodeRoleNotFound      = "role_not_found"
	CodeRoleAlreadyExists = "role_already_exists"

	// 并发错误
	CodeVersionConflict  = "version_conflict"
	CodeOperationTimeout = "operation_timeout"
)

// AppError 是应用程序的标准错误类型
type AppError struct {
	Type    string // 错误类型
	Code    string // 错误代码
	Message string // 错误消息
	Err     error  // 原始错误
}

// Error 实现error接口
func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %s", e.Message, e.Err.Error())
	}
	return e.Message
}

// Unwrap 实现errors.Unwrap接口
func (e *AppError) Unwrap() error {
	return e.Err
}

// HTTPStatusCode 返回对应的HTTP状态码
func (e *AppError) HTTPStatusCode() int {
	switch e.Type {
	case TypeValidation:
		return http.StatusBadRequest
	case TypeBusiness:
		return http.StatusBadRequest
	case TypeUnauthorized:
		return http.StatusUnauthorized
	case TypeForbidden:
		return http.StatusForbidden
	case TypeNotFound:
		return http.StatusNotFound
	case TypeConflict:
		return http.StatusConflict
	default:
		return http.StatusInternalServerError
	}
}

// NewAppError 创建新的应用错误
func NewAppError(errType, code, message string, err error) *AppError {
	return &AppError{
		Type:    errType,
		Code:    code,
		Message: message,
		Err:     err,
	}
}

// ValidationError 创建验证错误
func ValidationError(code, message string) *AppError {
	return NewAppError(TypeValidation, code, message, nil)
}

// BusinessError 创建业务错误
func BusinessError(code, message string) *AppError {
	return NewAppError(TypeBusiness, code, message, nil)
}

// DatabaseError 创建数据库错误
func DatabaseError(code, message string, err error) *AppError {
	return NewAppError(TypeDatabase, code, message, err)
}

// InternalError 创建内部错误
func InternalError(code, message string, err error) *AppError {
	return NewAppError(TypeInternal, code, message, err)
}

// UnauthorizedError 创建未授权错误
func UnauthorizedError(code, message string) *AppError {
	return NewAppError(TypeUnauthorized, code, message, nil)
}

// ForbiddenError 创建禁止访问错误
func ForbiddenError(code, message string) *AppError {
	return NewAppError(TypeForbidden, code, message, nil)
}

// NotFoundError 创建资源不存在错误
func NotFoundError(code, message string) *AppError {
	return NewAppError(TypeNotFound, code, message, nil)
}

// ConflictError 创建资源冲突错误
func ConflictError(code, message string) *AppError {
	return NewAppError(TypeConflict, code, message, nil)
}

// ThirdPartyError 创建第三方服务错误
func ThirdPartyError(code, message string, err error) *AppError {
	return NewAppError(TypeThirdParty, code, message, err)
}

// As 将普通错误转换为AppError
func As(err error) (*AppError, bool) {
	var appErr *AppError
	if err == nil {
		return nil, false
	}

	// 检查err是否已经是AppError
	if e, ok := err.(*AppError); ok {
		return e, true
	}

	// 尝试使用标准库的errors.As
	if ok := std_errors.As(err, &appErr); ok {
		return appErr, true
	}

	return nil, false
}

// Wrap 包装错误
func Wrap(err error, message string) error {
	if err == nil {
		return nil
	}

	if appErr, ok := As(err); ok {
		// 如果已经是AppError，则保留错误类型和代码，更新消息
		return NewAppError(appErr.Type, appErr.Code, message, appErr)
	}

	// 否则创建内部错误
	return InternalError(CodeUnknown, message, err)
}

// IsNotFound 检查是否为NotFound错误
func IsNotFound(err error) bool {
	if appErr, ok := As(err); ok {
		return appErr.Type == TypeNotFound
	}
	return false
}

// IsValidationError 检查是否为验证错误
func IsValidationError(err error) bool {
	if appErr, ok := As(err); ok {
		return appErr.Type == TypeValidation
	}
	return false
}

// IsUnauthorized 检查是否为未授权错误
func IsUnauthorized(err error) bool {
	if appErr, ok := As(err); ok {
		return appErr.Type == TypeUnauthorized
	}
	return false
}

// 为了兼容标准库的errors函数，提供以下函数
func New(text string) error {
	return BusinessError(CodeUnknown, text)
}

func Is(err, target error) bool {
	return std_errors.Is(err, target)
}
