package service

import (
	"context"

	"mall-go/services/user-service/application/dto"
)

// UserService defines the application service interface for user operations
type UserService interface {
	// Register creates a new user
	Register(ctx context.Context, req dto.UserCreateRequest) (string, error)

	// Login authenticates a user and returns a JWT token
	Login(ctx context.Context, req dto.UserLoginRequest) (*dto.UserLoginResponse, error)

	// GetUserInfo retrieves user information by ID
	GetUserInfo(ctx context.Context, id string) (*dto.UserDTO, error)

	// GetUserByUsername retrieves user information by username
	GetUserByUsername(ctx context.Context, username string) (*dto.UserDTO, error)

	// UpdateUser updates user information
	UpdateUser(ctx context.Context, id string, req dto.UserUpdateRequest) error

	// DeleteUser deletes a user
	DeleteUser(ctx context.Context, id string) error

	// ChangePassword changes a user's password
	ChangePassword(ctx context.Context, id string, req dto.UserChangePasswordRequest) error

	// ListUsers retrieves a paginated list of users
	ListUsers(ctx context.Context, page, pageSize int) (*dto.ListResponse, error)

	// SearchUsers searches for users based on criteria
	SearchUsers(ctx context.Context, query string, page, pageSize int) (*dto.ListResponse, error)
}

// RoleService defines the application service interface for role operations
type RoleService interface {
	// CreateRole creates a new role
	CreateRole(ctx context.Context, req dto.RoleCreateRequest) (string, error)

	// GetRole retrieves role information by ID
	GetRole(ctx context.Context, id string) (*dto.RoleDTO, error)

	// UpdateRole updates role information
	UpdateRole(ctx context.Context, id string, req dto.RoleUpdateRequest) error

	// DeleteRole deletes a role
	DeleteRole(ctx context.Context, id string) error

	// ListRoles retrieves a list of all roles
	ListRoles(ctx context.Context) ([]dto.RoleDTO, error)

	// AssignRolesToUser assigns roles to a user
	AssignRolesToUser(ctx context.Context, req dto.AssignRoleRequest) error

	// GetUserRoles retrieves roles for a user
	GetUserRoles(ctx context.Context, userID string) ([]dto.RoleDTO, error)
}
