package model

import (
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// User 用户聚合根
type User struct {
	ID            string
	Username      string
	Password      string
	Email         string
	Phone         string
	NickName      string
	Icon          string
	Gender        int
	Birthday      *time.Time
	City          string
	Job           string
	Status        UserStatus
	IntegrationPoints int
	Roles         []Role
	CreatedAt     time.Time
	UpdatedAt     time.Time
	LastLoginAt   *time.Time
}

// UserStatus 用户状态
type UserStatus string

const (
	UserStatusActive   UserStatus = "active"
	UserStatusInactive UserStatus = "inactive"
	UserStatusLocked   UserStatus = "locked"
)

// Role 用户角色
type Role struct {
	ID          string
	Name        string
	Description string
	Status      RoleStatus
	Sort        int
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// RoleStatus 角色状态
type RoleStatus string

const (
	RoleStatusActive   RoleStatus = "active"
	RoleStatusInactive RoleStatus = "inactive"
)

// Permission 权限实体
type Permission struct {
	ID          string
	Name        string
	Value       string
	Type        PermissionType
	Status      PermissionStatus
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// PermissionType 权限类型
type PermissionType string

const (
	PermissionTypeMenu     PermissionType = "menu"
	PermissionTypeButton   PermissionType = "button"
	PermissionTypeApi      PermissionType = "api"
)

// PermissionStatus 权限状态
type PermissionStatus string

const (
	PermissionStatusActive   PermissionStatus = "active"
	PermissionStatusInactive PermissionStatus = "inactive"
)

// 领域行为

// VerifyPassword 验证密码
func (u *User) VerifyPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}

// SetPassword 设置密码(哈希)
func (u *User) SetPassword(password string) error {
	if len(password) < 6 {
		return errors.New("password must be at least 6 characters")
	}
	
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	
	u.Password = string(hashedPassword)
	return nil
}

// ChangeStatus 修改用户状态
func (u *User) ChangeStatus(status UserStatus) {
	u.Status = status
	u.UpdatedAt = time.Now()
}

// UpdateLoginTime 更新登录时间
func (u *User) UpdateLoginTime() {
	now := time.Now()
	u.LastLoginAt = &now
	u.UpdatedAt = now
}

// HasRole 检查用户是否拥有指定角色
func (u *User) HasRole(roleID string) bool {
	for _, role := range u.Roles {
		if role.ID == roleID {
			return true
		}
	}
	return false
}

// AddRole 给用户添加角色
func (u *User) AddRole(role Role) {
	// 检查角色是否已存在
	if u.HasRole(role.ID) {
		return
	}
	
	u.Roles = append(u.Roles, role)
	u.UpdatedAt = time.Now()
}

// RemoveRole 移除用户的角色
func (u *User) RemoveRole(roleID string) {
	var newRoles []Role
	for _, role := range u.Roles {
		if role.ID != roleID {
			newRoles = append(newRoles, role)
		}
	}
	
	u.Roles = newRoles
	u.UpdatedAt = time.Now()
}

// IsActive 检查用户是否处于活跃状态
func (u *User) IsActive() bool {
	return u.Status == UserStatusActive
}