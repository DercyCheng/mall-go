package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"

	"mall-go/services/user-service/application/assembler"
	"mall-go/services/user-service/application/dto"
	"mall-go/services/user-service/domain/model"
	"mall-go/services/user-service/domain/repository"
)

// UserServiceImpl implements the UserService interface
type UserServiceImpl struct {
	userRepo    repository.UserRepository
	roleRepo    repository.RoleRepository
	jwtSecret   string
	tokenExpiry time.Duration
}

// NewUserService creates a new UserService implementation
func NewUserService(
	userRepo repository.UserRepository,
	roleRepo repository.RoleRepository,
	jwtSecret string,
	tokenExpiry time.Duration,
) UserService {
	return &UserServiceImpl{
		userRepo:    userRepo,
		roleRepo:    roleRepo,
		jwtSecret:   jwtSecret,
		tokenExpiry: tokenExpiry,
	}
}

// Register creates a new user
func (s *UserServiceImpl) Register(ctx context.Context, req dto.UserCreateRequest) (string, error) {
	// Check if username already exists
	existingUser, err := s.userRepo.FindByUsername(ctx, req.Username)
	if err != nil {
		return "", fmt.Errorf("error checking existing user: %w", err)
	}

	if existingUser != nil {
		return "", errors.New("username already exists")
	}

	// Check if email already exists
	existingEmail, err := s.userRepo.FindByEmail(ctx, req.Email)
	if err != nil {
		return "", fmt.Errorf("error checking existing email: %w", err)
	}

	if existingEmail != nil {
		return "", errors.New("email already in use")
	}

	// Create new user domain model
	user, err := assembler.ToUserModel(req)
	if err != nil {
		return "", fmt.Errorf("error creating user model: %w", err)
	}

	// Set icon if provided
	if req.Icon != "" {
		user.Icon = req.Icon
	}

	// Set note if provided
	if req.Note != "" {
		user.Note = req.Note
	}

	// Set status if provided
	if req.Status != 0 {
		user.Status = model.UserStatus(req.Status)
	}

	// Create user in repository
	if err := s.userRepo.Save(ctx, user); err != nil {
		return "", fmt.Errorf("error saving user: %w", err)
	}

	// Assign roles if provided
	if len(req.RoleIds) > 0 {
		for _, roleID := range req.RoleIds {
			if err := s.roleRepo.AssignRoleToUser(ctx, user.ID, roleID); err != nil {
				return "", fmt.Errorf("error assigning role: %w", err)
			}
		}
	}

	return user.ID, nil
}

// Login authenticates a user and returns a JWT token
func (s *UserServiceImpl) Login(ctx context.Context, req dto.UserLoginRequest) (*dto.UserLoginResponse, error) {
	// Find user by username
	user, err := s.userRepo.FindByUsername(ctx, req.Username)
	if err != nil {
		return nil, fmt.Errorf("error finding user: %w", err)
	}

	if user == nil {
		return nil, errors.New("invalid username or password")
	}

	// Verify password
	if !user.VerifyPassword(req.Password) {
		return nil, errors.New("invalid username or password")
	}

	// Check if user is active
	if user.Status != model.UserStatusActive {
		return nil, errors.New("user account is inactive")
	}

	// Record login
	user.RecordLogin()
	if err := s.userRepo.Update(ctx, user); err != nil {
		return nil, fmt.Errorf("error updating user login time: %w", err)
	}

	// Generate token
	expiryTime := time.Now().Add(s.tokenExpiry)
	claims := jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(expiryTime),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ID:        uuid.New().String(),
		Subject:   user.ID,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(s.jwtSecret))
	if err != nil {
		return nil, fmt.Errorf("error signing token: %w", err)
	}

	// Create response
	userDTO := assembler.ToUserDTO(user)
	response := &dto.UserLoginResponse{
		Token:     signedToken,
		TokenHead: "Bearer",
		ExpireAt:  expiryTime.Format(time.RFC3339),
		User:      userDTO,
	}

	return response, nil
}

// GetUserInfo retrieves user information by ID
func (s *UserServiceImpl) GetUserInfo(ctx context.Context, id string) (*dto.UserDTO, error) {
	user, err := s.userRepo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("error finding user: %w", err)
	}

	if user == nil {
		return nil, errors.New("user not found")
	}

	userDTO := assembler.ToUserDTO(user)
	return &userDTO, nil
}

// GetUserByUsername retrieves user information by username
func (s *UserServiceImpl) GetUserByUsername(ctx context.Context, username string) (*dto.UserDTO, error) {
	user, err := s.userRepo.FindByUsername(ctx, username)
	if err != nil {
		return nil, fmt.Errorf("error finding user: %w", err)
	}

	if user == nil {
		return nil, errors.New("user not found")
	}

	userDTO := assembler.ToUserDTO(user)
	return &userDTO, nil
}

// UpdateUser updates user information
func (s *UserServiceImpl) UpdateUser(ctx context.Context, id string, req dto.UserUpdateRequest) error {
	user, err := s.userRepo.FindByID(ctx, id)
	if err != nil {
		return fmt.Errorf("error finding user: %w", err)
	}

	if user == nil {
		return errors.New("user not found")
	}

	// Update profile fields
	if req.Email != "" || req.NickName != "" || req.Phone != "" || req.Icon != "" {
		user.UpdateProfile(req.NickName, req.Email, req.Phone, req.Icon)
	}

	// Update note if provided
	if req.Note != "" {
		user.Note = req.Note
	}

	// Update status if provided
	if req.Status != nil {
		user.Status = model.UserStatus(*req.Status)
	}

	// Save user updates
	if err := s.userRepo.Update(ctx, user); err != nil {
		return fmt.Errorf("error updating user: %w", err)
	}

	// Handle role assignments if provided
	if len(req.RoleIds) > 0 {
		// Get current user roles to determine what needs to be added or removed
		currentRoles, err := s.roleRepo.GetUserRoles(ctx, id)
		if err != nil {
			return fmt.Errorf("error getting user roles: %w", err)
		}

		// Create map of current role IDs for quick lookup
		currentRoleMap := make(map[string]bool)
		for _, role := range currentRoles {
			currentRoleMap[role.ID] = true
		}

		// Create map of requested role IDs for quick lookup
		requestedRoleMap := make(map[string]bool)
		for _, roleID := range req.RoleIds {
			requestedRoleMap[roleID] = true

			// Add role if not present
			if !currentRoleMap[roleID] {
				if err := s.roleRepo.AssignRoleToUser(ctx, id, roleID); err != nil {
					return fmt.Errorf("error assigning role: %w", err)
				}
			}
		}

		// Remove roles that are no longer requested
		for _, role := range currentRoles {
			if !requestedRoleMap[role.ID] {
				if err := s.roleRepo.RevokeRoleFromUser(ctx, id, role.ID); err != nil {
					return fmt.Errorf("error revoking role: %w", err)
				}
			}
		}
	}

	return nil
}

// DeleteUser deletes a user
func (s *UserServiceImpl) DeleteUser(ctx context.Context, id string) error {
	// Check if user exists
	user, err := s.userRepo.FindByID(ctx, id)
	if err != nil {
		return fmt.Errorf("error finding user: %w", err)
	}

	if user == nil {
		return errors.New("user not found")
	}

	// Delete user
	return s.userRepo.Delete(ctx, id)
}

// ChangePassword changes a user's password
func (s *UserServiceImpl) ChangePassword(ctx context.Context, id string, req dto.UserChangePasswordRequest) error {
	// Check if user exists
	user, err := s.userRepo.FindByID(ctx, id)
	if err != nil {
		return fmt.Errorf("error finding user: %w", err)
	}

	if user == nil {
		return errors.New("user not found")
	}

	// Change password
	if err := user.ChangePassword(req.OldPassword, req.NewPassword); err != nil {
		return err
	}

	// Save user updates
	return s.userRepo.Update(ctx, user)
}

// ListUsers retrieves a paginated list of users
func (s *UserServiceImpl) ListUsers(ctx context.Context, page, pageSize int) (*dto.ListResponse, error) {
	users, total, err := s.userRepo.List(ctx, page, pageSize)
	if err != nil {
		return nil, fmt.Errorf("error listing users: %w", err)
	}

	userDTOs := assembler.ToUserDTOList(users)
	return &dto.ListResponse{
		Total: total,
		List:  userDTOs,
	}, nil
}

// SearchUsers searches for users based on criteria
func (s *UserServiceImpl) SearchUsers(ctx context.Context, query string, page, pageSize int) (*dto.ListResponse, error) {
	users, total, err := s.userRepo.Search(ctx, query, page, pageSize)
	if err != nil {
		return nil, fmt.Errorf("error searching users: %w", err)
	}

	userDTOs := assembler.ToUserDTOList(users)
	return &dto.ListResponse{
		Total: total,
		List:  userDTOs,
	}, nil
}
