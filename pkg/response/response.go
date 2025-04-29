package response

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"mall-go/pkg/errors"
)

// Response 标准API响应结构
type Response struct {
	Code    int         `json:"code"`    // HTTP状态码
	Message string      `json:"message"` // 响应消息
	Data    interface{} `json:"data"`    // 响应数据
}

// Success 成功响应
func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code:    http.StatusOK,
		Message: "操作成功",
		Data:    data,
	})
}

// SuccessWithMessage 带自定义消息的成功响应
func SuccessWithMessage(c *gin.Context, message string, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code:    http.StatusOK,
		Message: message,
		Data:    data,
	})
}

// Created 创建成功响应
func Created(c *gin.Context, data interface{}) {
	c.JSON(http.StatusCreated, Response{
		Code:    http.StatusCreated,
		Message: "创建成功",
		Data:    data,
	})
}

// NoContent 无内容成功响应
func NoContent(c *gin.Context) {
	c.Status(http.StatusNoContent)
}

// Error 错误响应
func Error(c *gin.Context, err error) {
	appErr, ok := errors.As(err)
	if !ok {
		// 如果不是AppError，转换为内部错误
		appErr = errors.InternalError(errors.CodeInternalError, "服务器内部错误", err)
	}

	statusCode := appErr.HTTPStatusCode()
	c.JSON(statusCode, Response{
		Code:    statusCode,
		Message: appErr.Message,
		Data:    nil,
	})
}

// ErrorWithCode 带自定义状态码的错误响应
func ErrorWithCode(c *gin.Context, code int, message string) {
	c.JSON(code, Response{
		Code:    code,
		Message: message,
		Data:    nil,
	})
}

// BadRequest 请求参数错误响应
func BadRequest(c *gin.Context, message string) {
	c.JSON(http.StatusBadRequest, Response{
		Code:    http.StatusBadRequest,
		Message: message,
		Data:    nil,
	})
}

// Unauthorized 未授权错误响应
func Unauthorized(c *gin.Context, message string) {
	if message == "" {
		message = "未授权访问"
	}
	c.JSON(http.StatusUnauthorized, Response{
		Code:    http.StatusUnauthorized,
		Message: message,
		Data:    nil,
	})
}

// Forbidden 禁止访问错误响应
func Forbidden(c *gin.Context, message string) {
	if message == "" {
		message = "权限不足"
	}
	c.JSON(http.StatusForbidden, Response{
		Code:    http.StatusForbidden,
		Message: message,
		Data:    nil,
	})
}

// NotFound 资源不存在错误响应
func NotFound(c *gin.Context, message string) {
	if message == "" {
		message = "资源不存在"
	}
	c.JSON(http.StatusNotFound, Response{
		Code:    http.StatusNotFound,
		Message: message,
		Data:    nil,
	})
}

// Conflict 资源冲突错误响应
func Conflict(c *gin.Context, message string) {
	c.JSON(http.StatusConflict, Response{
		Code:    http.StatusConflict,
		Message: message,
		Data:    nil,
	})
}

// InternalServerError 服务器内部错误响应
func InternalServerError(c *gin.Context, message string) {
	if message == "" {
		message = "服务器内部错误"
	}
	c.JSON(http.StatusInternalServerError, Response{
		Code:    http.StatusInternalServerError,
		Message: message,
		Data:    nil,
	})
}

// ValidationError 验证错误响应
func ValidationError(c *gin.Context, message string) {
	c.JSON(http.StatusBadRequest, Response{
		Code:    http.StatusBadRequest,
		Message: message,
		Data:    nil,
	})
}

// Page 分页响应数据
type Page struct {
	List      interface{} `json:"list"`      // 分页数据
	Total     int64       `json:"total"`     // 总记录数
	Page      int         `json:"page"`      // 当前页码
	PageSize  int         `json:"pageSize"`  // 每页记录数
	TotalPage int         `json:"totalPage"` // 总页数
}

// PageSuccess 分页成功响应
func PageSuccess(c *gin.Context, list interface{}, total int64, page, pageSize int) {
	// 计算总页数
	totalPage := int(total) / pageSize
	if int(total)%pageSize > 0 {
		totalPage++
	}

	c.JSON(http.StatusOK, Response{
		Code:    http.StatusOK,
		Message: "操作成功",
		Data: Page{
			List:      list,
			Total:     total,
			Page:      page,
			PageSize:  pageSize,
			TotalPage: totalPage,
		},
	})
}
