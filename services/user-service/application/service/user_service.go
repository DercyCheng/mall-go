package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"

	"mall-go/pkg/auth"
	"mall-go/services/user-service/application/dto"
	"mall-go/services/user-service/domain/model"
	"mall-go/services/user-service/domain/repository"
)

// UserService 用户应用服务
type UserService struct {
	userRepo       repository.UserRepository
	roleRepo       repository.RoleRepository
	permissionRepo repository.PermissionRepository
}

// NewUserService 创建用户应用服务实例
func NewUserService(
	userRepo repository.UserRepository,
	roleRepo repository.RoleRepository,
	permissionRepo repository.PermissionRepository,
) *UserService {
	return &UserService{
		userRepo:       userRepo,
		roleRepo:       roleRepo,
		permissionRepo: permissionRepo,
	}
}

// Register 用户注册
func (s *UserService) Register(ctx context.Context, req dto.RegisterUserRequest) (*dto.UserResponse, error) {
	// 检查用户名是否已存在
	existingUser, _ := s.userRepo.FindByUsername(ctx, req.Username)
	if existingUser != nil {
		return nil, errors.New("username already exists")
	}

	// 检查邮箱是否已存在
	existingUser, _ = s.userRepo.FindByEmail(ctx, req.Email)
	if existingUser != nil {
		return nil, errors.New("email already exists")
	}

	// 检查手机号是否已存在
	existingUser, _ = s.userRepo.FindByPhone(ctx, req.Phone)
	if existingUser != nil {
		return nil, errors.New("phone already exists")
	}

	// 创建用户领域模型
	user := &model.User{
		ID:                uuid.New().String(),
		Username:          req.Username,
		Email:             req.Email,
		Phone:             req.Phone,
		NickName:          req.NickName,
		Gender:            req.Gender,
		City:              req.City,
		Job:               req.Job,
		Status:            model.UserStatusActive,
		IntegrationPoints: 0,
		Roles:             []model.Role{},
		CreatedAt:         time.Now(),
		UpdatedAt:         time.Now(),
	}

	// 设置密码(加密)
	if err := user.SetPassword(req.Password); err != nil {
		return nil, err
	}

	// 处理生日
	if req.Birthday != "" {
		birthDate, err := time.Parse("2006-01-02", req.Birthday)
		if err == nil {
			user.Birthday = &birthDate
		}
	}

	// 保存用户
	if err := s.userRepo.Save(ctx, user); err != nil {
		return nil, err
	}

	// 转换为响应DTO
	return s.toUserResponse(user), nil
}

// Login 用户登录
func (s *UserService) Login(ctx context.Context, req dto.LoginRequest) (*dto.LoginResponse, error) {
	// 根据用户名查找用户
	user, err := s.userRepo.FindByUsername(ctx, req.Username)
	if err != nil {
		return nil, errors.New("invalid username or password")
	}

	// 检查用户状态
	if user.Status != model.UserStatusActive {
		return nil, errors.New("account is inactive or locked")
	}

	// 验证密码
	if !user.VerifyPassword(req.Password) {
		return nil, errors.New("invalid username or password")
	}

	// 更新登录时间
	user.UpdateLoginTime()
	if err := s.userRepo.UpdateLastLoginTime(ctx, user.ID); err != nil {
		return nil, err
	}

	// 获取用户角色
	roles, err := s.userRepo.GetUserRoles(ctx, user.ID)
	if err == nil {
		user.Roles = roles
	}

	// 生成JWT令牌
	token, err := auth.GenerateToken(user.Username, uint(0)) // 使用UUID作为用户ID时需要修改JWT库
	if err != nil {
		return nil, err
	}

	// 创建登录响应
	userResp := s.toUserResponse(user)
	return &dto.LoginResponse{
		Token:     token,
		TokenType: "Bearer",
		ExpiresIn: 24 * 3600, // 24小时
		User:      *userResp,
	}, nil
}

// GetUserInfo 获取用户信息
func (s *UserService) GetUserInfo(ctx context.Context, id string) (*dto.UserResponse, error) {
	user, err := s.userRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// 获取用户角色
	roles, err := s.userRepo.GetUserRoles(ctx, user.ID)
	if err == nil {
		user.Roles = roles
	}

	return s.toUserResponse(user), nil
}

// UpdateUserInfo 更新用户信息
func (s *UserService) UpdateUserInfo(ctx context.Context, id string, req dto.UpdateUserRequest) (*dto.UserResponse, error) {
	// 获取原用户信息
	user, err := s.userRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// 检查邮箱是否已被其他用户使用
	if req.Email != "" && req.Email != user.Email {
		existingUser, _ := s.userRepo.FindByEmail(ctx, req.Email)
		if existingUser != nil && existingUser.ID != id {
			return nil, errors.New("email already in use by another user")
		}
		user.Email = req.Email
	}

	// 检查手机号是否已被其他用户使用
	if req.Phone != "" && req.Phone != user.Phone {
		existingUser, _ := s.userRepo.FindByPhone(ctx, req.Phone)
		if existingUser != nil && existingUser.ID != id {
			return nil, errors.New("phone already in use by another user")
		}
		user.Phone = req.Phone
	}

	// 更新用户信息
	if req.NickName != "" {
		user.NickName = req.NickName
	}
	if req.Icon != "" {
		user.Icon = req.Icon
	}
	if req.Gender >= 0 {
		user.Gender = req.Gender
	}
	if req.City != "" {
		user.City = req.City
	}
	if req.Job != "" {
		user.Job = req.Job
	}

	// 处理生日
	if req.Birthday != "" {
		birthDate, err := time.Parse("2006-01-02", req.Birthday)
		if err == nil {
			user.Birthday = &birthDate
		}
	}

	// 更新时间
	user.UpdatedAt = time.Now()

	// 保存更新
	if err := s.userRepo.Update(ctx, user); err != nil {
		return nil, err
	}

	return s.toUserResponse(user), nil
}

// UpdatePassword 修改密码
func (s *UserService) UpdatePassword(ctx context.Context, id string, req dto.UpdatePasswordRequest) error {
	// 获取用户信息
	user, err := s.userRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}

	// 验证旧密码
	if !user.VerifyPassword(req.OldPassword) {
		return errors.New("old password is incorrect")
	}

	// 设置新密码
	if err := user.SetPassword(req.NewPassword); err != nil {
		return err
	}

	// 更新时间
	user.UpdatedAt = time.Now()

	// 保存更新
	return s.userRepo.Update(ctx, user)
}

// UpdateStatus 更新用户状态
func (s *UserService) UpdateStatus(ctx context.Context, id string, req dto.UpdateStatusRequest) error {
	status := model.UserStatus(req.Status)
	return s.userRepo.UpdateStatus(ctx, id, status)
}

// AssignRoles 分配角色
func (s *UserService) AssignRoles(ctx context.Context, userID string, req dto.AssignRoleRequest) error {
	// 获取用户信息
	_, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return err
	}

	// 开始分配角色
	// 先获取所有当前用户的角色
	currentRoles, err := s.userRepo.GetUserRoles(ctx, userID)
	if err != nil {
		return err
	}

	// 构建当前角色ID集合
	currentRoleIDs := make(map[string]bool)
	for _, role := range currentRoles {
		currentRoleIDs[role.ID] = true
	}

	// 构建新角色ID集合
	newRoleIDs := make(map[string]bool)
	for _, roleID := range req.RoleIDs {
		newRoleIDs[roleID] = true
	}

	// 需要添加的角色
	for _, roleID := range req.RoleIDs {
		if !currentRoleIDs[roleID] {
			// 验证角色是否存在
			role, err := s.roleRepo.FindByID(ctx, roleID)
			if err != nil {
				return fmt.Errorf("invalid role ID: %s", roleID)
			}

			// 添加角色
			if err := s.userRepo.AddUserRole(ctx, userID, role.ID); err != nil {
				return err
			}
		}
	}

	// 需要删除的角色
	for _, role := range currentRoles {
		if !newRoleIDs[role.ID] {
			if err := s.userRepo.RemoveUserRole(ctx, userID, role.ID); err != nil {
				return err
			}
		}
	}

	return nil
}

// GetUsers 获取用户列表
func (s *UserService) GetUsers(ctx context.Context, page, pageSize int) (*dto.UserListResponse, error) {
	users, total, err := s.userRepo.FindAll(ctx, page, pageSize)
	if err != nil {
		return nil, err
	}

	userResponses := make([]dto.UserResponse, 0, len(users))
	for _, user := range users {
		// 获取用户角色
		roles, err := s.userRepo.GetUserRoles(ctx, user.ID)
		if err == nil {
			user.Roles = roles
		}
		userResponses = append(userResponses, *s.toUserResponse(user))
	}

	return &dto.UserListResponse{
		List:     userResponses,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	}, nil
}

// SearchUsers 搜索用户
func (s *UserService) SearchUsers(ctx context.Context, req dto.SearchUserRequest) (*dto.UserListResponse, error) {
	users, total, err := s.userRepo.Search(ctx, req.Keyword, req.Page, req.PageSize)
	if err != nil {
		return nil, err
	}

	userResponses := make([]dto.UserResponse, 0, len(users))
	for _, user := range users {
		// 获取用户角色
		roles, err := s.userRepo.GetUserRoles(ctx, user.ID)
		if err == nil {
			user.Roles = roles
		}
		userResponses = append(userResponses, *s.toUserResponse(user))
	}

	return &dto.UserListResponse{
		List:     userResponses,
		Total:    total,
		Page:     req.Page,
		PageSize: req.PageSize,
	}, nil
}

// DeleteUser 删除用户
func (s *UserService) DeleteUser(ctx context.Context, id string) error {
	return s.userRepo.Delete(ctx, id)
}

// 辅助方法: 将领域模型转换为响应DTO
func (s *UserService) toUserResponse(user *model.User) *dto.UserResponse {
	resp := &dto.UserResponse{
		ID:                user.ID,
		Username:          user.Username,
		Email:             user.Email,
		Phone:             user.Phone,
		NickName:          user.NickName,
		Icon:              user.Icon,
		Gender:            user.Gender,
		City:              user.City,
		Job:               user.Job,
		Status:            string(user.Status),
		IntegrationPoints: user.IntegrationPoints,
		CreatedAt:         user.CreatedAt,
	}

	// 处理可选的生日
	if user.Birthday != nil {
		birthday := user.Birthday.Format("2006-01-02")
		resp.Birthday = &birthday
	}

	// 处理最后登录时间
	if user.LastLoginAt != nil {
		lastLogin := user.LastLoginAt.Format("2006-01-02 15:04:05")
		resp.LastLoginAt = &lastLogin
	}

	// 处理角色
	if len(user.Roles) > 0 {
		roles := make([]dto.RoleDTO, 0, len(user.Roles))
		for _, role := range user.Roles {
			roles = append(roles, dto.RoleDTO{
				ID:          role.ID,
				Name:        role.Name,
				Description: role.Description,
				Status:      string(role.Status),
			})
		}
		resp.Roles = roles
	}

	return resp
}
