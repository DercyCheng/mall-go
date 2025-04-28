package service

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"

	"mall-go/services/user-service/application/assembler"
	"mall-go/services/user-service/application/dto"
	"mall-go/services/user-service/domain/model"
	"mall-go/services/user-service/domain/repository"
)

var (
	ErrPermissionNotFound          = errors.New("permission not found")
	ErrPermissionValueAlreadyExists = errors.New("permission value already exists")
	ErrInvalidPermissionStatus     = errors.New("invalid permission status")
)

// PermissionService 权限服务
type PermissionService struct {
	permissionRepo repository.PermissionRepository
	roleRepo       repository.RoleRepository
}

// NewPermissionService 创建权限服务
func NewPermissionService(
	permissionRepo repository.PermissionRepository,
	roleRepo repository.RoleRepository,
) *PermissionService {
	return &PermissionService{
		permissionRepo: permissionRepo,
		roleRepo:       roleRepo,
	}
}

// GetPermissionRepository 获取权限仓储接口
func (s *PermissionService) GetPermissionRepository() repository.PermissionRepository {
	return s.permissionRepo
}

// GetRoleRepository 获取角色仓储接口
func (s *PermissionService) GetRoleRepository() repository.RoleRepository {
	return s.roleRepo
}

// CreatePermission 创建权限
func (s *PermissionService) CreatePermission(ctx context.Context, req *dto.CreatePermissionRequest) (*dto.PermissionResponse, error) {
	// 检查权限值是否已存在
	exists, err := s.permissionRepo.ExistsByValue(ctx, req.Value)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, ErrPermissionValueAlreadyExists
	}

	// 创建权限模型
	permission := &model.Permission{
		ID:          uuid.New().String(),
		Name:        req.Name,
		Value:       req.Value,
		Type:        req.Type,
		Description: req.Description,
		Status:      model.StatusEnabled,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// 保存权限
	if err := s.permissionRepo.Create(ctx, permission); err != nil {
		return nil, err
	}

	// 转换为DTO
	return assembler.PermissionToDTO(permission), nil
}

// GetPermission 获取权限
func (s *PermissionService) GetPermission(ctx context.Context, id string) (*dto.PermissionResponse, error) {
	// 获取权限
	permission, err := s.permissionRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if permission == nil {
		return nil, ErrPermissionNotFound
	}

	// 转换为DTO
	return assembler.PermissionToDTO(permission), nil
}

// UpdatePermission 更新权限
func (s *PermissionService) UpdatePermission(ctx context.Context, id string, req *dto.UpdatePermissionRequest) (*dto.PermissionResponse, error) {
	// 获取权限
	permission, err := s.permissionRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if permission == nil {
		return nil, ErrPermissionNotFound
	}

	// 如果权限值被修改，检查新值是否已存在
	if req.Value != permission.Value {
		exists, err := s.permissionRepo.ExistsByValue(ctx, req.Value)
		if err != nil {
			return nil, err
		}
		if exists {
			return nil, ErrPermissionValueAlreadyExists
		}
	}

	// 更新权限信息
	if req.Name != "" {
		permission.Name = req.Name
	}
	if req.Value != "" {
		permission.Value = req.Value
	}
	if req.Type != "" {
		permission.Type = req.Type
	}
	if req.Description != nil {
		permission.Description = *req.Description
	}
	permission.UpdatedAt = time.Now()

	// 保存权限
	if err := s.permissionRepo.Update(ctx, permission); err != nil {
		return nil, err
	}

	// 转换为DTO
	return assembler.PermissionToDTO(permission), nil
}

// DeletePermission 删除权限
func (s *PermissionService) DeletePermission(ctx context.Context, id string) error {
	// 检查权限是否存在
	permission, err := s.permissionRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if permission == nil {
		return ErrPermissionNotFound
	}

	// 删除权限
	return s.permissionRepo.Delete(ctx, id)
}

// ListPermissions 获取权限列表
func (s *PermissionService) ListPermissions(ctx context.Context, page, size int) (*dto.PermissionListResponse, error) {
	// 获取总数
	total, err := s.permissionRepo.Count(ctx)
	if err != nil {
		return nil, err
	}

	// 获取权限列表
	permissions, err := s.permissionRepo.FindAll(ctx, page, size)
	if err != nil {
		return nil, err
	}

	// 转换为DTO
	permissionDTOs := make([]*dto.PermissionResponse, 0, len(permissions))
	for _, permission := range permissions {
		permissionDTOs = append(permissionDTOs, assembler.PermissionToDTO(permission))
	}

	// 构建响应
	return &dto.PermissionListResponse{
		Total:       total,
		Permissions: permissionDTOs,
	}, nil
}

// ListPermissionsByType 根据类型获取权限列表
func (s *PermissionService) ListPermissionsByType(ctx context.Context, permType string, page, size int) (*dto.PermissionListResponse, error) {
	// 获取总数
	total, err := s.permissionRepo.CountByType(ctx, permType)
	if err != nil {
		return nil, err
	}

	// 获取权限列表
	permissions, err := s.permissionRepo.FindByType(ctx, permType, page, size)
	if err != nil {
		return nil, err
	}

	// 转换为DTO
	permissionDTOs := make([]*dto.PermissionResponse, 0, len(permissions))
	for _, permission := range permissions {
		permissionDTOs = append(permissionDTOs, assembler.PermissionToDTO(permission))
	}

	// 构建响应
	return &dto.PermissionListResponse{
		Total:       total,
		Permissions: permissionDTOs,
	}, nil
}

// UpdatePermissionStatus 更新权限状态
func (s *PermissionService) UpdatePermissionStatus(ctx context.Context, id string, status int) error {
	// 检查状态是否有效
	if status != model.StatusEnabled && status != model.StatusDisabled {
		return ErrInvalidPermissionStatus
	}

	// 检查权限是否存在
	permission, err := s.permissionRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if permission == nil {
		return ErrPermissionNotFound
	}

	// 更新状态
	permission.Status = status
	return s.permissionRepo.Update(ctx, permission)
}