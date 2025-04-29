// Package docs 为用户服务生成Swagger文档。
//
// API文档提供了用户服务所有API的详细信息。
//
// 基本路径: /api/v1
//
// swagger:meta
package docs

import (
	"mall-go/services/user-service/application/dto"
)

// 用户注册请求
// swagger:parameters registerUser
type RegisterUserParams struct {
	// 在Body中
	// 必需: true
	// swagger:model
	Body dto.UserCreateRequest
}

// 用户注册成功响应
// swagger:response registerUserResponse
type RegisterUserResponse struct {
	// 在Body中
	// 例如: {"code":200,"message":"操作成功","data":{"id":"12345"}}
	Body struct {
		// HTTP状态码
		// 例如: 200
		// 必需: true
		Code int `json:"code"`
		// 状态消息
		// 例如: 操作成功
		// 必需: true
		Message string `json:"message"`
		// 响应数据
		// 必需: true
		Data struct {
			// 用户ID
			// 例如: 12345
			// 必需: true
			ID string `json:"id"`
		} `json:"data"`
	}
}

// 用户注册失败响应
// swagger:response badRequestResponse
type BadRequestResponse struct {
	// 在Body中
	// 例如: {"code":400,"message":"无效的请求参数: 用户名已存在","data":null}
	Body struct {
		// HTTP状态码
		// 例如: 400
		// 必需: true
		Code int `json:"code"`
		// 错误消息
		// 例如: 无效的请求参数: 用户名已存在
		// 必需: true
		Message string `json:"message"`
		// 响应数据
		// 必需: true
		Data interface{} `json:"data"`
	}
}

// 用户登录请求
// swagger:parameters loginUser
type LoginUserParams struct {
	// 在Body中
	// 必需: true
	// swagger:model
	Body dto.UserLoginRequest
}

// 用户登录成功响应
// swagger:response loginUserResponse
type LoginUserResponse struct {
	// 在Body中
	// 例如: {"code":200,"message":"操作成功","data":{"token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...","tokenHead":"Bearer","expireAt":"2023-05-01T12:00:00Z","user":{"id":"12345","username":"testuser","email":"test@example.com","nickName":"Test User","status":1}}}
	Body struct {
		// HTTP状态码
		// 例如: 200
		// 必需: true
		Code int `json:"code"`
		// 状态消息
		// 例如: 操作成功
		// 必需: true
		Message string `json:"message"`
		// 响应数据
		// 必需: true
		Data dto.UserLoginResponse `json:"data"`
	}
}

// 未授权响应
// swagger:response unauthorizedResponse
type UnauthorizedResponse struct {
	// 在Body中
	// 例如: {"code":401,"message":"用户名或密码错误","data":null}
	Body struct {
		// HTTP状态码
		// 例如: 401
		// 必需: true
		Code int `json:"code"`
		// 错误消息
		// 例如: 用户名或密码错误
		// 必需: true
		Message string `json:"message"`
		// 响应数据
		// 必需: true
		Data interface{} `json:"data"`
	}
}

// 获取用户信息请求
// swagger:parameters getUserInfo
type GetUserInfoParams struct {
	// 访问令牌
	// 在Header中: Authorization
	// 例如: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
	// 必需: true
	Authorization string `header:"Authorization"`
}

// 获取用户信息成功响应
// swagger:response getUserInfoResponse
type GetUserInfoResponse struct {
	// 在Body中
	// 例如: {"code":200,"message":"操作成功","data":{"id":"12345","username":"testuser","email":"test@example.com","nickName":"Test User","status":1}}
	Body struct {
		// HTTP状态码
		// 例如: 200
		// 必需: true
		Code int `json:"code"`
		// 状态消息
		// 例如: 操作成功
		// 必需: true
		Message string `json:"message"`
		// 响应数据
		// 必需: true
		Data dto.UserDTO `json:"data"`
	}
}

// 更新用户请求
// swagger:parameters updateUser
type UpdateUserParams struct {
	// 用户ID
	// 在Path中
	// 例如: 12345
	// 必需: true
	ID string `path:"id"`

	// 访问令牌
	// 在Header中: Authorization
	// 例如: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
	// 必需: true
	Authorization string `header:"Authorization"`

	// 在Body中
	// 必需: true
	// swagger:model
	Body dto.UserUpdateRequest
}

// 更新用户成功响应
// swagger:response updateUserResponse
type UpdateUserResponse struct {
	// 在Body中
	// 例如: {"code":200,"message":"操作成功","data":null}
	Body struct {
		// HTTP状态码
		// 例如: 200
		// 必需: true
		Code int `json:"code"`
		// 状态消息
		// 例如: 操作成功
		// 必需: true
		Message string `json:"message"`
		// 响应数据
		// 必需: true
		Data interface{} `json:"data"`
	}
}

// 删除用户请求
// swagger:parameters deleteUser
type DeleteUserParams struct {
	// 用户ID
	// 在Path中
	// 例如: 12345
	// 必需: true
	ID string `path:"id"`

	// 访问令牌
	// 在Header中: Authorization
	// 例如: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
	// 必需: true
	Authorization string `header:"Authorization"`
}

// 删除用户成功响应
// swagger:response deleteUserResponse
type DeleteUserResponse struct {
	// 在Body中
	// 例如: {"code":200,"message":"操作成功","data":null}
	Body struct {
		// HTTP状态码
		// 例如: 200
		// 必需: true
		Code int `json:"code"`
		// 状态消息
		// 例如: 操作成功
		// 必需: true
		Message string `json:"message"`
		// 响应数据
		// 必需: true
		Data interface{} `json:"data"`
	}
}

// 服务器内部错误响应
// swagger:response internalServerErrorResponse
type InternalServerErrorResponse struct {
	// 在Body中
	// 例如: {"code":500,"message":"服务器内部错误","data":null}
	Body struct {
		// HTTP状态码
		// 例如: 500
		// 必需: true
		Code int `json:"code"`
		// 错误消息
		// 例如: 服务器内部错误
		// 必需: true
		Message string `json:"message"`
		// 响应数据
		// 必需: true
		Data interface{} `json:"data"`
	}
}
