package dto

// UserDTO represents data transfer object for User
type UserDTO struct {
	ID        string   `json:"id,omitempty"`
	Username  string   `json:"username"`
	Email     string   `json:"email"`
	NickName  string   `json:"nickName"`
	Phone     string   `json:"phone,omitempty"` // Added Phone field
	Icon      string   `json:"icon,omitempty"`
	Status    int      `json:"status"`
	Note      string   `json:"note,omitempty"`
	RoleIds   []string `json:"roleIds,omitempty"`
	CreatedAt string   `json:"createdAt,omitempty"`
	LastLogin string   `json:"lastLogin,omitempty"`
}

// UserCreateRequest represents request to create a new user
type UserCreateRequest struct {
	Username string   `json:"username" binding:"required" validate:"required,min=3,max=20"`
	Password string   `json:"password" binding:"required" validate:"required,strong_password"`
	Email    string   `json:"email" binding:"required,email" validate:"required,email"`
	Phone    string   `json:"phone" binding:"omitempty" validate:"omitempty,phone_cn"`
	NickName string   `json:"nickName" validate:"required,max=30"`
	Icon     string   `json:"icon,omitempty"`
	Note     string   `json:"note,omitempty" validate:"max=200"`
	Status   int      `json:"status" validate:"oneof=0 1"`
	RoleIds  []string `json:"roleIds,omitempty"`
}

// UserUpdateRequest represents request to update an existing user
type UserUpdateRequest struct {
	Email    string   `json:"email,omitempty" binding:"omitempty,email" validate:"omitempty,email"`
	Phone    string   `json:"phone,omitempty" validate:"omitempty,phone_cn"`
	NickName string   `json:"nickName,omitempty" validate:"omitempty,max=30"`
	Icon     string   `json:"icon,omitempty"`
	Note     string   `json:"note,omitempty" validate:"omitempty,max=200"`
	Status   *int     `json:"status,omitempty" validate:"omitempty,oneof=0 1"`
	RoleIds  []string `json:"roleIds,omitempty"`
}

// UserLoginRequest represents request to login
type UserLoginRequest struct {
	Username string `json:"username" binding:"required" validate:"required"`
	Password string `json:"password" binding:"required" validate:"required"`
}

// UserLoginResponse represents response after successful login
type UserLoginResponse struct {
	Token     string  `json:"token"`
	TokenHead string  `json:"tokenHead"`
	ExpireAt  string  `json:"expireAt"`
	User      UserDTO `json:"user"`
}

// UserChangePasswordRequest represents request to change password
type UserChangePasswordRequest struct {
	OldPassword string `json:"oldPassword" binding:"required" validate:"required"`
	NewPassword string `json:"newPassword" binding:"required" validate:"required,strong_password"`
}

// ListResponse represents a paginated response
type ListResponse struct {
	Total int64       `json:"total"`
	List  interface{} `json:"list"`
}

// RoleDTO represents data transfer object for Role
type RoleDTO struct {
	ID          string `json:"id,omitempty"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	CreatedAt   string `json:"createdAt,omitempty"`
}

// RoleCreateRequest represents request to create a new role
type RoleCreateRequest struct {
	Name        string `json:"name" binding:"required" validate:"required,min=2,max=50"`
	Description string `json:"description,omitempty" validate:"max=200"`
}

// RoleUpdateRequest represents request to update an existing role
type RoleUpdateRequest struct {
	Name        string `json:"name,omitempty" validate:"omitempty,min=2,max=50"`
	Description string `json:"description,omitempty" validate:"omitempty,max=200"`
}

// AssignRoleRequest represents request to assign role to user
type AssignRoleRequest struct {
	UserID  string   `json:"userId" binding:"required" validate:"required"`
	RoleIds []string `json:"roleIds" binding:"required" validate:"required,min=1"`
}
