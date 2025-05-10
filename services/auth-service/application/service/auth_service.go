package service

import (
	"context"
	"errors"
	"time"

	"mall-go/pkg/auth"
	"mall-go/services/auth-service/application/dto"
	"mall-go/services/auth-service/domain/model"
	"mall-go/services/auth-service/domain/repository"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// AuthService defines the authentication service
type AuthService interface {
	// Login authenticates a user and returns a token
	Login(ctx context.Context, req dto.LoginRequest) (*dto.LoginResponse, error)

	// Register creates a new user account
	Register(ctx context.Context, req dto.RegisterRequest) (*dto.RegisterResponse, error)

	// ValidateToken verifies a JWT token
	ValidateToken(ctx context.Context, req dto.ValidateTokenRequest) (*dto.ValidateTokenResponse, error)

	// RefreshToken refreshes an existing JWT token
	RefreshToken(ctx context.Context, req dto.RefreshTokenRequest) (*dto.RefreshTokenResponse, error)

	// Logout invalidates a JWT token
	Logout(ctx context.Context, req dto.LogoutRequest) (*dto.LogoutResponse, error)
}

// authServiceImpl is the implementation of AuthService
type authServiceImpl struct {
	userRepository  repository.UserRepository
	tokenRepository repository.TokenRepository
	jwtSecret       string
	tokenExpiry     int64
	refreshExpiry   int64
}

// NewAuthService creates a new instance of AuthService
func NewAuthService(
	userRepository repository.UserRepository,
	tokenRepository repository.TokenRepository,
	jwtSecret string,
	tokenExpiry int64,
	refreshExpiry int64,
) AuthService {
	return &authServiceImpl{
		userRepository:  userRepository,
		tokenRepository: tokenRepository,
		jwtSecret:       jwtSecret,
		tokenExpiry:     tokenExpiry,
		refreshExpiry:   refreshExpiry,
	}
}

// Login authenticates a user and returns a token
func (s *authServiceImpl) Login(ctx context.Context, req dto.LoginRequest) (*dto.LoginResponse, error) {
	// Find user by username
	user, err := s.userRepository.FindByUsername(ctx, req.Username)
	if err != nil {
		return &dto.LoginResponse{
			Success: false,
			Message: "Invalid username or password",
		}, nil
	}

	// Check if user is locked
	if user.IsLocked() {
		return &dto.LoginResponse{
			Success: false,
			Message: "Account is locked. Please try again later or contact support.",
		}, nil
	}

	// Verify password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	
	// Record login attempt
	loginHistory := model.LoginHistory{
		ID:        uuid.New().String(),
		UserID:    user.ID,
		IP:        req.ClientIP,
		UserAgent: req.UserAgent,
		Status:    err == nil,
		Message:   "Login attempt",
		CreatedAt: time.Now(),
	}
	
	// Record login history
	_ = s.userRepository.CreateLoginHistory(ctx, loginHistory)
	
	if err != nil {
		// Record failed login attempt
		user.RecordLoginAttempt(false)
		_ = s.userRepository.Update(ctx, user)
		
		return &dto.LoginResponse{
			Success: false,
			Message: "Invalid username or password",
		}, nil
	}

	// Record successful login
	user.RecordLoginAttempt(true)
	_ = s.userRepository.Update(ctx, user)

	// Generate JWT token
	expiresAt := time.Now().Add(time.Duration(s.tokenExpiry) * time.Second)
	tokenString, err := auth.GenerateToken(user.ID, user.Username, convertRoles(user.Roles), s.jwtSecret, expiresAt)
	if err != nil {
		return nil, err
	}
	
	// Generate refresh token
	refreshToken := uuid.New().String()
	
	// Store tokens
	err = s.tokenRepository.StoreToken(ctx, user.ID, "access", tokenString, expiresAt.Unix())
	if err != nil {
		return nil, err
	}
	
	err = s.tokenRepository.StoreRefreshToken(ctx, user.ID, refreshToken, tokenString, time.Now().Add(time.Duration(s.refreshExpiry)*time.Second).Unix())
	if err != nil {
		return nil, err
	}

	// Convert roles to string slice
	roleStrings := make([]string, len(user.Roles))
	for i, role := range user.Roles {
		roleStrings[i] = role.Code
	}

	// Create response
	response := &dto.LoginResponse{
		Success:      true,
		Message:      "Login successful",
		Token:        tokenString,
		RefreshToken: refreshToken,
		ExpiresAt:    expiresAt.Unix(),
		UserInfo: dto.UserInfo{
			UserID:      user.ID,
			Username:    user.Username,
			Email:       user.Email,
			Phone:       user.Phone,
			UserType:    int(user.UserType),
			Roles:       roleStrings,
			LastLoginAt: user.LastLoginAt.Unix(),
		},
	}

	return response, nil
}

// Register creates a new user account
func (s *authServiceImpl) Register(ctx context.Context, req dto.RegisterRequest) (*dto.RegisterResponse, error) {
	// Check if username exists
	existingUser, err := s.userRepository.FindByUsername(ctx, req.Username)
	if err == nil && existingUser != nil {
		return &dto.RegisterResponse{
			Success: false,
			Message: "Username already exists",
		}, nil
	}

	// Check if email exists
	existingUser, err = s.userRepository.FindByEmail(ctx, req.Email)
	if err == nil && existingUser != nil {
		return &dto.RegisterResponse{
			Success: false,
			Message: "Email already exists",
		}, nil
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// Create user
	userType := model.UserTypeCustomer
	if req.UserType == int(model.UserTypeAdmin) {
		userType = model.UserTypeAdmin
	}
	
	user := &model.User{
		ID:           uuid.New().String(),
		Username:     req.Username,
		Password:     string(hashedPassword),
		Email:        req.Email,
		Phone:        req.Phone,
		UserType:     userType,
		Status:       model.UserStatusActive,
		Roles:        []model.Role{},
		LastLoginAt:  time.Time{},
		LoginAttempts: 0,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	// Validate user
	if err := user.Validate(); err != nil {
		return &dto.RegisterResponse{
			Success: false,
			Message: err.Error(),
		}, nil
	}

	// Save user
	if err := s.userRepository.Create(ctx, user); err != nil {
		return nil, err
	}

	// Create response
	return &dto.RegisterResponse{
		Success: true,
		Message: "Registration successful",
		UserID:  user.ID,
	}, nil
}

// ValidateToken verifies a JWT token
func (s *authServiceImpl) ValidateToken(ctx context.Context, req dto.ValidateTokenRequest) (*dto.ValidateTokenResponse, error) {
	// Parse token
	claims, err := auth.ParseToken(req.Token, s.jwtSecret)
	if err != nil {
		return &dto.ValidateTokenResponse{
			Valid:   false,
			Message: "Invalid token",
		}, nil
	}

	// Check if token exists in repository
	userID, err := s.tokenRepository.ValidateToken(ctx, "access", req.Token)
	if err != nil || userID == "" {
		return &dto.ValidateTokenResponse{
			Valid:   false,
			Message: "Token has been invalidated",
		}, nil
	}

	// Get user details
	user, err := s.userRepository.FindByID(ctx, claims.UserID)
	if err != nil {
		return &dto.ValidateTokenResponse{
			Valid:   false,
			Message: "User not found",
		}, nil
	}

	// Check if user is still active
	if user.Status != model.UserStatusActive {
		return &dto.ValidateTokenResponse{
			Valid:   false,
			Message: "User is not active",
		}, nil
	}

	// Extract roles as string slice
	roleStrings := make([]string, len(user.Roles))
	for i, role := range user.Roles {
		roleStrings[i] = role.Code
	}

	// Convert claims to map
	claimMap := make(map[string]string)
	for key, value := range claims.CustomClaims {
		if str, ok := value.(string); ok {
			claimMap[key] = str
		}
	}

	return &dto.ValidateTokenResponse{
		Valid:    true,
		Message:  "Token is valid",
		UserID:   claims.UserID,
		Username: claims.Username,
		Roles:    roleStrings,
		Claims:   claimMap,
	}, nil
}

// RefreshToken refreshes an existing JWT token
func (s *authServiceImpl) RefreshToken(ctx context.Context, req dto.RefreshTokenRequest) (*dto.RefreshTokenResponse, error) {
	// Get user ID from refresh token
	userID, err := s.tokenRepository.GetUserIDByRefreshToken(ctx, req.RefreshToken)
	if err != nil || userID == "" {
		return &dto.RefreshTokenResponse{
			Success: false,
			Message: "Invalid refresh token",
		}, nil
	}

	// Get user details
	user, err := s.userRepository.FindByID(ctx, userID)
	if err != nil {
		return &dto.RefreshTokenResponse{
			Success: false,
			Message: "User not found",
		}, nil
	}

	// Check if user is still active
	if user.Status != model.UserStatusActive {
		return &dto.RefreshTokenResponse{
			Success: false,
			Message: "User is not active",
		}, nil
	}

	// Get user roles
	roles, err := s.userRepository.FindUserRoles(ctx, userID)
	if err != nil {
		return nil, err
	}
	user.Roles = roles

	// Invalidate old refresh token
	_ = s.tokenRepository.InvalidateRefreshToken(ctx, req.RefreshToken)

	// Generate new JWT token
	expiresAt := time.Now().Add(time.Duration(s.tokenExpiry) * time.Second)
	tokenString, err := auth.GenerateToken(user.ID, user.Username, convertRoles(user.Roles), s.jwtSecret, expiresAt)
	if err != nil {
		return nil, err
	}
	
	// Generate new refresh token
	refreshToken := uuid.New().String()
	
	// Store new tokens
	err = s.tokenRepository.StoreToken(ctx, user.ID, "access", tokenString, expiresAt.Unix())
	if err != nil {
		return nil, err
	}
	
	err = s.tokenRepository.StoreRefreshToken(ctx, user.ID, refreshToken, tokenString, time.Now().Add(time.Duration(s.refreshExpiry)*time.Second).Unix())
	if err != nil {
		return nil, err
	}

	return &dto.RefreshTokenResponse{
		Success:      true,
		Message:      "Token refreshed",
		Token:        tokenString,
		RefreshToken: refreshToken,
		ExpiresAt:    expiresAt.Unix(),
	}, nil
}

// Logout invalidates a JWT token
func (s *authServiceImpl) Logout(ctx context.Context, req dto.LogoutRequest) (*dto.LogoutResponse, error) {
	// Invalidate token
	err := s.tokenRepository.InvalidateToken(ctx, "access", req.Token)
	if err != nil {
		return nil, err
	}

	return &dto.LogoutResponse{
		Success: true,
		Message: "Logout successful",
	}, nil
}

// Helper function to convert domain roles to string array
func convertRoles(roles []model.Role) []string {
	result := make([]string, len(roles))
	for i, role := range roles {
		result[i] = role.Code
	}
	return result
}
