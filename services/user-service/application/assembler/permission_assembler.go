// Package assembler provides mapping between domain models and DTOs
package assembler

import (
	"time"

	"mall-go/services/user-service/application/dto"
	"mall-go/services/user-service/domain/model"
)

// PermissionToDTO converts a Permission domain model to a PermissionDetailResponse DTO
func PermissionToDTO(permission *model.Permission) *dto.PermissionDetailResponse {
	if permission == nil {
		return nil
	}

	// Format time as string
	createdAt := ""
	updatedAt := ""

	if !permission.CreatedAt.IsZero() {
		createdAt = permission.CreatedAt.Format(time.RFC3339)
	}

	if !permission.UpdatedAt.IsZero() {
		updatedAt = permission.UpdatedAt.Format(time.RFC3339)
	}

	// Convert status to int
	var status int
	if permission.Status == model.PermissionStatusActive {
		status = 1
	} else {
		status = 0
	}

	return &dto.PermissionDetailResponse{
		ID:          permission.ID,
		Name:        permission.Name,
		Value:       permission.Value,
		Description: "", // The domain model doesn't have Description
		Type:        string(permission.Type),
		ParentID:    "", // The domain model doesn't have ParentID
		Sort:        0,  // The domain model doesn't have Sort
		Status:      status,
		CreatedAt:   createdAt,
		UpdatedAt:   updatedAt,
	}
}

// PermissionsToDTOs converts a slice of Permission domain models to a slice of PermissionDetailResponse DTOs
func PermissionsToDTOs(permissions []*model.Permission) []dto.PermissionDetailResponse {
	result := make([]dto.PermissionDetailResponse, 0, len(permissions))

	for _, permission := range permissions {
		dto := PermissionToDTO(permission)
		if dto != nil {
			result = append(result, *dto)
		}
	}

	return result
}
