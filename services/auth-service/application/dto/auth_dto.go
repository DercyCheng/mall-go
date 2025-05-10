package dto

// LoginRequest represents the login request data transfer object
type LoginRequest struct {
	Username  string `json:"username" binding:"required"`
	Password  string `json:"password" binding:"required"`
	ClientIP  string `json:"client_ip"`
	UserAgent string `json:"user_agent"`
}

// LoginResponse represents the login response data transfer object
type LoginResponse struct {
	Success     bool     `json:"success"`
	Message     string   `json:"message"`
	Token       string   `json:"token,omitempty"`
	RefreshToken string  `json:"refresh_token,omitempty"`
	ExpiresAt   int64    `json:"expires_at,omitempty"`
	UserInfo    UserInfo `json:"user_info,omitempty"`
}

// RegisterRequest represents the registration request data transfer object
type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=4,max=32"`
	Password string `json:"password" binding:"required,min=6"`
	Email    string `json:"email" binding:"required,email"`
	Phone    string `json:"phone"`
	UserType int    `json:"user_type"`
}

// RegisterResponse represents the registration response data transfer object
type RegisterResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	UserID  string `json:"user_id,omitempty"`
}

// ValidateTokenRequest represents the token validation request data transfer object
type ValidateTokenRequest struct {
	Token string `json:"token" binding:"required"`
}

// ValidateTokenResponse represents the token validation response data transfer object
type ValidateTokenResponse struct {
	Valid    bool              `json:"valid"`
	Message  string            `json:"message"`
	UserID   string            `json:"user_id,omitempty"`
	Username string            `json:"username,omitempty"`
	Roles    []string          `json:"roles,omitempty"`
	Claims   map[string]string `json:"claims,omitempty"`
}

// RefreshTokenRequest represents the token refresh request data transfer object
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

// RefreshTokenResponse represents the token refresh response data transfer object
type RefreshTokenResponse struct {
	Success      bool   `json:"success"`
	Message      string `json:"message"`
	Token        string `json:"token,omitempty"`
	RefreshToken string `json:"refresh_token,omitempty"`
	ExpiresAt    int64  `json:"expires_at,omitempty"`
}

// LogoutRequest represents the logout request data transfer object
type LogoutRequest struct {
	Token  string `json:"token" binding:"required"`
	UserID string `json:"user_id"`
}

// LogoutResponse represents the logout response data transfer object
type LogoutResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// UserInfo represents simplified user information data transfer object
type UserInfo struct {
	UserID       string   `json:"user_id"`
	Username     string   `json:"username"`
	Email        string   `json:"email"`
	Phone        string   `json:"phone,omitempty"`
	UserType     int      `json:"user_type"`
	Roles        []string `json:"roles,omitempty"`
	LastLoginAt  int64    `json:"last_login_time,omitempty"`
}
