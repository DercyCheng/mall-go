// filepath: /Users/dercyc/go/src/pro/mall-go/services/user-service/application/service/user_service_gomock_test.go
package service_test

import (
	"context"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"

	"mall-go/services/user-service/application/dto"
	"mall-go/services/user-service/application/service"
	"mall-go/services/user-service/domain/model"
	"mall-go/services/user-service/mocks"
)

// 使用gomock测试UserService的Register方法
func TestUserService_Register_WithGomock(t *testing.T) {
	// 创建gomock控制器
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// 创建mock的用户仓库和角色仓库
	mockUserRepo := mocks.NewMockUserRepository(ctrl)
	mockRoleRepo := mocks.NewMockRoleRepository(ctrl)

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

		// 设置mock期望行为
		// 1. 检查用户名是否已存在
		mockUserRepo.EXPECT().
			FindByUsername(gomock.Any(), "testuser").
			Return(nil, nil)

		// 2. 检查邮箱是否已存在
		mockUserRepo.EXPECT().
			FindByEmail(gomock.Any(), "test@example.com").
			Return(nil, nil)

		// 3. 保存用户
		mockUserRepo.EXPECT().
			Save(gomock.Any(), gomock.Any()).
			DoAndReturn(func(_ context.Context, user *model.User) error {
				// 验证传递给Save的用户对象
				assert.Equal(t, "testuser", user.Username)
				assert.Equal(t, "test@example.com", user.Email)
				assert.Equal(t, "Test User", user.NickName)
				assert.Equal(t, model.UserStatus(1), user.Status)
				assert.NotEmpty(t, user.Password) // 密码应已加密

				// 模拟设置ID
				user.ID = "user-123"
				return nil
			})

		// 执行测试
		id, err := userService.Register(context.Background(), req)

		// 验证结果
		assert.NoError(t, err)
		assert.Equal(t, "user-123", id)
	})

	// 测试用例2：用户名已存在
	t.Run("Username_Already_Exists", func(t *testing.T) {
		// 准备测试请求
		req := dto.UserCreateRequest{
			Username: "existinguser",
			Password: "Password123!",
			Email:    "new@example.com",
			NickName: "Existing User",
			Status:   1,
		}

		// 设置mock期望行为
		// 1. 模拟用户名已存在的情况
		existingUser := &model.User{
			ID:       "user-456",
			Username: "existinguser",
		}
		mockUserRepo.EXPECT().
			FindByUsername(gomock.Any(), "existinguser").
			Return(existingUser, nil)

		// 注意：由于用户名已存在，不应该调用FindByEmail和Save方法

		// 执行测试
		id, err := userService.Register(context.Background(), req)

		// 验证结果
		assert.Error(t, err)
		assert.Empty(t, id)
		assert.Equal(t, "username already exists", err.Error())
	})

	// 测试用例3：邮箱已存在
	t.Run("Email_Already_Exists", func(t *testing.T) {
		// 准备测试请求
		req := dto.UserCreateRequest{
			Username: "newuser",
			Password: "Password123!",
			Email:    "existing@example.com",
			NickName: "New User",
			Status:   1,
		}

		// 设置mock期望行为
		// 1. 检查用户名是否已存在（模拟不存在）
		mockUserRepo.EXPECT().
			FindByUsername(gomock.Any(), "newuser").
			Return(nil, nil)

		// 2. 模拟邮箱已存在的情况
		existingUser := &model.User{
			ID:    "user-789",
			Email: "existing@example.com",
		}
		mockUserRepo.EXPECT().
			FindByEmail(gomock.Any(), "existing@example.com").
			Return(existingUser, nil)

		// 注意：由于邮箱已存在，不应该调用Save方法

		// 执行测试
		id, err := userService.Register(context.Background(), req)

		// 验证结果
		assert.Error(t, err)
		assert.Empty(t, id)
		assert.Equal(t, "email already in use", err.Error())
	})
}

// 使用gomock测试UserService的Login方法
func TestUserService_Login_WithGomock(t *testing.T) {
	// 创建gomock控制器
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// 创建mock的用户仓库和角色仓库
	mockUserRepo := mocks.NewMockUserRepository(ctrl)
	mockRoleRepo := mocks.NewMockRoleRepository(ctrl)

	// 设置测试数据
	testSecret := "test-secret"
	testExpiry := 24 * time.Hour

	// 创建服务实例
	userService := service.NewUserService(mockUserRepo, mockRoleRepo, testSecret, testExpiry)

	// 生成测试用密码哈希
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("Password123!"), bcrypt.DefaultCost)

	// 测试用例1：成功登录
	t.Run("Success", func(t *testing.T) {
		// 准备测试请求
		req := dto.UserLoginRequest{
			Username: "testuser",
			Password: "Password123!",
		}

		// 创建一个预期的用户模型
		expectedUser := &model.User{
			ID:        "user-123",
			Username:  "testuser",
			Password:  string(hashedPassword),
			Email:     "test@example.com",
			NickName:  "Test User",
			Status:    model.UserStatusActive,
			CreatedAt: time.Now(),
		}

		// 设置mock期望行为
		// 1. 查找用户
		mockUserRepo.EXPECT().
			FindByUsername(gomock.Any(), "testuser").
			Return(expectedUser, nil)

		// 2. 更新登录时间
		mockUserRepo.EXPECT().
			Update(gomock.Any(), gomock.Any()).
			Return(nil)

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
		// 准备测试请求
		req := dto.UserLoginRequest{
			Username: "nonexistentuser",
			Password: "Password123!",
		}

		// 设置mock期望行为
		// 模拟用户不存在的情况
		mockUserRepo.EXPECT().
			FindByUsername(gomock.Any(), "nonexistentuser").
			Return(nil, nil)

		// 执行测试
		resp, err := userService.Login(context.Background(), req)

		// 验证结果
		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.Equal(t, "invalid username or password", err.Error())
	})

	// 测试用例3：密码错误
	t.Run("Incorrect_Password", func(t *testing.T) {
		// 准备测试请求
		req := dto.UserLoginRequest{
			Username: "testuser",
			Password: "WrongPassword123!",
		}

		// 创建一个预期的用户模型
		expectedUser := &model.User{
			ID:        "user-123",
			Username:  "testuser",
			Password:  string(hashedPassword), // 这是"Password123!"的哈希
			Email:     "test@example.com",
			NickName:  "Test User",
			Status:    model.UserStatusActive,
			CreatedAt: time.Now(),
		}

		// 设置mock期望行为
		// 1. 查找用户
		mockUserRepo.EXPECT().
			FindByUsername(gomock.Any(), "testuser").
			Return(expectedUser, nil)

		// 执行测试
		resp, err := userService.Login(context.Background(), req)

		// 验证结果
		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.Equal(t, "invalid username or password", err.Error())
	})

	// 测试用例4：账号已禁用
	t.Run("Account_Disabled", func(t *testing.T) {
		// 准备测试请求
		req := dto.UserLoginRequest{
			Username: "disableduser",
			Password: "Password123!",
		}

		// 创建一个预期的用户模型
		expectedUser := &model.User{
			ID:        "user-456",
			Username:  "disableduser",
			Password:  string(hashedPassword),
			Email:     "disabled@example.com",
			NickName:  "Disabled User",
			Status:    model.UserStatusInactive, // 账号禁用
			CreatedAt: time.Now(),
		}

		// 设置mock期望行为
		// 1. 查找用户
		mockUserRepo.EXPECT().
			FindByUsername(gomock.Any(), "disableduser").
			Return(expectedUser, nil)

		// 执行测试
		resp, err := userService.Login(context.Background(), req)

		// 验证结果
		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.Equal(t, "user account is inactive", err.Error())
	})
}