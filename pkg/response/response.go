package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Response 统一API响应结构
type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

// Success 成功响应
func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code:    200,
		Message: "操作成功",
		Data:    data,
	})
}

// Error 错误响应
func Error(c *gin.Context, code int, message string) {
	c.JSON(http.StatusOK, Response{
		Code:    code,
		Message: message,
		Data:    nil,
	})
}

// ServerError 服务器错误
func ServerError(c *gin.Context) {
	c.JSON(http.StatusInternalServerError, Response{
		Code:    500,
		Message: "服务器内部错误",
		Data:    nil,
	})
}

// Unauthorized 未授权
func Unauthorized(c *gin.Context) {
	c.JSON(http.StatusUnauthorized, Response{
		Code:    401,
		Message: "未授权或授权已过期",
		Data:    nil,
	})
}

// Forbidden 权限不足
func Forbidden(c *gin.Context) {
	c.JSON(http.StatusForbidden, Response{
		Code:    403,
		Message: "权限不足",
		Data:    nil,
	})
}

// NotFound 资源不存在
func NotFound(c *gin.Context) {
	c.JSON(http.StatusNotFound, Response{
		Code:    404,
		Message: "请求的资源不存在",
		Data:    nil,
	})
}

// BadRequest 请求参数错误
func BadRequest(c *gin.Context, message string) {
	c.JSON(http.StatusBadRequest, Response{
		Code:    400,
		Message: message,
		Data:    nil,
	})
}

// SuccessWithPagination 带分页信息的成功响应
func SuccessWithPagination(c *gin.Context, data interface{}, total int64, pageNum, pageSize int) {
	c.JSON(http.StatusOK, Response{
		Code:    200,
		Message: "操作成功",
		Data: gin.H{
			"list":     data,
			"total":    total,
			"pageNum":  pageNum,
			"pageSize": pageSize,
		},
	})
}