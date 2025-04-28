package dto

import (
	"errors"
	"time"

	"mall-go/services/user-service/domain/model"
)

// 请求DTO

// RegisterUserRequest 用户注册请求
type RegisterUserRequest struct {
	Username  string `json:"username" binding:"required,min=3,max=20"`
	Password  string `json:"password" binding:"required,min=6"`
	Email     string `json:"email" binding:"required,email"`
	Phone     string `json:"phone" binding:"required"`
	NickName  string `json:"nickName"`
	Gender    int    `json:"gender" binding:"oneof=0 1 2"` // 0-未知 1-男 2-女
	Birthday  string `json:"birthday" binding:"omitempty,datetime=2006-01-02"`
	City      string `json:"city"`
	Job       string `json:"job"`
}

// LoginRequest 用户登录请求
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// UpdateUserRequest 更新用户信息请求
type UpdateUserRequest struct {
	Email    string `json:"email" binding:"omitempty,email"`
	Phone    string `json:"phone"`
	NickName string `json:"nickName"`
	Icon     string `json:"icon"`
	Gender   int    `json:"gender" binding:"omitempty,oneof=0 1 2"`
	Birthday string `json:"birthday" binding:"omitempty,datetime=2006-01-02"`
	City     string `json:"city"`
	Job      string `json:"job"`
}

// UpdatePasswordRequest 修改密码请求
type UpdatePasswordRequest struct {
	OldPassword string `json:"oldPassword" binding:"required"`
	NewPassword string `json:"newPassword" binding:"required,min=6"`
}

// UpdateStatusRequest 更新用户状态请求
type UpdateStatusRequest struct {
	Status string `json:"status" binding:"required,oneof=active inactive locked"`
}

// AssignRoleRequest 分配角色请求
type AssignRoleRequest struct {
	RoleIDs []string `json:"roleIds" binding:"required"`
}

// SearchUserRequest 搜索用户请求
type SearchUserRequest struct {
	Keyword  string `json:"keyword" form:"keyword"`
	Page     int    `json:"page" form:"page,default=1" binding:"min=1"`
	PageSize int    `json:"pageSize" form:"pageSize,default=10" binding:"min=1,max=100"`
}

// 响应DTO

// UserResponse 用户信息响应
type UserResponse struct {
	ID               string    `json:"id"`
	Username         string    `json:"username"`
	Email            string    `json:"email"`
	Phone            string    `json:"phone"`
	NickName         string    `json:"nickName"`
	Icon             string    `json:"icon"`
	Gender           int       `json:"gender"`
	Birthday         *string   `json:"birthday,omitempty"`
	City             string    `json:"city"`
	Job              string    `json:"job"`
	Status           string    `json:"status"`
	IntegrationPoints int      `json:"integrationPoints"`
	Roles            []RoleDTO `json:"roles,omitempty"`
	CreatedAt        time.Time `json:"createdAt"`
	LastLoginAt      *string   `json:"lastLoginAt,omitempty"`
}

// RoleDTO 角色信息
type RoleDTO struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Status      string `json:"status"`
}

// PermissionDTO 权限信息
type PermissionDTO struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Value  string `json:"value"`
	Type   string `json:"type"`
	Status string `json:"status"`
}

// LoginResponse 登录响应
type LoginResponse struct {
	Token     string       `json:"token"`
	TokenType string       `json:"tokenType"`
	ExpiresIn int          `json:"expiresIn"` // 过期时间(秒)
	User      UserResponse `json:"user"`
}

// UserListResponse 用户列表响应
type UserListResponse struct {
	List     []UserResponse `json:"list"`
	Total    int64          `json:"total"`
	Page     int            `json:"page"`
	PageSize int            `json:"pageSize"`
}

// RoleListResponse 角色列表响应
type RoleListResponse struct {
	List     []RoleDTO `json:"list"`
	Total    int64     `json:"total"`
	Page     int       `json:"page"`
	PageSize int       `json:"pageSize"`
}

// PermissionListResponse 权限列表响应
type PermissionListResponse struct {
	List     []PermissionDTO `json:"list"`
	Total    int64           `json:"total"`
	Page     int             `json:"page"`
	PageSize int             `json:"pageSize"`
}

// AssignPermissionsToRoleRequest 分配权限到角色请求
type AssignPermissionsToRoleRequest struct {
	PermissionIDs []string `json:"permissionIds" binding:"required"`
}

// Validate 验证请求
func (req *AssignPermissionsToRoleRequest) Validate() error {
	if len(req.PermissionIDs) == 0 {
		return errors.New("permission IDs cannot be empty")
	}
	return nil
}