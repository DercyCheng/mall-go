package repository

import (
	"context"

	"mall-go/services/auth-service/domain/model"
)

// UserRepository defines the repository interface for user operations
type UserRepository interface {
	// FindByID finds a user by ID
	FindByID(ctx context.Context, id string) (*model.User, error)

	// FindByUsername finds a user by username
	FindByUsername(ctx context.Context, username string) (*model.User, error)

	// FindByEmail finds a user by email
	FindByEmail(ctx context.Context, email string) (*model.User, error)

	// Create creates a new user
	Create(ctx context.Context, user *model.User) error

	// Update updates an existing user
	Update(ctx context.Context, user *model.User) error

	// Delete deletes a user by ID
	Delete(ctx context.Context, id string) error

	// FindUserRoles finds all roles for a user
	FindUserRoles(ctx context.Context, userID string) ([]model.Role, error)

	// UpdateUserRoles updates user roles
	UpdateUserRoles(ctx context.Context, userID string, roleIDs []string) error

	// CreateLoginHistory creates a login history record
	CreateLoginHistory(ctx context.Context, history model.LoginHistory) error

	// GetLoginHistory gets a user's login history
	GetLoginHistory(ctx context.Context, userID string, limit, offset int) ([]model.LoginHistory, int64, error)
}

// TokenRepository defines the repository interface for token operations
type TokenRepository interface {
	// StoreToken stores a token with expiration
	StoreToken(ctx context.Context, userID, tokenType, token string, expiration int64) error

	// ValidateToken validates if a token exists and is valid
	ValidateToken(ctx context.Context, tokenType, token string) (string, error)

	// InvalidateToken invalidates a token
	InvalidateToken(ctx context.Context, tokenType, token string) error

	// InvalidateUserTokens invalidates all tokens for a user
	InvalidateUserTokens(ctx context.Context, userID string) error

	// StoreRefreshToken stores a refresh token with expiration
	StoreRefreshToken(ctx context.Context, userID, refreshToken, accessToken string, expiration int64) error

	// GetUserIDByRefreshToken gets a user ID by refresh token
	GetUserIDByRefreshToken(ctx context.Context, refreshToken string) (string, error)

	// InvalidateRefreshToken invalidates a refresh token
	InvalidateRefreshToken(ctx context.Context, refreshToken string) error
}
