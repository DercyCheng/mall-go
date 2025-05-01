package service_test

import (
	"context"
	"errors"
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

func TestUserServiceImpl_Register(t *testing.T) {
	// 创建 gomock 控制器
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// 创建模拟的依赖对象
	mockUserRepo := mocks.NewMockUserRepository(ctrl)
	mockRoleRepo := mocks.NewMockRoleRepository(ctrl)

	// 创建被测试的服务
	userService := service.NewUserService(
		mockUserRepo,
		mockRoleRepo,
		"test-secret",
		time.Hour*24,
	)

	// 测试场景：注册成功
	t.Run("Register_Success", func(t *testing.T) {
		// 设置输入
		req := dto.UserCreateRequest{
			Username: "testuser",
			Password: "password123",
			Email:    "test@example.com",
			NickName: "Test User",
			Phone:    "1234567890",
		}

		// 设置模拟行为
		mockUserRepo.EXPECT().FindByUsername(gomock.Any(), req.Username).Return(nil, nil)
		mockUserRepo.EXPECT().FindByEmail(gomock.Any(), req.Email).Return(nil, nil)
		mockUserRepo.EXPECT().Save(gomock.Any(), gomock.Any()).DoAndReturn(
			func(ctx context.Context, user *model.User) error {
				// 验证模型转换是否正确
				assert.Equal(t, req.Username, user.Username)
				assert.Equal(t, req.Email, user.Email)
				assert.Equal(t, req.NickName, user.NickName)
				assert.Equal(t, req.Phone, user.Phone)
				assert.NotEmpty(t, user.ID)
				assert.NotEmpty(t, user.Password)
				// 模拟数据库保存成功后设置ID
				if user.ID == "" {
					user.ID = "generated-user-id"
				}
				return nil
			})

		// 执行方法
		userID, err := userService.Register(context.Background(), req)

		// 断言结果
		assert.NoError(t, err)
		assert.NotEmpty(t, userID)
	})

	// 测试场景：用户名已存在
	t.Run("Register_UsernameAlreadyExists", func(t *testing.T) {
		// 设置输入
		req := dto.UserCreateRequest{
			Username: "existinguser",
			Password: "password123",
			Email:    "test@example.com",
		}

		// 设置模拟行为 - 模拟用户名已存在
		existingUser := &model.User{
			ID:       "existing-id",
			Username: "existinguser",
		}
		mockUserRepo.EXPECT().FindByUsername(gomock.Any(), req.Username).Return(existingUser, nil)

		// 执行方法
		userID, err := userService.Register(context.Background(), req)

		// 断言结果
		assert.Error(t, err)
		assert.Equal(t, "username already exists", err.Error())
		assert.Empty(t, userID)
	})

	// 测试场景：邮箱已存在
	t.Run("Register_EmailAlreadyExists", func(t *testing.T) {
		// 设置输入
		req := dto.UserCreateRequest{
			Username: "newuser",
			Password: "password123",
			Email:    "existing@example.com",
		}

		// 设置模拟行为 - 模拟邮箱已存在
		mockUserRepo.EXPECT().FindByUsername(gomock.Any(), req.Username).Return(nil, nil)
		existingUser := &model.User{
			ID:    "existing-id",
			Email: "existing@example.com",
		}
		mockUserRepo.EXPECT().FindByEmail(gomock.Any(), req.Email).Return(existingUser, nil)

		// 执行方法
		userID, err := userService.Register(context.Background(), req)

		// 断言结果
		assert.Error(t, err)
		assert.Equal(t, "email already in use", err.Error())
		assert.Empty(t, userID)
	})

	// 测试场景：保存用户时发生错误
	t.Run("Register_SaveError", func(t *testing.T) {
		// 设置输入
		req := dto.UserCreateRequest{
			Username: "newuser",
			Password: "password123",
			Email:    "new@example.com",
		}

		// 设置模拟行为
		mockUserRepo.EXPECT().FindByUsername(gomock.Any(), req.Username).Return(nil, nil)
		mockUserRepo.EXPECT().FindByEmail(gomock.Any(), req.Email).Return(nil, nil)
		mockUserRepo.EXPECT().Save(gomock.Any(), gomock.Any()).Return(errors.New("database error"))

		// 执行方法
		userID, err := userService.Register(context.Background(), req)

		// 断言结果
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "error saving user")
		assert.Empty(t, userID)
	})
}

func TestUserServiceImpl_Login(t *testing.T) {
	// 创建 gomock 控制器
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// 创建模拟的依赖对象
	mockUserRepo := mocks.NewMockUserRepository(ctrl)
	mockRoleRepo := mocks.NewMockRoleRepository(ctrl)

	// 创建被测试的服务
	userService := service.NewUserService(
		mockUserRepo,
		mockRoleRepo,
		"test-secret",
		time.Hour*24,
	)

	// 测试场景：登录成功
	t.Run("Login_Success", func(t *testing.T) {
		// 设置输入
		req := dto.UserLoginRequest{
			Username: "testuser",
			Password: "password123",
		}

		// 创建一个有效的用户模型
		user := &model.User{
			ID:        "user-123",
			Username:  "testuser",
			Email:     "test@example.com",
			NickName:  "Test User",
			Status:    model.UserStatusActive,
			CreatedAt: time.Now().Add(-24 * time.Hour),
			LastLogin: time.Now().Add(-12 * time.Hour),
		}
		// 设置密码
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
		user.Password = string(hashedPassword)

		// 记录更新前的时间用于比较
		beforeLoginTime := user.LastLogin

		// 设置模拟行为
		mockUserRepo.EXPECT().FindByUsername(gomock.Any(), req.Username).Return(user, nil)
		mockUserRepo.EXPECT().Update(gomock.Any(), gomock.Any()).DoAndReturn(
			func(ctx context.Context, updatedUser *model.User) error {
				// 验证最后登录时间是否已更新 - 使用更可靠的方式
				assert.NotEqual(t, beforeLoginTime, updatedUser.LastLogin)
				assert.False(t, updatedUser.LastLogin.IsZero()) // 确保不是零值
				assert.True(t, time.Now().After(beforeLoginTime)) // 确保当前时间在旧登录时间之后
				return nil
			})

		// 执行方法
		response, err := userService.Login(context.Background(), req)

		// 断言结果
		assert.NoError(t, err)
		assert.NotNil(t, response)
		assert.NotEmpty(t, response.Token)
		assert.Equal(t, "Bearer", response.TokenHead)
		assert.NotEmpty(t, response.ExpireAt)
		assert.Equal(t, user.ID, response.User.ID)
		assert.Equal(t, user.Username, response.User.Username)
		assert.Equal(t, user.Email, response.User.Email)
	})

	// 测试场景：用户不存在
	t.Run("Login_UserNotFound", func(t *testing.T) {
		// 设置输入
		req := dto.UserLoginRequest{
			Username: "nonexistent",
			Password: "password123",
		}

		// 设置模拟行为
		mockUserRepo.EXPECT().FindByUsername(gomock.Any(), req.Username).Return(nil, nil)

		// 执行方法
		response, err := userService.Login(context.Background(), req)

		// 断言结果
		assert.Error(t, err)
		assert.Equal(t, "invalid username or password", err.Error())
		assert.Nil(t, response)
	})

	// 测试场景：密码错误
	t.Run("Login_WrongPassword", func(t *testing.T) {
		// 设置输入
		req := dto.UserLoginRequest{
			Username: "testuser",
			Password: "wrongpassword",
		}

		// 创建一个有效的用户模型
		user := &model.User{
			ID:       "user-123",
			Username: "testuser",
			Email:    "test@example.com",
			Status:   model.UserStatusActive,
		}
		// 设置正确的密码
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
		user.Password = string(hashedPassword)

		// 设置模拟行为
		mockUserRepo.EXPECT().FindByUsername(gomock.Any(), req.Username).Return(user, nil)

		// 执行方法
		response, err := userService.Login(context.Background(), req)

		// 断言结果
		assert.Error(t, err)
		assert.Equal(t, "invalid username or password", err.Error())
		assert.Nil(t, response)
	})

	// 测试场景：用户账号未激活
	t.Run("Login_InactiveUser", func(t *testing.T) {
		// 设置输入
		req := dto.UserLoginRequest{
			Username: "inactiveuser",
			Password: "password123",
		}

		// 创建一个未激活的用户模型
		user := &model.User{
			ID:       "user-123",
			Username: "inactiveuser",
			Email:    "inactive@example.com",
			Status:   model.UserStatusInactive,
		}
		// 设置密码
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
		user.Password = string(hashedPassword)

		// 设置模拟行为
		mockUserRepo.EXPECT().FindByUsername(gomock.Any(), req.Username).Return(user, nil)

		// 执行方法
		response, err := userService.Login(context.Background(), req)

		// 断言结果
		assert.Error(t, err)
		assert.Equal(t, "user account is inactive", err.Error())
		assert.Nil(t, response)
	})

	// 测试场景：更新用户登录时间失败
	t.Run("Login_UpdateError", func(t *testing.T) {
		// 设置输入
		req := dto.UserLoginRequest{
			Username: "testuser",
			Password: "password123",
		}

		// 创建一个有效的用户模型
		user := &model.User{
			ID:       "user-123",
			Username: "testuser",
			Email:    "test@example.com",
			Status:   model.UserStatusActive,
		}
		// 设置密码
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
		user.Password = string(hashedPassword)

		// 设置模拟行为
		mockUserRepo.EXPECT().FindByUsername(gomock.Any(), req.Username).Return(user, nil)
		mockUserRepo.EXPECT().Update(gomock.Any(), gomock.Any()).Return(errors.New("database error"))

		// 执行方法
		response, err := userService.Login(context.Background(), req)

		// 断言结果
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "error updating user login time")
		assert.Nil(t, response)
	})
}

func TestUserServiceImpl_GetUserInfo(t *testing.T) {
	// 创建 gomock 控制器
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// 创建模拟的依赖对象
	mockUserRepo := mocks.NewMockUserRepository(ctrl)
	mockRoleRepo := mocks.NewMockRoleRepository(ctrl)

	// 创建被测试的服务
	userService := service.NewUserService(
		mockUserRepo,
		mockRoleRepo,
		"test-secret",
		time.Hour*24,
	)

	// 测试场景：成功获取用户信息
	t.Run("GetUserInfo_Success", func(t *testing.T) {
		// 准备测试数据
		userID := "user-123"
		user := &model.User{
			ID:        userID,
			Username:  "testuser",
			Email:     "test@example.com",
			NickName:  "Test User",
			Phone:     "1234567890",
			Status:    model.UserStatusActive,
			CreatedAt: time.Now().Add(-24 * time.Hour),
		}

		// 设置模拟行为
		mockUserRepo.EXPECT().FindByID(gomock.Any(), userID).Return(user, nil)

		// 执行方法
		userDTO, err := userService.GetUserInfo(context.Background(), userID)

		// 断言结果
		assert.NoError(t, err)
		assert.NotNil(t, userDTO)
		assert.Equal(t, userID, userDTO.ID)
		assert.Equal(t, user.Username, userDTO.Username)
		assert.Equal(t, user.Email, userDTO.Email)
		assert.Equal(t, user.NickName, userDTO.NickName)
		assert.Equal(t, user.Phone, userDTO.Phone)
	})

	// 测试场景：用户不存在
	t.Run("GetUserInfo_UserNotFound", func(t *testing.T) {
		// 准备测试数据
		userID := "nonexistent-user"

		// 设置模拟行为
		mockUserRepo.EXPECT().FindByID(gomock.Any(), userID).Return(nil, nil)

		// 执行方法
		userDTO, err := userService.GetUserInfo(context.Background(), userID)

		// 断言结果
		assert.Error(t, err)
		assert.Equal(t, "user not found", err.Error())
		assert.Nil(t, userDTO)
	})

	// 测试场景：数据库错误
	t.Run("GetUserInfo_DatabaseError", func(t *testing.T) {
		// 准备测试数据
		userID := "user-123"

		// 设置模拟行为
		mockUserRepo.EXPECT().FindByID(gomock.Any(), userID).Return(nil, errors.New("database connection error"))

		// 执行方法
		userDTO, err := userService.GetUserInfo(context.Background(), userID)

		// 断言结果
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "error finding user")
		assert.Nil(t, userDTO)
	})
}

func TestUserServiceImpl_UpdateUser(t *testing.T) {
	// 创建 gomock 控制器
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// 创建模拟的依赖对象
	mockUserRepo := mocks.NewMockUserRepository(ctrl)
	mockRoleRepo := mocks.NewMockRoleRepository(ctrl)

	// 创建被测试的服务
	userService := service.NewUserService(
		mockUserRepo,
		mockRoleRepo,
		"test-secret",
		time.Hour*24,
	)

	// 测试场景：成功更新用户信息
	t.Run("UpdateUser_Success", func(t *testing.T) {
		// 准备测试数据
		userID := "user-123"
		user := &model.User{
			ID:        userID,
			Username:  "testuser",
			Email:     "old@example.com",
			NickName:  "Old Name",
			Phone:     "1234567890",
			Status:    model.UserStatusActive,
		}

		updateReq := dto.UserUpdateRequest{
			Email:    "new@example.com",
			NickName: "New Name",
			Phone:    "9876543210",
		}

		// 设置模拟行为
		mockUserRepo.EXPECT().FindByID(gomock.Any(), userID).Return(user, nil)
		mockUserRepo.EXPECT().Update(gomock.Any(), gomock.Any()).DoAndReturn(
			func(ctx context.Context, updatedUser *model.User) error {
				// 验证更新是否应用
				assert.Equal(t, updateReq.Email, updatedUser.Email)
				assert.Equal(t, updateReq.NickName, updatedUser.NickName)
				assert.Equal(t, updateReq.Phone, updatedUser.Phone)
				return nil
			})

		// 执行方法
		err := userService.UpdateUser(context.Background(), userID, updateReq)

		// 断言结果
		assert.NoError(t, err)
	})

	// 测试场景：用户不存在
	t.Run("UpdateUser_UserNotFound", func(t *testing.T) {
		// 准备测试数据
		userID := "nonexistent-user"
		updateReq := dto.UserUpdateRequest{
			Email: "new@example.com",
		}

		// 设置模拟行为
		mockUserRepo.EXPECT().FindByID(gomock.Any(), userID).Return(nil, nil)

		// 执行方法
		err := userService.UpdateUser(context.Background(), userID, updateReq)

		// 断言结果
		assert.Error(t, err)
		assert.Equal(t, "user not found", err.Error())
	})

	// 测试场景：更新数据库失败
	t.Run("UpdateUser_DatabaseError", func(t *testing.T) {
		// 准备测试数据
		userID := "user-123"
		user := &model.User{
			ID:        userID,
			Username:  "testuser",
			Email:     "old@example.com",
			NickName:  "Old Name",
		}

		updateReq := dto.UserUpdateRequest{
			Email: "new@example.com",
		}

		// 设置模拟行为
		mockUserRepo.EXPECT().FindByID(gomock.Any(), userID).Return(user, nil)
		mockUserRepo.EXPECT().Update(gomock.Any(), gomock.Any()).Return(errors.New("database error"))

		// 执行方法
		err := userService.UpdateUser(context.Background(), userID, updateReq)

		// 断言结果
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "error updating user")
	})

	// 测试场景：更新用户角色
	t.Run("UpdateUser_WithRoles", func(t *testing.T) {
		// 准备测试数据
		userID := "user-123"
		user := &model.User{
			ID:        userID,
			Username:  "testuser",
			Email:     "test@example.com",
		}

		// 使用正确的类型 []*model.Role 而不是 []model.Role
		currentRoles := []*model.Role{
			{ID: "role-1", Name: "admin"},
			{ID: "role-2", Name: "user"},
		}

		updateReq := dto.UserUpdateRequest{
			RoleIds: []string{"role-1", "role-3"}, // 保留 role-1，移除 role-2，添加 role-3
		}

		// 设置模拟行为
		mockUserRepo.EXPECT().FindByID(gomock.Any(), userID).Return(user, nil)
		mockUserRepo.EXPECT().Update(gomock.Any(), gomock.Any()).Return(nil)
		mockRoleRepo.EXPECT().GetUserRoles(gomock.Any(), userID).Return(currentRoles, nil)
		
		// 移除 role-2
		mockRoleRepo.EXPECT().RevokeRoleFromUser(gomock.Any(), userID, "role-2").Return(nil)
		
		// 添加 role-3
		mockRoleRepo.EXPECT().AssignRoleToUser(gomock.Any(), userID, "role-3").Return(nil)

		// 执行方法
		err := userService.UpdateUser(context.Background(), userID, updateReq)

		// 断言结果
		assert.NoError(t, err)
	})
}

func TestUserServiceImpl_DeleteUser(t *testing.T) {
	// 创建 gomock 控制器
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// 创建模拟的依赖对象
	mockUserRepo := mocks.NewMockUserRepository(ctrl)
	mockRoleRepo := mocks.NewMockRoleRepository(ctrl)

	// 创建被测试的服务
	userService := service.NewUserService(
		mockUserRepo,
		mockRoleRepo,
		"test-secret",
		time.Hour*24,
	)

	// 测试场景：成功删除用户
	t.Run("DeleteUser_Success", func(t *testing.T) {
		// 准备测试数据
		userID := "user-123"
		user := &model.User{
			ID:       userID,
			Username: "testuser",
		}

		// 设置模拟行为
		mockUserRepo.EXPECT().FindByID(gomock.Any(), userID).Return(user, nil)
		mockUserRepo.EXPECT().Delete(gomock.Any(), userID).Return(nil)

		// 执行方法
		err := userService.DeleteUser(context.Background(), userID)

		// 断言结果
		assert.NoError(t, err)
	})

	// 测试场景：用户不存在
	t.Run("DeleteUser_UserNotFound", func(t *testing.T) {
		// 准备测试数据
		userID := "nonexistent-user"

		// 设置模拟行为
		mockUserRepo.EXPECT().FindByID(gomock.Any(), userID).Return(nil, nil)

		// 执行方法
		err := userService.DeleteUser(context.Background(), userID)

		// 断言结果
		assert.Error(t, err)
		assert.Equal(t, "user not found", err.Error())
	})

	// 测试场景：删除用户时数据库错误
	t.Run("DeleteUser_DatabaseError", func(t *testing.T) {
		// 准备测试数据
		userID := "user-123"
		user := &model.User{
			ID:       userID,
			Username: "testuser",
		}

		// 设置模拟行为
		mockUserRepo.EXPECT().FindByID(gomock.Any(), userID).Return(user, nil)
		mockUserRepo.EXPECT().Delete(gomock.Any(), userID).Return(errors.New("database error"))

		// 执行方法
		err := userService.DeleteUser(context.Background(), userID)

		// 断言结果
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "error deleting user")
	})
}

func TestUserServiceImpl_ChangePassword(t *testing.T) {
	// 创建 gomock 控制器
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// 创建模拟的依赖对象
	mockUserRepo := mocks.NewMockUserRepository(ctrl)
	mockRoleRepo := mocks.NewMockRoleRepository(ctrl)

	// 创建被测试的服务
	userService := service.NewUserService(
		mockUserRepo,
		mockRoleRepo,
		"test-secret",
		time.Hour*24,
	)

	// 测试场景：成功修改密码
	t.Run("ChangePassword_Success", func(t *testing.T) {
		// 准备测试数据
		userID := "user-123"
		user := &model.User{
			ID:       userID,
			Username: "testuser",
		}
		
		// 设置旧密码
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("oldpassword"), bcrypt.DefaultCost)
		user.Password = string(hashedPassword)
		
		// 旧密码的哈希值，用于验证密码是否被修改
		oldPasswordHash := user.Password

		req := dto.UserChangePasswordRequest{
			OldPassword: "oldpassword",
			NewPassword: "newpassword",
		}

		// 设置模拟行为
		mockUserRepo.EXPECT().FindByID(gomock.Any(), userID).Return(user, nil)
		mockUserRepo.EXPECT().Update(gomock.Any(), gomock.Any()).DoAndReturn(
			func(ctx context.Context, updatedUser *model.User) error {
				// 验证密码是否被修改
				assert.NotEqual(t, oldPasswordHash, updatedUser.Password)
				return nil
			})

		// 执行方法
		err := userService.ChangePassword(context.Background(), userID, req)

		// 断言结果
		assert.NoError(t, err)
	})

	// 测试场景：用户不存在
	t.Run("ChangePassword_UserNotFound", func(t *testing.T) {
		// 准备测试数据
		userID := "nonexistent-user"
		req := dto.UserChangePasswordRequest{
			OldPassword: "oldpassword",
			NewPassword: "newpassword",
		}

		// 设置模拟行为
		mockUserRepo.EXPECT().FindByID(gomock.Any(), userID).Return(nil, nil)

		// 执行方法
		err := userService.ChangePassword(context.Background(), userID, req)

		// 断言结果
		assert.Error(t, err)
		assert.Equal(t, "user not found", err.Error())
	})

	// 测试场景：旧密码错误
	t.Run("ChangePassword_WrongOldPassword", func(t *testing.T) {
		// 准备测试数据
		userID := "user-123"
		user := &model.User{
			ID:       userID,
			Username: "testuser",
		}
		
		// 设置旧密码
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("correctoldpassword"), bcrypt.DefaultCost)
		user.Password = string(hashedPassword)

		req := dto.UserChangePasswordRequest{
			OldPassword: "wrongoldpassword",
			NewPassword: "newpassword",
		}

		// 设置模拟行为
		mockUserRepo.EXPECT().FindByID(gomock.Any(), userID).Return(user, nil)

		// 执行方法
		err := userService.ChangePassword(context.Background(), userID, req)

		// 断言结果
		assert.Error(t, err)
		assert.Equal(t, "current password is incorrect", err.Error())
	})

	// 测试场景：更新密码时数据库错误
	t.Run("ChangePassword_DatabaseError", func(t *testing.T) {
		// 准备测试数据
		userID := "user-123"
		user := &model.User{
			ID:       userID,
			Username: "testuser",
		}
		
		// 设置旧密码
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("oldpassword"), bcrypt.DefaultCost)
		user.Password = string(hashedPassword)

		req := dto.UserChangePasswordRequest{
			OldPassword: "oldpassword",
			NewPassword: "newpassword",
		}

		// 设置模拟行为
		mockUserRepo.EXPECT().FindByID(gomock.Any(), userID).Return(user, nil)
		mockUserRepo.EXPECT().Update(gomock.Any(), gomock.Any()).Return(errors.New("database error"))

		// 执行方法
		err := userService.ChangePassword(context.Background(), userID, req)

		// 断言结果
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "database error")
	})
}

func TestUserServiceImpl_GetUserByUsername(t *testing.T) {
	// 创建 gomock 控制器
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// 创建模拟的依赖对象
	mockUserRepo := mocks.NewMockUserRepository(ctrl)
	mockRoleRepo := mocks.NewMockRoleRepository(ctrl)

	// 创建被测试的服务
	userService := service.NewUserService(
		mockUserRepo,
		mockRoleRepo,
		"test-secret",
		time.Hour*24,
	)

	// 测试场景：成功通过用户名获取用户信息
	t.Run("GetUserByUsername_Success", func(t *testing.T) {
		// 准备测试数据
		username := "testuser"
		user := &model.User{
			ID:        "user-123",
			Username:  username,
			Email:     "test@example.com",
			NickName:  "Test User",
			Phone:     "1234567890",
			Status:    model.UserStatusActive,
			CreatedAt: time.Now().Add(-24 * time.Hour),
		}

		// 设置模拟行为
		mockUserRepo.EXPECT().FindByUsername(gomock.Any(), username).Return(user, nil)

		// 执行方法
		userDTO, err := userService.GetUserByUsername(context.Background(), username)

		// 断言结果
		assert.NoError(t, err)
		assert.NotNil(t, userDTO)
		assert.Equal(t, user.ID, userDTO.ID)
		assert.Equal(t, username, userDTO.Username)
		assert.Equal(t, user.Email, userDTO.Email)
		assert.Equal(t, user.NickName, userDTO.NickName)
		assert.Equal(t, user.Phone, userDTO.Phone)
	})

	// 测试场景：用户名不存在
	t.Run("GetUserByUsername_UserNotFound", func(t *testing.T) {
		// 准备测试数据
		username := "nonexistentuser"

		// 设置模拟行为
		mockUserRepo.EXPECT().FindByUsername(gomock.Any(), username).Return(nil, nil)

		// 执行方法
		userDTO, err := userService.GetUserByUsername(context.Background(), username)

		// 断言结果
		assert.Error(t, err)
		assert.Equal(t, "user not found", err.Error())
		assert.Nil(t, userDTO)
	})

	// 测试场景：数据库错误
	t.Run("GetUserByUsername_DatabaseError", func(t *testing.T) {
		// 准备测试数据
		username := "testuser"

		// 设置模拟行为
		mockUserRepo.EXPECT().FindByUsername(gomock.Any(), username).Return(nil, errors.New("database error"))

		// 执行方法
		userDTO, err := userService.GetUserByUsername(context.Background(), username)

		// 断言结果
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "error finding user")
		assert.Nil(t, userDTO)
	})
}

func TestUserServiceImpl_ListUsers(t *testing.T) {
	// 创建 gomock 控制器
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// 创建模拟的依赖对象
	mockUserRepo := mocks.NewMockUserRepository(ctrl)
	mockRoleRepo := mocks.NewMockRoleRepository(ctrl)

	// 创建被测试的服务
	userService := service.NewUserService(
		mockUserRepo,
		mockRoleRepo,
		"test-secret",
		time.Hour*24,
	)

	// 测试场景：成功获取用户列表
	t.Run("ListUsers_Success", func(t *testing.T) {
		// 准备测试数据
		page := 1
		pageSize := 10
		users := []*model.User{
			{
				ID:        "user-1",
				Username:  "user1",
				Email:     "user1@example.com",
				NickName:  "User One",
				Status:    model.UserStatusActive,
				CreatedAt: time.Now().Add(-24 * time.Hour),
			},
			{
				ID:        "user-2",
				Username:  "user2",
				Email:     "user2@example.com",
				NickName:  "User Two",
				Status:    model.UserStatusActive,
				CreatedAt: time.Now().Add(-48 * time.Hour),
			},
		}
		total := int64(2)

		// 设置模拟行为
		mockUserRepo.EXPECT().List(gomock.Any(), page, pageSize).Return(users, total, nil)

		// 执行方法
		response, err := userService.ListUsers(context.Background(), page, pageSize)

		// 断言结果
		assert.NoError(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, total, response.Total)
		assert.Len(t, response.List, len(users))
		
		// 验证返回的用户列表是否正确
		userDTOs, ok := response.List.([]dto.UserDTO)
		assert.True(t, ok, "response.List should be of type []dto.UserDTO")
		for i, user := range users {
			assert.Equal(t, user.ID, userDTOs[i].ID)
			assert.Equal(t, user.Username, userDTOs[i].Username)
			assert.Equal(t, user.Email, userDTOs[i].Email)
		}
	})

	// 测试场景：获取用户列表时数据库错误
	t.Run("ListUsers_DatabaseError", func(t *testing.T) {
		// 准备测试数据
		page := 1
		pageSize := 10

		// 设置模拟行为
		mockUserRepo.EXPECT().List(gomock.Any(), page, pageSize).Return(nil, int64(0), errors.New("database error"))

		// 执行方法
		response, err := userService.ListUsers(context.Background(), page, pageSize)

		// 断言结果
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "error listing users")
		assert.Nil(t, response)
	})

	// 测试场景：用户列表为空
	t.Run("ListUsers_EmptyList", func(t *testing.T) {
		// 准备测试数据
		page := 1
		pageSize := 10
		users := []*model.User{}
		total := int64(0)

		// 设置模拟行为
		mockUserRepo.EXPECT().List(gomock.Any(), page, pageSize).Return(users, total, nil)

		// 执行方法
		response, err := userService.ListUsers(context.Background(), page, pageSize)

		// 断言结果
		assert.NoError(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, total, response.Total)
		assert.Empty(t, response.List)
	})
}

func TestUserServiceImpl_SearchUsers(t *testing.T) {
	// 创建 gomock 控制器
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// 创建模拟的依赖对象
	mockUserRepo := mocks.NewMockUserRepository(ctrl)
	mockRoleRepo := mocks.NewMockRoleRepository(ctrl)

	// 创建被测试的服务
	userService := service.NewUserService(
		mockUserRepo,
		mockRoleRepo,
		"test-secret",
		time.Hour*24,
	)

	// 测试场景：成功搜索用户
	t.Run("SearchUsers_Success", func(t *testing.T) {
		// 准备测试数据
		query := "test"
		page := 1
		pageSize := 10
		users := []*model.User{
			{
				ID:        "user-1",
				Username:  "testuser1",
				Email:     "test1@example.com",
				NickName:  "Test User One",
				Status:    model.UserStatusActive,
				CreatedAt: time.Now().Add(-24 * time.Hour),
			},
			{
				ID:        "user-2",
				Username:  "testuser2",
				Email:     "test2@example.com",
				NickName:  "Test User Two",
				Status:    model.UserStatusActive,
				CreatedAt: time.Now().Add(-48 * time.Hour),
			},
		}
		total := int64(2)

		// 设置模拟行为
		mockUserRepo.EXPECT().Search(gomock.Any(), query, page, pageSize).Return(users, total, nil)

		// 执行方法
		response, err := userService.SearchUsers(context.Background(), query, page, pageSize)

		// 断言结果
		assert.NoError(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, total, response.Total)
		assert.Len(t, response.List, len(users))
		
		// 验证返回的用户列表是否正确
		userDTOs, ok := response.List.([]dto.UserDTO)
		assert.True(t, ok, "response.List should be of type []dto.UserDTO")
		for i, user := range users {
			assert.Equal(t, user.ID, userDTOs[i].ID)
			assert.Equal(t, user.Username, userDTOs[i].Username)
			assert.Equal(t, user.Email, userDTOs[i].Email)
		}
	})

	// 测试场景：搜索用户时数据库错误
	t.Run("SearchUsers_DatabaseError", func(t *testing.T) {
		// 准备测试数据
		query := "test"
		page := 1
		pageSize := 10

		// 设置模拟行为
		mockUserRepo.EXPECT().Search(gomock.Any(), query, page, pageSize).Return(nil, int64(0), errors.New("database error"))

		// 执行方法
		response, err := userService.SearchUsers(context.Background(), query, page, pageSize)

		// 断言结果
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "error searching users")
		assert.Nil(t, response)
	})

	// 测试场景：搜索结果为空
	t.Run("SearchUsers_EmptyResults", func(t *testing.T) {
		// 准备测试数据
		query := "nonexistent"
		page := 1
		pageSize := 10
		users := []*model.User{}
		total := int64(0)

		// 设置模拟行为
		mockUserRepo.EXPECT().Search(gomock.Any(), query, page, pageSize).Return(users, total, nil)

		// 执行方法
		response, err := userService.SearchUsers(context.Background(), query, page, pageSize)

		// 断言结果
		assert.NoError(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, total, response.Total)
		assert.Empty(t, response.List)
	})
}