// Package dto contains data transfer objects for the user service
package dto

import (
	"errors"
	"strings"
)

// CreatePermissionRequest represents a request to create a new permission
type CreatePermissionRequest struct {
	Name        string `json:"name" binding:"required"`
	Value       string `json:"value" binding:"required"`
	Description string `json:"description"`
	Type        string `json:"type" binding:"required"` // api, menu, button
	ParentID    string `json:"parentId"`
	Sort        int    `json:"sort" binding:"gte=0"`
}

// Validate validates the create permission request
func (r *CreatePermissionRequest) Validate() error {
	if strings.TrimSpace(r.Name) == "" {
		return errors.New("permission name cannot be empty")
	}
	if strings.TrimSpace(r.Value) == "" {
		return errors.New("permission value cannot be empty")
	}
	if !isValidPermissionType(r.Type) {
		return errors.New("invalid permission type, must be api, menu, or button")
	}
	return nil
}

// UpdatePermissionRequest represents a request to update an existing permission
type UpdatePermissionRequest struct {
	Name        string `json:"name"`
	Value       string `json:"value"`
	Description string `json:"description"`
	Type        string `json:"type"`
	ParentID    string `json:"parentId"`
	Sort        int    `json:"sort" binding:"gte=0"`
}

// Validate validates the update permission request
func (r *UpdatePermissionRequest) Validate() error {
	if strings.TrimSpace(r.Name) == "" && strings.TrimSpace(r.Value) == "" &&
		strings.TrimSpace(r.Type) == "" && r.Sort == 0 {
		return errors.New("at least one field must be provided")
	}
	if r.Type != "" && !isValidPermissionType(r.Type) {
		return errors.New("invalid permission type, must be api, menu, or button")
	}
	return nil
}

// UpdatePermissionStatusRequest represents a request to update a permission's status
type UpdatePermissionStatusRequest struct {
	Status string `json:"status" binding:"required"`
}

// Validate validates the update permission status request
func (r *UpdatePermissionStatusRequest) Validate() error {
	status := strings.ToLower(strings.TrimSpace(r.Status))
	if status != "0" && status != "1" && status != "enable" && status != "disable" {
		return errors.New("invalid status, must be 0, 1, enable, or disable")
	}
	return nil
}

// PermissionAssignmentRequest represents a request to assign permissions to a role
type PermissionAssignmentRequest struct {
	PermissionIDs []string `json:"permissionIds" binding:"required"`
}

// Validate validates the permission assignment request
func (r *PermissionAssignmentRequest) Validate() error {
	if len(r.PermissionIDs) == 0 {
		return errors.New("permission IDs cannot be empty")
	}
	return nil
}

// PermissionDetailResponse represents the response for a permission
type PermissionDetailResponse struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Value       string `json:"value"`
	Description string `json:"description"`
	Type        string `json:"type"`
	ParentID    string `json:"parentId"`
	Sort        int    `json:"sort"`
	Status      int    `json:"status"`
	CreatedAt   string `json:"createdAt"`
	UpdatedAt   string `json:"updatedAt"`
}

// PermissionPageResponse represents the response for a list of permissions
type PermissionPageResponse struct {
	List     []PermissionDetailResponse `json:"list"`
	Total    int64                      `json:"total"`
	Page     int                        `json:"page"`
	PageSize int                        `json:"pageSize"`
}

// Helper function to validate permission types
func isValidPermissionType(permType string) bool {
	permType = strings.ToLower(permType)
	return permType == "api" || permType == "menu" || permType == "button"
}
