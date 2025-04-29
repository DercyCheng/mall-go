package repository

import (
	"context"

	"mall-go/services/user-service/domain/model"
)

// UserRepository defines the interface for user data storage operations
type UserRepository interface {
	// Save persists a user to the database
	Save(ctx context.Context, user *model.User) error

	// FindByID finds a user by their unique ID
	FindByID(ctx context.Context, id string) (*model.User, error)

	// FindByUsername finds a user by their username
	FindByUsername(ctx context.Context, username string) (*model.User, error)

	// FindByEmail finds a user by their email
	FindByEmail(ctx context.Context, email string) (*model.User, error)

	// Update updates a user's information
	Update(ctx context.Context, user *model.User) error

	// Delete removes a user from the database
	Delete(ctx context.Context, id string) error

	// List returns a list of users with pagination
	List(ctx context.Context, page, pageSize int) ([]*model.User, int64, error)

	// Search searches for users based on criteria
	Search(ctx context.Context, query string, page, pageSize int) ([]*model.User, int64, error)
}

// RoleRepository defines the interface for role data storage operations
type RoleRepository interface {
	// Save persists a role to the database
	Save(ctx context.Context, role *model.Role) error

	// FindByID finds a role by its unique ID
	FindByID(ctx context.Context, id string) (*model.Role, error)

	// FindByName finds a role by its name
	FindByName(ctx context.Context, name string) (*model.Role, error)

	// Update updates a role's information
	Update(ctx context.Context, role *model.Role) error

	// Delete removes a role from the database
	Delete(ctx context.Context, id string) error

	// List returns a list of roles
	List(ctx context.Context) ([]*model.Role, error)

	// GetUserRoles returns all roles for a given user
	GetUserRoles(ctx context.Context, userID string) ([]*model.Role, error)

	// AssignRoleToUser assigns a role to a user
	AssignRoleToUser(ctx context.Context, userID, roleID string) error

	// RevokeRoleFromUser removes a role from a user
	RevokeRoleFromUser(ctx context.Context, userID, roleID string) error
}
