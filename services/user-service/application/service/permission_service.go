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
	ErrPermissionNotFound           = errors.New("permission not found")
	ErrPermissionValueAlreadyExists = errors.New("permission value already exists")
	ErrInvalidPermissionStatus      = errors.New("invalid permission status")
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
func (s *PermissionService) CreatePermission(ctx context.Context, req *dto.CreatePermissionRequest) (*dto.PermissionDetailResponse, error) {
	// 检查权限值是否已存在
	existingPerm, err := s.permissionRepo.FindByValue(ctx, req.Value)
	if err != nil {
		return nil, err
	}
	if existingPerm != nil {
		return nil, ErrPermissionValueAlreadyExists
	}

	// 创建权限模型
	permission := &model.Permission{
		ID:        uuid.New().String(),
		Name:      req.Name,
		Value:     req.Value,
		Type:      model.PermissionType(req.Type),
		Status:    model.PermissionStatusActive,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// 保存权限
	if err := s.permissionRepo.Save(ctx, permission); err != nil {
		return nil, err
	}

	// 转换为DTO
	return assembler.PermissionToDTO(permission), nil
}

// GetPermission 获取权限
func (s *PermissionService) GetPermission(ctx context.Context, id string) (*dto.PermissionDetailResponse, error) {
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
func (s *PermissionService) UpdatePermission(ctx context.Context, id string, req *dto.UpdatePermissionRequest) (*dto.PermissionDetailResponse, error) {
	// 获取权限
	permission, err := s.permissionRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if permission == nil {
		return nil, ErrPermissionNotFound
	}

	// 如果权限值被修改，检查新值是否已存在
	if req.Value != "" && req.Value != permission.Value {
		existingPerm, err := s.permissionRepo.FindByValue(ctx, req.Value)
		if err != nil {
			return nil, err
		}
		if existingPerm != nil && existingPerm.ID != id {
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
		permission.Type = model.PermissionType(req.Type)
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
func (s *PermissionService) ListPermissions(ctx context.Context, page, size int) (*dto.PermissionPageResponse, error) {
	// 获取权限列表和总数
	permissions, total, err := s.permissionRepo.FindAll(ctx, page, size)
	if err != nil {
		return nil, err
	}

	// 转换为DTO
	permissionDTOs := assembler.PermissionsToDTOs(permissions)

	// 构建响应
	return &dto.PermissionPageResponse{
		List:     permissionDTOs,
		Total:    total,
		Page:     page,
		PageSize: size,
	}, nil
}

// ListPermissionsByType 根据类型获取权限列表
func (s *PermissionService) ListPermissionsByType(ctx context.Context, permType string, page, size int) (*dto.PermissionPageResponse, error) {
	// 转换类型字符串为领域模型的类型
	domainPermType := model.PermissionType(permType)

	// 获取权限列表
	permissions, err := s.permissionRepo.FindByType(ctx, domainPermType)
	if err != nil {
		return nil, err
	}

	// 分页处理 (因为FindByType方法不支持分页，我们在内存中处理)
	total := int64(len(permissions))
	start := (page - 1) * size
	end := start + size

	if start >= len(permissions) {
		permissions = []*model.Permission{}
	} else if end > len(permissions) {
		permissions = permissions[start:]
	} else {
		permissions = permissions[start:end]
	}

	// 转换为DTO
	permissionDTOs := assembler.PermissionsToDTOs(permissions)

	// 构建响应
	return &dto.PermissionPageResponse{
		List:     permissionDTOs,
		Total:    total,
		Page:     page,
		PageSize: size,
	}, nil
}

// UpdatePermissionStatus 更新权限状态
func (s *PermissionService) UpdatePermissionStatus(ctx context.Context, id string, statusValue string) error {
	// 将状态字符串转换为领域模型状态
	var status model.PermissionStatus
	if statusValue == "1" || statusValue == "enable" {
		status = model.PermissionStatusActive
	} else if statusValue == "0" || statusValue == "disable" {
		status = model.PermissionStatusInactive
	} else {
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
	return s.permissionRepo.UpdateStatus(ctx, id, status)
}
