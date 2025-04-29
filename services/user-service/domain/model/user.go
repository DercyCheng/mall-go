package model

import (
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// UserStatus represents the status of a user
type UserStatus int

const (
	UserStatusActive   UserStatus = 1
	UserStatusInactive UserStatus = 0
)

// User represents a user domain entity
type User struct {
	ID        string     `json:"id"`
	Username  string     `json:"username"`
	Password  string     `json:"-"` // Password is never exposed in JSON
	Email     string     `json:"email"`
	NickName  string     `json:"nickName"`
	Phone     string     `json:"phone"` // Added Phone field
	Icon      string     `json:"icon"`
	Status    UserStatus `json:"status"`
	Version   int        `json:"version"` // Added Version field for optimistic locking
	Note      string     `json:"note"`
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt time.Time  `json:"updatedAt"`
	LastLogin time.Time  `json:"lastLogin"`
	Roles     []Role     `json:"roles"`
}

// Role represents a user role domain entity
type Role struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

// NewUser creates a new user with encrypted password
func NewUser(username, password, email, nickName string) (*User, error) {
	if username == "" || password == "" || email == "" {
		return nil, errors.New("username, password and email are required")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	return &User{
		Username:  username,
		Password:  string(hashedPassword),
		Email:     email,
		NickName:  nickName,
		Phone:     "", // Initialize Phone field
		Status:    UserStatusActive,
		Version:   1, // Initialize Version field for optimistic locking
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}

// VerifyPassword checks if the provided password matches the stored hash
func (u *User) VerifyPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}

// ChangePassword changes user's password
func (u *User) ChangePassword(currentPassword, newPassword string) error {
	if !u.VerifyPassword(currentPassword) {
		return errors.New("current password is incorrect")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	u.Password = string(hashedPassword)
	u.UpdatedAt = time.Now()
	u.Version++ // Increment version for optimistic locking
	return nil
}

// Activate activates the user
func (u *User) Activate() {
	u.Status = UserStatusActive
	u.UpdatedAt = time.Now()
	u.Version++ // Increment version for optimistic locking
}

// Deactivate deactivates the user
func (u *User) Deactivate() {
	u.Status = UserStatusInactive
	u.UpdatedAt = time.Now()
	u.Version++ // Increment version for optimistic locking
}

// UpdateProfile updates user profile information
func (u *User) UpdateProfile(nickName, email, phone, icon string) {
	if nickName != "" {
		u.NickName = nickName
	}

	if email != "" {
		u.Email = email
	}

	if phone != "" {
		u.Phone = phone
	}

	if icon != "" {
		u.Icon = icon
	}

	u.UpdatedAt = time.Now()
	u.Version++ // Increment version for optimistic locking
}

// RecordLogin records a login event
func (u *User) RecordLogin() {
	u.LastLogin = time.Now()
}

// HasRole checks if the user has a specific role
func (u *User) HasRole(roleName string) bool {
	for _, role := range u.Roles {
		if role.Name == roleName {
			return true
		}
	}
	return false
}

// AddRole adds a role to the user if it doesn't already exist
func (u *User) AddRole(role Role) {
	// Check if user already has this role
	for _, r := range u.Roles {
		if r.ID == role.ID {
			return
		}
	}

	u.Roles = append(u.Roles, role)
	u.UpdatedAt = time.Now()
	u.Version++ // Increment version for optimistic locking
}

// RemoveRole removes a role from the user
func (u *User) RemoveRole(roleID string) {
	var newRoles []Role
	for _, r := range u.Roles {
		if r.ID != roleID {
			newRoles = append(newRoles, r)
		}
	}

	u.Roles = newRoles
	u.UpdatedAt = time.Now()
	u.Version++ // Increment version for optimistic locking
}
