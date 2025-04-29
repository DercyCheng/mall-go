package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"

	"mall-go/services/user-service/application/assembler"
	"mall-go/services/user-service/application/dto"
	"mall-go/services/user-service/domain/repository"
)

// RoleServiceImpl implements the RoleService interface
type RoleServiceImpl struct {
	roleRepo repository.RoleRepository
}

// NewRoleService creates a new RoleService implementation
func NewRoleService(roleRepo repository.RoleRepository) RoleService {
	return &RoleServiceImpl{
		roleRepo: roleRepo,
	}
}

// CreateRole creates a new role
func (s *RoleServiceImpl) CreateRole(ctx context.Context, req dto.RoleCreateRequest) (string, error) {
	// Check if role name already exists
	existingRole, err := s.roleRepo.FindByName(ctx, req.Name)
	if err != nil {
		return "", fmt.Errorf("error checking existing role: %w", err)
	}

	if existingRole != nil {
		return "", errors.New("role name already exists")
	}

	// Create role
	role := assembler.ToRoleModel(req)
	role.ID = uuid.New().String()

	// Save role
	if err := s.roleRepo.Save(ctx, &role); err != nil {
		return "", fmt.Errorf("error saving role: %w", err)
	}

	return role.ID, nil
}

// GetRole retrieves role information by ID
func (s *RoleServiceImpl) GetRole(ctx context.Context, id string) (*dto.RoleDTO, error) {
	role, err := s.roleRepo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("error finding role: %w", err)
	}

	if role == nil {
		return nil, errors.New("role not found")
	}

	roleDTO := assembler.ToRoleDTO(*role)
	return &roleDTO, nil
}

// UpdateRole updates role information
func (s *RoleServiceImpl) UpdateRole(ctx context.Context, id string, req dto.RoleUpdateRequest) error {
	// Check if role exists
	role, err := s.roleRepo.FindByID(ctx, id)
	if err != nil {
		return fmt.Errorf("error finding role: %w", err)
	}

	if role == nil {
		return errors.New("role not found")
	}

	// Check if new name conflicts with an existing role
	if req.Name != "" && req.Name != role.Name {
		existingRole, err := s.roleRepo.FindByName(ctx, req.Name)
		if err != nil {
			return fmt.Errorf("error checking existing role: %w", err)
		}

		if existingRole != nil && existingRole.ID != id {
			return errors.New("role name already exists")
		}

		role.Name = req.Name
	}

	// Update description if provided
	if req.Description != "" {
		role.Description = req.Description
	}

	// Update timestamp
	role.UpdatedAt = time.Now()

	// Save role
	return s.roleRepo.Update(ctx, role)
}

// DeleteRole deletes a role
func (s *RoleServiceImpl) DeleteRole(ctx context.Context, id string) error {
	// Check if role exists
	role, err := s.roleRepo.FindByID(ctx, id)
	if err != nil {
		return fmt.Errorf("error finding role: %w", err)
	}

	if role == nil {
		return errors.New("role not found")
	}

	// Delete role
	return s.roleRepo.Delete(ctx, id)
}

// ListRoles retrieves a list of all roles
func (s *RoleServiceImpl) ListRoles(ctx context.Context) ([]dto.RoleDTO, error) {
	roles, err := s.roleRepo.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("error listing roles: %w", err)
	}

	roleDTOs := assembler.ToRoleDTOList(roles)
	return roleDTOs, nil
}

// AssignRolesToUser assigns roles to a user
func (s *RoleServiceImpl) AssignRolesToUser(ctx context.Context, req dto.AssignRoleRequest) error {
	// Get current user roles to determine what needs to be added or removed
	currentRoles, err := s.roleRepo.GetUserRoles(ctx, req.UserID)
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
			if err := s.roleRepo.AssignRoleToUser(ctx, req.UserID, roleID); err != nil {
				return fmt.Errorf("error assigning role: %w", err)
			}
		}
	}

	// Remove roles that are no longer requested
	for _, role := range currentRoles {
		if !requestedRoleMap[role.ID] {
			if err := s.roleRepo.RevokeRoleFromUser(ctx, req.UserID, role.ID); err != nil {
				return fmt.Errorf("error revoking role: %w", err)
			}
		}
	}

	return nil
}

// GetUserRoles retrieves roles for a user
func (s *RoleServiceImpl) GetUserRoles(ctx context.Context, userID string) ([]dto.RoleDTO, error) {
	roles, err := s.roleRepo.GetUserRoles(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("error getting user roles: %w", err)
	}

	roleDTOs := assembler.ToRoleDTOList(roles)
	return roleDTOs, nil
}
