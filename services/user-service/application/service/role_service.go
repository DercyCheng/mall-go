package service

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"

	"mall-go/services/user-service/application/dto"
	"mall-go/services/user-service/domain/model"
	"mall-go/services/user-service/domain/repository"
)

// RoleService 角色应用服务
type RoleService struct {
	roleRepo       repository.RoleRepository
	permissionRepo repository.PermissionRepository
}

// NewRoleService 创建角色应用服务实例
func NewRoleService(
	roleRepo repository.RoleRepository,
	permissionRepo repository.PermissionRepository,
) *RoleService {
	return &RoleService{
		roleRepo:       roleRepo,
		permissionRepo: permissionRepo,
	}
}

// CreateRole 创建角色
func (s *RoleService) CreateRole(ctx context.Context, name, description string) (*dto.RoleDTO, error) {
	// 检查角色名是否已存在
	existingRole, _ := s.roleRepo.FindByName(ctx, name)
	if existingRole != nil {
		return nil, errors.New("role name already exists")
	}

	// 创建角色领域模型
	role := &model.Role{
		ID:          uuid.New().String(),
		Name:        name,
		Description: description,
		Status:      model.RoleStatusActive,
		Sort:        0,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// 保存角色
	if err := s.roleRepo.Save(ctx, role); err != nil {
		return nil, err
	}

	// 转换为响应DTO
	return &dto.RoleDTO{
		ID:          role.ID,
		Name:        role.Name,
		Description: role.Description,
		Status:      string(role.Status),
	}, nil
}

// GetRole 获取角色详情
func (s *RoleService) GetRole(ctx context.Context, id string) (*dto.RoleDTO, error) {
	role, err := s.roleRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// 转换为响应DTO
	return &dto.RoleDTO{
		ID:          role.ID,
		Name:        role.Name,
		Description: role.Description,
		Status:      string(role.Status),
	}, nil
}

// UpdateRole 更新角色信息
func (s *RoleService) UpdateRole(ctx context.Context, id, name, description string) (*dto.RoleDTO, error) {
	// 获取原角色信息
	role, err := s.roleRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// 检查角色名是否已被其他角色使用
	if name != "" && name != role.Name {
		existingRole, _ := s.roleRepo.FindByName(ctx, name)
		if existingRole != nil && existingRole.ID != id {
			return nil, errors.New("role name already in use by another role")
		}
		role.Name = name
	}

	// 更新角色信息
	if description != "" {
		role.Description = description
	}

	// 更新时间
	role.UpdatedAt = time.Now()

	// 保存更新
	if err := s.roleRepo.Update(ctx, role); err != nil {
		return nil, err
	}

	// 转换为响应DTO
	return &dto.RoleDTO{
		ID:          role.ID,
		Name:        role.Name,
		Description: role.Description,
		Status:      string(role.Status),
	}, nil
}

// DeleteRole 删除角色
func (s *RoleService) DeleteRole(ctx context.Context, id string) error {
	return s.roleRepo.Delete(ctx, id)
}

// UpdateRoleStatus 更新角色状态
func (s *RoleService) UpdateRoleStatus(ctx context.Context, id string, status model.RoleStatus) error {
	return s.roleRepo.UpdateStatus(ctx, id, status)
}

// ListRoles 获取角色列表
func (s *RoleService) ListRoles(ctx context.Context, page, pageSize int) (*dto.RoleListResponse, error) {
	roles, total, err := s.roleRepo.FindAll(ctx, page, pageSize)
	if err != nil {
		return nil, err
	}

	roleDTOs := make([]dto.RoleDTO, 0, len(roles))
	for _, role := range roles {
		roleDTOs = append(roleDTOs, dto.RoleDTO{
			ID:          role.ID,
			Name:        role.Name,
			Description: role.Description,
			Status:      string(role.Status),
		})
	}

	return &dto.RoleListResponse{
		List:     roleDTOs,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	}, nil
}

// AssignPermissions 分配权限
func (s *RoleService) AssignPermissions(ctx context.Context, roleID string, permissionIDs []string) error {
	// 验证角色是否存在
	_, err := s.roleRepo.FindByID(ctx, roleID)
	if err != nil {
		return err
	}

	// 获取当前角色的所有权限
	currentPermissions, err := s.roleRepo.GetRolePermissions(ctx, roleID)
	if err != nil {
		return err
	}

	// 构建当前权限ID集合
	currentPermissionIDs := make(map[string]bool)
	for _, perm := range currentPermissions {
		currentPermissionIDs[perm.ID] = true
	}

	// 构建新权限ID集合
	newPermissionIDs := make(map[string]bool)
	for _, permID := range permissionIDs {
		newPermissionIDs[permID] = true
	}

	// 需要添加的权限
	for _, permID := range permissionIDs {
		if !currentPermissionIDs[permID] {
			// 验证权限是否存在
			_, err := s.permissionRepo.FindByID(ctx, permID)
			if err != nil {
				return errors.New("invalid permission ID: " + permID)
			}

			// 添加权限
			if err := s.roleRepo.AddRolePermission(ctx, roleID, permID); err != nil {
				return err
			}
		}
	}

	// 需要删除的权限
	for _, perm := range currentPermissions {
		if !newPermissionIDs[perm.ID] {
			if err := s.roleRepo.RemoveRolePermission(ctx, roleID, perm.ID); err != nil {
				return err
			}
		}
	}

	return nil
}

// GetRolePermissions 获取角色的权限
func (s *RoleService) GetRolePermissions(ctx context.Context, roleID string) ([]dto.PermissionDTO, error) {
	// 验证角色是否存在
	_, err := s.roleRepo.FindByID(ctx, roleID)
	if err != nil {
		return nil, err
	}

	// 获取角色的所有权限
	permissions, err := s.roleRepo.GetRolePermissions(ctx, roleID)
	if err != nil {
		return nil, err
	}

	// 转换为响应DTO
	permissionDTOs := make([]dto.PermissionDTO, 0, len(permissions))
	for _, perm := range permissions {
		permissionDTOs = append(permissionDTOs, dto.PermissionDTO{
			ID:     perm.ID,
			Name:   perm.Name,
			Value:  perm.Value,
			Type:   string(perm.Type),
			Status: string(perm.Status),
		})
	}

	return permissionDTOs, nil
}