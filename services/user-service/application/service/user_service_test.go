package service_test

import (
	"context"
	"testing"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/stretchr/testify/assert"

	"mall-go/pkg/errors"
	"mall-go/services/user-service/application/dto"
	"mall-go/services/user-service/application/service"
	"mall-go/services/user-service/domain/model"
)

// hashPassword is a helper function for testing to create password hashes
func hashPassword(password string) (string, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedBytes), nil
}

// MockUserRepository is a mock implementation of repository.UserRepository
type MockUserRepository struct {
	users        map[string]*model.User
	usersByName  map[string]*model.User
	usersByEmail map[string]*model.User
	lastUser     *model.User
	lastID       string
}

// NewMockUserRepository creates a new mock user repository
func NewMockUserRepository() *MockUserRepository {
	return &MockUserRepository{
		users:        make(map[string]*model.User),
		usersByName:  make(map[string]*model.User),
		usersByEmail: make(map[string]*model.User),
	}
}

// Save implements repository.UserRepository.Save
func (m *MockUserRepository) Save(ctx context.Context, user *model.User) error {
	m.lastUser = user
	if user.ID == "" {
		user.ID = "user-123" // Mock generated ID
	}
	m.users[user.ID] = user
	m.usersByName[user.Username] = user
	m.usersByEmail[user.Email] = user
	return nil
}

// FindByID implements repository.UserRepository.FindByID
func (m *MockUserRepository) FindByID(ctx context.Context, id string) (*model.User, error) {
	m.lastID = id
	return m.users[id], nil
}

// FindByUsername implements repository.UserRepository.FindByUsername
func (m *MockUserRepository) FindByUsername(ctx context.Context, username string) (*model.User, error) {
	return m.usersByName[username], nil
}

// FindByEmail implements repository.UserRepository.FindByEmail
func (m *MockUserRepository) FindByEmail(ctx context.Context, email string) (*model.User, error) {
	return m.usersByEmail[email], nil
}

// Update implements repository.UserRepository.Update
func (m *MockUserRepository) Update(ctx context.Context, user *model.User) error {
	m.lastUser = user
	m.users[user.ID] = user
	m.usersByName[user.Username] = user
	m.usersByEmail[user.Email] = user
	return nil
}

// Delete implements repository.UserRepository.Delete
func (m *MockUserRepository) Delete(ctx context.Context, id string) error {
	m.lastID = id
	if user, ok := m.users[id]; ok {
		delete(m.usersByName, user.Username)
		delete(m.usersByEmail, user.Email)
		delete(m.users, id)
	}
	return nil
}

// List implements repository.UserRepository.List
func (m *MockUserRepository) List(ctx context.Context, page, pageSize int) ([]*model.User, int64, error) {
	users := make([]*model.User, 0, len(m.users))
	for _, user := range m.users {
		users = append(users, user)
	}
	return users, int64(len(users)), nil
}

// Search implements repository.UserRepository.Search
func (m *MockUserRepository) Search(ctx context.Context, query string, page, pageSize int) ([]*model.User, int64, error) {
	// Simple mock implementation that returns all users
	users := make([]*model.User, 0, len(m.users))
	for _, user := range m.users {
		users = append(users, user)
	}
	return users, int64(len(users)), nil
}

// MockRoleRepository is a mock implementation of repository.RoleRepository
type MockRoleRepository struct {
	roles       map[string]*model.Role
	rolesByName map[string]*model.Role
	userRoles   map[string][]*model.Role
	lastRole    *model.Role
	lastID      string
}

// NewMockRoleRepository creates a new mock role repository
func NewMockRoleRepository() *MockRoleRepository {
	return &MockRoleRepository{
		roles:       make(map[string]*model.Role),
		rolesByName: make(map[string]*model.Role),
		userRoles:   make(map[string][]*model.Role),
	}
}

// Save implements repository.RoleRepository.Save
func (m *MockRoleRepository) Save(ctx context.Context, role *model.Role) error {
	m.lastRole = role
	if role.ID == "" {
		role.ID = "role-123" // Mock generated ID
	}
	m.roles[role.ID] = role
	m.rolesByName[role.Name] = role
	return nil
}

// FindByID implements repository.RoleRepository.FindByID
func (m *MockRoleRepository) FindByID(ctx context.Context, id string) (*model.Role, error) {
	m.lastID = id
	return m.roles[id], nil
}

// FindByName implements repository.RoleRepository.FindByName
func (m *MockRoleRepository) FindByName(ctx context.Context, name string) (*model.Role, error) {
	return m.rolesByName[name], nil
}

// Update implements repository.RoleRepository.Update
func (m *MockRoleRepository) Update(ctx context.Context, role *model.Role) error {
	m.lastRole = role
	m.roles[role.ID] = role
	m.rolesByName[role.Name] = role
	return nil
}

// Delete implements repository.RoleRepository.Delete
func (m *MockRoleRepository) Delete(ctx context.Context, id string) error {
	m.lastID = id
	if role, ok := m.roles[id]; ok {
		delete(m.rolesByName, role.Name)
		delete(m.roles, id)
	}
	return nil
}

// List implements repository.RoleRepository.List
func (m *MockRoleRepository) List(ctx context.Context) ([]*model.Role, error) {
	roles := make([]*model.Role, 0, len(m.roles))
	for _, role := range m.roles {
		roles = append(roles, role)
	}
	return roles, nil
}

// GetUserRoles implements repository.RoleRepository.GetUserRoles
func (m *MockRoleRepository) GetUserRoles(ctx context.Context, userID string) ([]*model.Role, error) {
	return m.userRoles[userID], nil
}

// AssignRoleToUser implements repository.RoleRepository.AssignRoleToUser
func (m *MockRoleRepository) AssignRoleToUser(ctx context.Context, userID, roleID string) error {
	role, ok := m.roles[roleID]
	if !ok {
		return errors.New("role not found")
	}

	roles, ok := m.userRoles[userID]
	if !ok {
		roles = make([]*model.Role, 0)
	}

	// Check if role is already assigned
	for _, r := range roles {
		if r.ID == roleID {
			return nil
		}
	}

	roles = append(roles, role)
	m.userRoles[userID] = roles
	return nil
}

// RevokeRoleFromUser implements repository.RoleRepository.RevokeRoleFromUser
func (m *MockRoleRepository) RevokeRoleFromUser(ctx context.Context, userID, roleID string) error {
	roles, ok := m.userRoles[userID]
	if !ok {
		return nil
	}

	// Remove role from user's roles
	for i, role := range roles {
		if role.ID == roleID {
			roles = append(roles[:i], roles[i+1:]...)
			m.userRoles[userID] = roles
			return nil
		}
	}

	return nil
}

func TestUserService_Register(t *testing.T) {
	// 创建模拟的仓库
	mockUserRepo := NewMockUserRepository()
	mockRoleRepo := NewMockRoleRepository()

	// 设置测试数据
	testSecret := "test-secret"
	testExpiry := 24 * time.Hour

	// 创建服务实例
	userService := service.NewUserService(mockUserRepo, mockRoleRepo, testSecret, testExpiry)

	// 测试用例1：成功注册
	t.Run("Success", func(t *testing.T) {
		// 准备测试请求
		req := dto.UserCreateRequest{
			Username: "testuser",
			Password: "Password123!",
			Email:    "test@example.com",
			NickName: "Test User",
			Status:   1,
		}

		// 执行测试
		id, err := userService.Register(context.Background(), req)

		// 验证结果
		assert.NoError(t, err)
		assert.Equal(t, "user-123", id)

		// 验证用户是否被正确保存
		createdUser := mockUserRepo.lastUser
		assert.Equal(t, "testuser", createdUser.Username)
		assert.Equal(t, "test@example.com", createdUser.Email)
		assert.Equal(t, "Test User", createdUser.NickName)
		assert.Equal(t, model.UserStatus(1), createdUser.Status)
		assert.NotEmpty(t, createdUser.Password) // 密码应已加密
	})

	// 测试用例2：用户名已存在
	t.Run("Username_Already_Exists", func(t *testing.T) {
		// 重置模拟仓库
		mockUserRepo = NewMockUserRepository()
		userService = service.NewUserService(mockUserRepo, mockRoleRepo, testSecret, testExpiry)

		// 设置预先存在的用户
		existingUser := &model.User{
			ID:       "user-456",
			Username: "existinguser",
		}
		mockUserRepo.users["user-456"] = existingUser
		mockUserRepo.usersByName["existinguser"] = existingUser

		// 准备测试请求
		req := dto.UserCreateRequest{
			Username: "existinguser",
			Password: "Password123!",
			Email:    "new@example.com",
			NickName: "Existing User",
			Status:   1,
		}

		// 执行测试
		id, err := userService.Register(context.Background(), req)

		// 验证结果
		assert.Error(t, err)
		assert.Empty(t, id)
		appErr, ok := errors.As(err)
		assert.True(t, ok)
		assert.Equal(t, errors.CodeUserAlreadyExists, appErr.Code)
	})

	// 测试用例3：邮箱已存在
	t.Run("Email_Already_Exists", func(t *testing.T) {
		// 重置模拟仓库
		mockUserRepo = NewMockUserRepository()
		userService = service.NewUserService(mockUserRepo, mockRoleRepo, testSecret, testExpiry)

		// 设置预先存在的用户
		existingUser := &model.User{
			ID:    "user-789",
			Email: "existing@example.com",
		}
		mockUserRepo.users["user-789"] = existingUser
		mockUserRepo.usersByEmail["existing@example.com"] = existingUser

		// 准备测试请求
		req := dto.UserCreateRequest{
			Username: "newuser",
			Password: "Password123!",
			Email:    "existing@example.com",
			NickName: "New User",
			Status:   1,
		}

		// 执行测试
		id, err := userService.Register(context.Background(), req)

		// 验证结果
		assert.Error(t, err)
		assert.Empty(t, id)
		appErr, ok := errors.As(err)
		assert.True(t, ok)
		assert.Equal(t, errors.CodeUserAlreadyExists, appErr.Code)
	})
}

func TestUserService_Login(t *testing.T) {
	// 创建模拟的仓库
	mockUserRepo := NewMockUserRepository()
	mockRoleRepo := NewMockRoleRepository()

	// 设置测试数据
	testSecret := "test-secret"
	testExpiry := 24 * time.Hour

	// 创建服务实例
	userService := service.NewUserService(mockUserRepo, mockRoleRepo, testSecret, testExpiry)

	// 测试用例1：成功登录
	t.Run("Success", func(t *testing.T) {
		// 准备测试请求
		req := dto.UserLoginRequest{
			Username: "testuser",
			Password: "Password123!",
		}

		// 设置预先存在的用户
		hashedPassword, _ := hashPassword("Password123!")
		existingUser := &model.User{
			ID:        "user-123",
			Username:  "testuser",
			Password:  hashedPassword,
			Email:     "test@example.com",
			NickName:  "Test User",
			Status:    model.UserStatusActive,
			CreatedAt: time.Now(),
		}
		mockUserRepo.users["user-123"] = existingUser
		mockUserRepo.usersByName["testuser"] = existingUser

		// 执行测试
		resp, err := userService.Login(context.Background(), req)

		// 验证结果
		assert.NoError(t, err)
		assert.NotEmpty(t, resp.Token)
		assert.Equal(t, "Bearer", resp.TokenHead)
		assert.NotEmpty(t, resp.ExpireAt)
		assert.Equal(t, "testuser", resp.User.Username)
		assert.Equal(t, "test@example.com", resp.User.Email)
	})

	// 测试用例2：用户不存在
	t.Run("User_Not_Found", func(t *testing.T) {
		// 重置模拟仓库
		mockUserRepo = NewMockUserRepository()
		userService = service.NewUserService(mockUserRepo, mockRoleRepo, testSecret, testExpiry)

		// 准备测试请求
		req := dto.UserLoginRequest{
			Username: "nonexistentuser",
			Password: "Password123!",
		}

		// 执行测试
		resp, err := userService.Login(context.Background(), req)

		// 验证结果
		assert.Error(t, err)
		assert.Nil(t, resp)
		appErr, ok := errors.As(err)
		assert.True(t, ok)
		assert.Equal(t, errors.CodeInvalidCredentials, appErr.Code)
	})

	// 测试用例3：密码错误
	t.Run("Incorrect_Password", func(t *testing.T) {
		// 重置模拟仓库
		mockUserRepo = NewMockUserRepository()
		userService = service.NewUserService(mockUserRepo, mockRoleRepo, testSecret, testExpiry)

		// 设置预先存在的用户
		hashedPassword, _ := hashPassword("Password123!")
		existingUser := &model.User{
			ID:        "user-123",
			Username:  "testuser",
			Password:  hashedPassword,
			Email:     "test@example.com",
			NickName:  "Test User",
			Status:    model.UserStatusActive,
			CreatedAt: time.Now(),
		}
		mockUserRepo.users["user-123"] = existingUser
		mockUserRepo.usersByName["testuser"] = existingUser

		// 准备测试请求
		req := dto.UserLoginRequest{
			Username: "testuser",
			Password: "WrongPassword123!",
		}

		// 执行测试
		resp, err := userService.Login(context.Background(), req)

		// 验证结果
		assert.Error(t, err)
		assert.Nil(t, resp)
		appErr, ok := errors.As(err)
		assert.True(t, ok)
		assert.Equal(t, errors.CodeInvalidCredentials, appErr.Code)
	})

	// 测试用例4：账号已禁用
	t.Run("Account_Disabled", func(t *testing.T) {
		// 重置模拟仓库
		mockUserRepo = NewMockUserRepository()
		userService = service.NewUserService(mockUserRepo, mockRoleRepo, testSecret, testExpiry)

		// 设置预先存在的用户
		hashedPassword, _ := hashPassword("Password123!")
		existingUser := &model.User{
			ID:        "user-456",
			Username:  "disableduser",
			Password:  hashedPassword,
			Email:     "disabled@example.com",
			NickName:  "Disabled User",
			Status:    0, // 账号禁用
			CreatedAt: time.Now(),
		}
		mockUserRepo.users["user-456"] = existingUser
		mockUserRepo.usersByName["disableduser"] = existingUser

		// 准备测试请求
		req := dto.UserLoginRequest{
			Username: "disableduser",
			Password: "Password123!",
		}

		// 执行测试
		resp, err := userService.Login(context.Background(), req)

		// 验证结果
		assert.Error(t, err)
		assert.Nil(t, resp)
		appErr, ok := errors.As(err)
		assert.True(t, ok)
		assert.Equal(t, errors.CodeInvalidCredentials, appErr.Code)
	})
}
