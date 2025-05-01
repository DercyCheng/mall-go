package response

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	appErrors "mall-go/pkg/errors"
)

func init() {
	// 设置Gin为测试模式
	gin.SetMode(gin.TestMode)
}

// 创建测试用的gin上下文
func getTestContext() (*gin.Context, *httptest.ResponseRecorder) {
	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)
	ctx.Request = httptest.NewRequest("GET", "/test", nil)
	return ctx, recorder
}

func TestSuccess(t *testing.T) {
	c, recorder := getTestContext()

	// 测试数据
	testData := map[string]interface{}{
		"name": "John",
		"age":  30,
	}

	// 调用测试函数
	Success(c, testData)

	// 断言
	assert.Equal(t, http.StatusOK, recorder.Code)
	assert.Contains(t, recorder.Body.String(), `"code":200`)
	assert.Contains(t, recorder.Body.String(), `"message":"操作成功"`)
	assert.Contains(t, recorder.Body.String(), `"name":"John"`)
	assert.Contains(t, recorder.Body.String(), `"age":30`)
}

func TestSuccessWithMessage(t *testing.T) {
	c, recorder := getTestContext()

	// 调用测试函数
	SuccessWithMessage(c, "自定义成功消息", "success data")

	// 断言
	assert.Equal(t, http.StatusOK, recorder.Code)
	assert.Contains(t, recorder.Body.String(), `"code":200`)
	assert.Contains(t, recorder.Body.String(), `"message":"自定义成功消息"`)
	assert.Contains(t, recorder.Body.String(), `"data":"success data"`)
}

func TestCreated(t *testing.T) {
	c, recorder := getTestContext()

	// 调用测试函数
	Created(c, "created data")

	// 断言
	assert.Equal(t, http.StatusCreated, recorder.Code)
	assert.Contains(t, recorder.Body.String(), `"code":201`)
	assert.Contains(t, recorder.Body.String(), `"message":"创建成功"`)
	assert.Contains(t, recorder.Body.String(), `"data":"created data"`)
}

func TestNoContent(t *testing.T) {
	// 创建一个新的 gin 路由
	router := gin.New()
	router.GET("/test", func(c *gin.Context) {
		NoContent(c)
	})

	// 创建请求
	req := httptest.NewRequest("GET", "/test", nil)
	resp := httptest.NewRecorder()
	
	// 执行请求
	router.ServeHTTP(resp, req)
	
	// 断言
	assert.Equal(t, http.StatusNoContent, resp.Code)
	assert.Empty(t, resp.Body.String())
}

func TestError(t *testing.T) {
	c, recorder := getTestContext()

	// 创建应用错误
	appErr := appErrors.ValidationError(appErrors.CodeValidationError, "验证失败")

	// 调用测试函数
	Error(c, appErr)

	// 断言
	assert.Equal(t, http.StatusBadRequest, recorder.Code)
	assert.Contains(t, recorder.Body.String(), `"code":400`)
	assert.Contains(t, recorder.Body.String(), `"message":"验证失败"`)
	assert.Contains(t, recorder.Body.String(), `"data":null`)
}

func TestErrorWithStandardError(t *testing.T) {
	c, recorder := getTestContext()

	// 创建标准错误
	stdErr := errors.New("标准错误")

	// 调用测试函数
	Error(c, stdErr)

	// 断言
	assert.Equal(t, http.StatusInternalServerError, recorder.Code)
	assert.Contains(t, recorder.Body.String(), `"code":500`)
	assert.Contains(t, recorder.Body.String(), `"message":"服务器内部错误"`)
}

func TestErrorWithCode(t *testing.T) {
	c, recorder := getTestContext()

	// 调用测试函数
	ErrorWithCode(c, 403, "自定义错误")

	// 断言
	assert.Equal(t, http.StatusForbidden, recorder.Code)
	assert.Contains(t, recorder.Body.String(), `"code":403`)
	assert.Contains(t, recorder.Body.String(), `"message":"自定义错误"`)
	assert.Contains(t, recorder.Body.String(), `"data":null`)
}

func TestBadRequest(t *testing.T) {
	c, recorder := getTestContext()

	// 调用测试函数
	BadRequest(c, "请求参数错误")

	// 断言
	assert.Equal(t, http.StatusBadRequest, recorder.Code)
	assert.Contains(t, recorder.Body.String(), `"code":400`)
	assert.Contains(t, recorder.Body.String(), `"message":"请求参数错误"`)
}

func TestUnauthorized(t *testing.T) {
	t.Run("WithMessage", func(t *testing.T) {
		c, recorder := getTestContext()

		// 调用测试函数
		Unauthorized(c, "令牌已过期")

		// 断言
		assert.Equal(t, http.StatusUnauthorized, recorder.Code)
		assert.Contains(t, recorder.Body.String(), `"code":401`)
		assert.Contains(t, recorder.Body.String(), `"message":"令牌已过期"`)
	})

	t.Run("DefaultMessage", func(t *testing.T) {
		c, recorder := getTestContext()

		// 调用测试函数
		Unauthorized(c, "")

		// 断言
		assert.Equal(t, http.StatusUnauthorized, recorder.Code)
		assert.Contains(t, recorder.Body.String(), `"code":401`)
		assert.Contains(t, recorder.Body.String(), `"message":"未授权访问"`)
	})
}

func TestForbidden(t *testing.T) {
	t.Run("WithMessage", func(t *testing.T) {
		c, recorder := getTestContext()

		// 调用测试函数
		Forbidden(c, "权限不足")

		// 断言
		assert.Equal(t, http.StatusForbidden, recorder.Code)
		assert.Contains(t, recorder.Body.String(), `"code":403`)
		assert.Contains(t, recorder.Body.String(), `"message":"权限不足"`)
	})

	t.Run("DefaultMessage", func(t *testing.T) {
		c, recorder := getTestContext()

		// 调用测试函数
		Forbidden(c, "")

		// 断言
		assert.Equal(t, http.StatusForbidden, recorder.Code)
		assert.Contains(t, recorder.Body.String(), `"code":403`)
		assert.Contains(t, recorder.Body.String(), `"message":"权限不足"`)
	})
}

func TestNotFound(t *testing.T) {
	t.Run("WithMessage", func(t *testing.T) {
		c, recorder := getTestContext()

		// 调用测试函数
		NotFound(c, "用户不存在")

		// 断言
		assert.Equal(t, http.StatusNotFound, recorder.Code)
		assert.Contains(t, recorder.Body.String(), `"code":404`)
		assert.Contains(t, recorder.Body.String(), `"message":"用户不存在"`)
	})

	t.Run("DefaultMessage", func(t *testing.T) {
		c, recorder := getTestContext()

		// 调用测试函数
		NotFound(c, "")

		// 断言
		assert.Equal(t, http.StatusNotFound, recorder.Code)
		assert.Contains(t, recorder.Body.String(), `"code":404`)
		assert.Contains(t, recorder.Body.String(), `"message":"资源不存在"`)
	})
}

func TestConflict(t *testing.T) {
	c, recorder := getTestContext()

	// 调用测试函数
	Conflict(c, "资源已存在")

	// 断言
	assert.Equal(t, http.StatusConflict, recorder.Code)
	assert.Contains(t, recorder.Body.String(), `"code":409`)
	assert.Contains(t, recorder.Body.String(), `"message":"资源已存在"`)
}

func TestInternalServerError(t *testing.T) {
	t.Run("WithMessage", func(t *testing.T) {
		c, recorder := getTestContext()

		// 调用测试函数
		InternalServerError(c, "数据库连接失败")

		// 断言
		assert.Equal(t, http.StatusInternalServerError, recorder.Code)
		assert.Contains(t, recorder.Body.String(), `"code":500`)
		assert.Contains(t, recorder.Body.String(), `"message":"数据库连接失败"`)
	})

	t.Run("DefaultMessage", func(t *testing.T) {
		c, recorder := getTestContext()

		// 调用测试函数
		InternalServerError(c, "")

		// 断言
		assert.Equal(t, http.StatusInternalServerError, recorder.Code)
		assert.Contains(t, recorder.Body.String(), `"code":500`)
		assert.Contains(t, recorder.Body.String(), `"message":"服务器内部错误"`)
	})
}

func TestValidationError(t *testing.T) {
	c, recorder := getTestContext()

	// 调用测试函数
	ValidationError(c, "用户名不能为空")

	// 断言
	assert.Equal(t, http.StatusBadRequest, recorder.Code)
	assert.Contains(t, recorder.Body.String(), `"code":400`)
	assert.Contains(t, recorder.Body.String(), `"message":"用户名不能为空"`)
}

func TestPageSuccess(t *testing.T) {
	c, recorder := getTestContext()

	// 测试数据
	testData := []map[string]interface{}{
		{"id": 1, "name": "Item 1"},
		{"id": 2, "name": "Item 2"},
		{"id": 3, "name": "Item 3"},
	}

	// 调用测试函数
	PageSuccess(c, testData, 25, 1, 10)

	// 断言
	assert.Equal(t, http.StatusOK, recorder.Code)
	assert.Contains(t, recorder.Body.String(), `"code":200`)
	assert.Contains(t, recorder.Body.String(), `"message":"操作成功"`)
	assert.Contains(t, recorder.Body.String(), `"total":25`)
	assert.Contains(t, recorder.Body.String(), `"page":1`)
	assert.Contains(t, recorder.Body.String(), `"pageSize":10`)
	assert.Contains(t, recorder.Body.String(), `"totalPage":3`)
	assert.Contains(t, recorder.Body.String(), `"Item 1"`)
	assert.Contains(t, recorder.Body.String(), `"Item 3"`)
}