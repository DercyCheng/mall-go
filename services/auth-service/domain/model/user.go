package model

import (
	"errors"
	"time"
)

// User represents the user domain model
type User struct {
	ID           string
	Username     string
	Password     string
	Email        string
	Phone        string
	UserType     UserType
	Status       UserStatus
	Roles        []Role
	LastLoginAt  time.Time
	LoginAttempts int
	LockedUntil  *time.Time
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// UserType represents the type of user
type UserType int

const (
	UserTypeAdmin UserType = 1
	UserTypeCustomer UserType = 2
)

// UserStatus represents the status of a user
type UserStatus int

const (
	UserStatusActive UserStatus = 1
	UserStatusInactive UserStatus = 0
	UserStatusLocked UserStatus = 2
)

// Role represents a user role
type Role struct {
	ID         string
	Name       string
	Code       string
	Status     int
	Sort       int
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

// PasswordReset represents a password reset request
type PasswordReset struct {
	ID         string
	UserID     string
	Token      string
	ExpiredAt  time.Time
	UsedAt     *time.Time
	CreatedAt  time.Time
}

// LoginHistory represents a user login history record
type LoginHistory struct {
	ID         string
	UserID     string
	IP         string
	UserAgent  string
	Status     bool
	Message    string
	CreatedAt  time.Time
}

// Validate validates user data
func (u *User) Validate() error {
	if u.Username == "" {
		return errors.New("username is required")
	}

	if len(u.Username) < 4 || len(u.Username) > 32 {
		return errors.New("username must be between 4 and 32 characters")
	}

	if u.Password == "" {
		return errors.New("password is required")
	}

	if len(u.Password) < 6 {
		return errors.New("password must be at least 6 characters")
	}

	if u.Email == "" {
		return errors.New("email is required")
	}

	// Basic email format validation
	if !containsChar(u.Email, '@') {
		return errors.New("invalid email format")
	}

	return nil
}

// IsLocked checks if a user account is locked
func (u *User) IsLocked() bool {
	return u.Status == UserStatusLocked || (u.LockedUntil != nil && u.LockedUntil.After(time.Now()))
}

// HasRole checks if a user has a specific role
func (u *User) HasRole(roleCode string) bool {
	for _, role := range u.Roles {
		if role.Code == roleCode {
			return true
		}
	}
	return false
}

// RecordLoginAttempt records a login attempt and locks the account if too many failures
func (u *User) RecordLoginAttempt(success bool) {
	if success {
		u.LoginAttempts = 0
		u.LastLoginAt = time.Now()
		u.LockedUntil = nil
		return
	}
	
	u.LoginAttempts++
	
	// Lock account after 5 failed attempts
	if u.LoginAttempts >= 5 {
		lockUntil := time.Now().Add(time.Minute * 15)
		u.LockedUntil = &lockUntil
		u.Status = UserStatusLocked
	}
}

// Helper function to check if a string contains a character
func containsChar(s string, c byte) bool {
	for i := 0; i < len(s); i++ {
		if s[i] == c {
			return true
		}
	}
	return false
}
