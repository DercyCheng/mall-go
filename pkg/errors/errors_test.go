package errors

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAppError(t *testing.T) {
	// 测试创建不同类型的错误
	t.Run("CreateErrors", func(t *testing.T) {
		// 验证错误
		validationErr := ValidationError(CodeValidationError, "输入验证失败")
		assert.Equal(t, TypeValidation, validationErr.Type)
		assert.Equal(t, CodeValidationError, validationErr.Code)
		assert.Equal(t, "输入验证失败", validationErr.Message)
		assert.Equal(t, http.StatusBadRequest, validationErr.HTTPStatusCode())

		// 业务错误
		businessErr := BusinessError(CodeUserAlreadyExists, "用户已存在")
		assert.Equal(t, TypeBusiness, businessErr.Type)
		assert.Equal(t, CodeUserAlreadyExists, businessErr.Code)
		assert.Equal(t, "用户已存在", businessErr.Message)
		assert.Equal(t, http.StatusBadRequest, businessErr.HTTPStatusCode())

		// 数据库错误
		origErr := fmt.Errorf("database connection failed")
		dbErr := DatabaseError(CodeDatabaseError, "数据库操作失败", origErr)
		assert.Equal(t, TypeDatabase, dbErr.Type)
		assert.Equal(t, CodeDatabaseError, dbErr.Code)
		assert.Equal(t, "数据库操作失败: database connection failed", dbErr.Error())
		assert.Equal(t, http.StatusInternalServerError, dbErr.HTTPStatusCode())

		// 未授权错误
		authErr := UnauthorizedError(CodeInvalidToken, "无效的令牌")
		assert.Equal(t, TypeUnauthorized, authErr.Type)
		assert.Equal(t, http.StatusUnauthorized, authErr.HTTPStatusCode())

		// 禁止访问错误
		forbiddenErr := ForbiddenError(CodePermissionDenied, "权限不足")
		assert.Equal(t, TypeForbidden, forbiddenErr.Type)
		assert.Equal(t, http.StatusForbidden, forbiddenErr.HTTPStatusCode())
	})

	t.Run("ErrorWrapping", func(t *testing.T) {
		// 测试错误包装
		origErr := fmt.Errorf("原始错误")
		wrappedErr := Wrap(origErr, "包装的错误消息")
		
		assert.Contains(t, wrappedErr.Error(), "包装的错误消息")
		assert.Contains(t, wrappedErr.Error(), "原始错误")
		
		// 测试包装AppError
		appErr := BusinessError(CodeUserAlreadyExists, "用户已存在")
		wrappedAppErr := Wrap(appErr, "注册失败")
		
		if unwrappedErr, ok := As(wrappedAppErr); ok {
			assert.Equal(t, TypeBusiness, unwrappedErr.Type)
			assert.Equal(t, CodeUserAlreadyExists, unwrappedErr.Code)
		} else {
			t.Fatal("无法正确解包AppError")
		}
	})

	t.Run("ErrorHelpers", func(t *testing.T) {
		// 测试错误帮助函数
		notFoundErr := NotFoundError(CodeUserNotFound, "用户不存在")
		assert.True(t, IsNotFound(notFoundErr))
		assert.False(t, IsValidationError(notFoundErr))
		
		validationErr := ValidationError(CodeValidationError, "输入验证失败")
		assert.True(t, IsValidationError(validationErr))
		assert.False(t, IsUnauthorized(validationErr))
		
		unauthorizedErr := UnauthorizedError(CodeInvalidToken, "无效的令牌")
		assert.True(t, IsUnauthorized(unauthorizedErr))
	})

	t.Run("ErrorCompatibility", func(t *testing.T) {
		// 测试与标准库的兼容性
		stdErr := New("标准错误")
		assert.Contains(t, stdErr.Error(), "标准错误")
		
		err1 := fmt.Errorf("错误1")
		err2 := fmt.Errorf("包装: %w", err1)
		assert.True(t, Is(err2, err1))
	})
}