package service_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"mall-go/services/user-service/application/dto"
	"mall-go/services/user-service/application/service"
	"mall-go/services/user-service/domain/model"
	"mall-go/services/user-service/mocks"
)

func TestRoleServiceImpl_CreateRole(t *testing.T) {
	// 创建 gomock 控制器
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// 创建模拟的依赖对象
	mockRoleRepo := mocks.NewMockRoleRepository(ctrl)

	// 创建被测试的服务
	roleService := service.NewRoleService(mockRoleRepo)

	// 测试场景：成功创建角色
	t.Run("CreateRole_Success", func(t *testing.T) {
		// 设置输入
		req := dto.RoleCreateRequest{
			Name:        "admin",
			Description: "Administrator role",
		}

		// 设置模拟行为
		mockRoleRepo.EXPECT().FindByName(gomock.Any(), req.Name).Return(nil, nil)
		mockRoleRepo.EXPECT().Save(gomock.Any(), gomock.Any()).DoAndReturn(
			func(ctx context.Context, role *model.Role) error {
				// 验证模型转换是否正确
				assert.Equal(t, req.Name, role.Name)
				assert.Equal(t, req.Description, role.Description)
				assert.NotEmpty(t, role.ID)
				return nil
			})

		// 执行方法
		roleID, err := roleService.CreateRole(context.Background(), req)

		// 断言结果
		assert.NoError(t, err)
		assert.NotEmpty(t, roleID)
	})

	// 测试场景：角色名已存在
	t.Run("CreateRole_NameAlreadyExists", func(t *testing.T) {
		// 设置输入
		req := dto.RoleCreateRequest{
			Name:        "existing_role",
			Description: "Role description",
		}

		// 设置模拟行为 - 模拟角色名已存在
		existingRole := &model.Role{
			ID:   "existing-id",
			Name: "existing_role",
		}
		mockRoleRepo.EXPECT().FindByName(gomock.Any(), req.Name).Return(existingRole, nil)

		// 执行方法
		roleID, err := roleService.CreateRole(context.Background(), req)

		// 断言结果
		assert.Error(t, err)
		assert.Equal(t, "role name already exists", err.Error())
		assert.Empty(t, roleID)
	})

	// 测试场景：保存角色时发生错误
	t.Run("CreateRole_SaveError", func(t *testing.T) {
		// 设置输入
		req := dto.RoleCreateRequest{
			Name:        "new_role",
			Description: "Role description",
		}

		// 设置模拟行为
		mockRoleRepo.EXPECT().FindByName(gomock.Any(), req.Name).Return(nil, nil)
		mockRoleRepo.EXPECT().Save(gomock.Any(), gomock.Any()).Return(errors.New("database error"))

		// 执行方法
		roleID, err := roleService.CreateRole(context.Background(), req)

		// 断言结果
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "error saving role")
		assert.Empty(t, roleID)
	})
}

func TestRoleServiceImpl_GetRole(t *testing.T) {
	// 创建 gomock 控制器
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// 创建模拟的依赖对象
	mockRoleRepo := mocks.NewMockRoleRepository(ctrl)

	// 创建被测试的服务
	roleService := service.NewRoleService(mockRoleRepo)

	// 测试场景：成功获取角色
	t.Run("GetRole_Success", func(t *testing.T) {
		roleID := "role-123"
		expectedRole := &model.Role{
			ID:          roleID,
			Name:        "admin",
			Description: "Administrator role",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}

		mockRoleRepo.EXPECT().FindByID(gomock.Any(), roleID).Return(expectedRole, nil)

		result, err := roleService.GetRole(context.Background(), roleID)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, roleID, result.ID)
		assert.Equal(t, "admin", result.Name)
		assert.Equal(t, "Administrator role", result.Description)
	})

	// 测试场景：角色不存在
	t.Run("GetRole_NotFound", func(t *testing.T) {
		roleID := "nonexistent"
		mockRoleRepo.EXPECT().FindByID(gomock.Any(), roleID).Return(nil, nil)

		result, err := roleService.GetRole(context.Background(), roleID)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, "role not found", err.Error())
	})

	// 测试场景：数据库错误
	t.Run("GetRole_DatabaseError", func(t *testing.T) {
		roleID := "role-123"
		mockRoleRepo.EXPECT().FindByID(gomock.Any(), roleID).Return(nil, errors.New("database error"))

		result, err := roleService.GetRole(context.Background(), roleID)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "error finding role")
	})
}

func TestRoleServiceImpl_UpdateRole(t *testing.T) {
	// 创建 gomock 控制器
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// 创建模拟的依赖对象
	mockRoleRepo := mocks.NewMockRoleRepository(ctrl)

	// 创建被测试的服务
	roleService := service.NewRoleService(mockRoleRepo)

	// 测试场景：成功更新角色
	t.Run("UpdateRole_Success", func(t *testing.T) {
		roleID := "role-123"
		req := dto.RoleUpdateRequest{
			Name:        "updated-role",
			Description: "Updated role description",
		}

		// 使用较早的时间点，以确保任何新的更新时间都会比这个时间晚
		oldTime := time.Now().Add(-24 * time.Hour)
		currentRole := &model.Role{
			ID:          roleID,
			Name:        "old-role",
			Description: "Old description",
			CreatedAt:   oldTime,
			UpdatedAt:   oldTime,
		}

		mockRoleRepo.EXPECT().FindByID(gomock.Any(), roleID).Return(currentRole, nil)
		mockRoleRepo.EXPECT().FindByName(gomock.Any(), req.Name).Return(nil, nil)
		mockRoleRepo.EXPECT().Update(gomock.Any(), gomock.Any()).DoAndReturn(
			func(ctx context.Context, role *model.Role) error {
				// 验证更新是否正确应用
				assert.Equal(t, roleID, role.ID)
				assert.Equal(t, req.Name, role.Name)
				assert.Equal(t, req.Description, role.Description)
				// 确保更新时间是在旧时间之后，而不是硬编码比较
				assert.True(t, role.UpdatedAt.After(oldTime), "UpdatedAt should be after the original time")
				return nil
			})

		err := roleService.UpdateRole(context.Background(), roleID, req)

		assert.NoError(t, err)
	})

	// 测试场景：角色不存在
	t.Run("UpdateRole_RoleNotFound", func(t *testing.T) {
		roleID := "nonexistent"
		req := dto.RoleUpdateRequest{
			Name:        "updated-role",
			Description: "Updated description",
		}

		mockRoleRepo.EXPECT().FindByID(gomock.Any(), roleID).Return(nil, nil)

		err := roleService.UpdateRole(context.Background(), roleID, req)

		assert.Error(t, err)
		assert.Equal(t, "role not found", err.Error())
	})

	// 测试场景：角色名冲突
	t.Run("UpdateRole_NameConflict", func(t *testing.T) {
		roleID := "role-123"
		req := dto.RoleUpdateRequest{
			Name:        "existing-role",
			Description: "Updated description",
		}

		currentRole := &model.Role{
			ID:          roleID,
			Name:        "old-role",
			Description: "Old description",
			CreatedAt:   time.Now().Add(-24 * time.Hour),
			UpdatedAt:   time.Now().Add(-24 * time.Hour),
		}

		existingRole := &model.Role{
			ID:          "role-456",
			Name:        "existing-role",
			Description: "Another role",
		}

		mockRoleRepo.EXPECT().FindByID(gomock.Any(), roleID).Return(currentRole, nil)
		mockRoleRepo.EXPECT().FindByName(gomock.Any(), req.Name).Return(existingRole, nil)

		err := roleService.UpdateRole(context.Background(), roleID, req)

		assert.Error(t, err)
		assert.Equal(t, "role name already exists", err.Error())
	})
}

func TestRoleServiceImpl_DeleteRole(t *testing.T) {
	// 创建 gomock 控制器
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// 创建模拟的依赖对象
	mockRoleRepo := mocks.NewMockRoleRepository(ctrl)

	// 创建被测试的服务
	roleService := service.NewRoleService(mockRoleRepo)

	// 测试场景：成功删除角色
	t.Run("DeleteRole_Success", func(t *testing.T) {
		roleID := "role-123"
		role := &model.Role{
			ID:   roleID,
			Name: "test-role",
		}

		mockRoleRepo.EXPECT().FindByID(gomock.Any(), roleID).Return(role, nil)
		mockRoleRepo.EXPECT().Delete(gomock.Any(), roleID).Return(nil)

		err := roleService.DeleteRole(context.Background(), roleID)

		assert.NoError(t, err)
	})

	// 测试场景：角色不存在
	t.Run("DeleteRole_NotFound", func(t *testing.T) {
		roleID := "nonexistent"
		mockRoleRepo.EXPECT().FindByID(gomock.Any(), roleID).Return(nil, nil)

		err := roleService.DeleteRole(context.Background(), roleID)

		assert.Error(t, err)
		assert.Equal(t, "role not found", err.Error())
	})

	// 测试场景：删除角色时发生错误
	t.Run("DeleteRole_Error", func(t *testing.T) {
		roleID := "role-123"
		role := &model.Role{
			ID:   roleID,
			Name: "test-role",
		}

		mockRoleRepo.EXPECT().FindByID(gomock.Any(), roleID).Return(role, nil)
		mockRoleRepo.EXPECT().Delete(gomock.Any(), roleID).Return(errors.New("database error"))

		err := roleService.DeleteRole(context.Background(), roleID)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "error deleting role")
	})
}

func TestRoleServiceImpl_ListRoles(t *testing.T) {
	// 创建 gomock 控制器
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// 创建模拟的依赖对象
	mockRoleRepo := mocks.NewMockRoleRepository(ctrl)

	// 创建被测试的服务
	roleService := service.NewRoleService(mockRoleRepo)

	// 测试场景：成功获取角色列表
	t.Run("ListRoles_Success", func(t *testing.T) {
		roles := []*model.Role{
			{
				ID:          "role-1",
				Name:        "admin",
				Description: "Administrator role",
				CreatedAt:   time.Now(),
			},
			{
				ID:          "role-2",
				Name:        "user",
				Description: "Regular user role",
				CreatedAt:   time.Now(),
			},
		}

		mockRoleRepo.EXPECT().List(gomock.Any()).Return(roles, nil)

		result, err := roleService.ListRoles(context.Background())

		assert.NoError(t, err)
		assert.Len(t, result, 2)
		assert.Equal(t, "role-1", result[0].ID)
		assert.Equal(t, "admin", result[0].Name)
		assert.Equal(t, "role-2", result[1].ID)
		assert.Equal(t, "user", result[1].Name)
	})

	// 测试场景：返回空列表
	t.Run("ListRoles_Empty", func(t *testing.T) {
		mockRoleRepo.EXPECT().List(gomock.Any()).Return([]*model.Role{}, nil)

		result, err := roleService.ListRoles(context.Background())

		assert.NoError(t, err)
		assert.Empty(t, result)
	})

	// 测试场景：数据库错误
	t.Run("ListRoles_DatabaseError", func(t *testing.T) {
		mockRoleRepo.EXPECT().List(gomock.Any()).Return(nil, errors.New("database error"))

		result, err := roleService.ListRoles(context.Background())

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "error listing roles")
	})
}

func TestRoleServiceImpl_AssignRolesToUser(t *testing.T) {
	// 创建 gomock 控制器
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// 创建模拟的依赖对象
	mockRoleRepo := mocks.NewMockRoleRepository(ctrl)

	// 创建被测试的服务
	roleService := service.NewRoleService(mockRoleRepo)

	// 测试场景：为用户分配新角色
	t.Run("AssignRolesToUser_NewRoles", func(t *testing.T) {
		userID := "user-123"
		req := dto.AssignRoleRequest{
			UserID:  userID,
			RoleIds: []string{"role-1", "role-2"},
		}

		// 模拟用户当前没有角色
		mockRoleRepo.EXPECT().GetUserRoles(gomock.Any(), userID).Return([]*model.Role{}, nil)
		
		// 期望分配两个新角色
		mockRoleRepo.EXPECT().AssignRoleToUser(gomock.Any(), userID, "role-1").Return(nil)
		mockRoleRepo.EXPECT().AssignRoleToUser(gomock.Any(), userID, "role-2").Return(nil)

		err := roleService.AssignRolesToUser(context.Background(), req)

		assert.NoError(t, err)
	})

	// 测试场景：更新用户角色（添加新角色，移除已有角色）
	t.Run("AssignRolesToUser_UpdateRoles", func(t *testing.T) {
		userID := "user-123"
		req := dto.AssignRoleRequest{
			UserID:  userID,
			RoleIds: []string{"role-1", "role-3"}, // 保留role-1，添加role-3，移除role-2
		}

		// 模拟用户当前有角色
		currentRoles := []*model.Role{
			{ID: "role-1", Name: "admin"},
			{ID: "role-2", Name: "user"},
		}
		mockRoleRepo.EXPECT().GetUserRoles(gomock.Any(), userID).Return(currentRoles, nil)
		
		// 期望添加新角色
		mockRoleRepo.EXPECT().AssignRoleToUser(gomock.Any(), userID, "role-3").Return(nil)
		
		// 期望移除已有角色
		mockRoleRepo.EXPECT().RevokeRoleFromUser(gomock.Any(), userID, "role-2").Return(nil)

		err := roleService.AssignRolesToUser(context.Background(), req)

		assert.NoError(t, err)
	})

	// 测试场景：获取用户角色时发生错误
	t.Run("AssignRolesToUser_GetRolesError", func(t *testing.T) {
		userID := "user-123"
		req := dto.AssignRoleRequest{
			UserID:  userID,
			RoleIds: []string{"role-1"},
		}

		mockRoleRepo.EXPECT().GetUserRoles(gomock.Any(), userID).Return(nil, errors.New("database error"))

		err := roleService.AssignRolesToUser(context.Background(), req)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "error getting user roles")
	})
}

func TestRoleServiceImpl_GetUserRoles(t *testing.T) {
	// 创建 gomock 控制器
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// 创建模拟的依赖对象
	mockRoleRepo := mocks.NewMockRoleRepository(ctrl)

	// 创建被测试的服务
	roleService := service.NewRoleService(mockRoleRepo)

	// 测试场景：成功获取用户角色
	t.Run("GetUserRoles_Success", func(t *testing.T) {
		userID := "user-123"
		roles := []*model.Role{
			{
				ID:          "role-1",
				Name:        "admin",
				Description: "Administrator role",
				CreatedAt:   time.Now(),
			},
			{
				ID:          "role-2",
				Name:        "user",
				Description: "Regular user role",
				CreatedAt:   time.Now(),
			},
		}

		mockRoleRepo.EXPECT().GetUserRoles(gomock.Any(), userID).Return(roles, nil)

		result, err := roleService.GetUserRoles(context.Background(), userID)

		assert.NoError(t, err)
		assert.Len(t, result, 2)
		assert.Equal(t, "role-1", result[0].ID)
		assert.Equal(t, "admin", result[0].Name)
		assert.Equal(t, "role-2", result[1].ID)
		assert.Equal(t, "user", result[1].Name)
	})

	// 测试场景：用户没有角色
	t.Run("GetUserRoles_NoRoles", func(t *testing.T) {
		userID := "user-123"
		mockRoleRepo.EXPECT().GetUserRoles(gomock.Any(), userID).Return([]*model.Role{}, nil)

		result, err := roleService.GetUserRoles(context.Background(), userID)

		assert.NoError(t, err)
		assert.Empty(t, result)
	})

	// 测试场景：获取角色时发生错误
	t.Run("GetUserRoles_Error", func(t *testing.T) {
		userID := "user-123"
		mockRoleRepo.EXPECT().GetUserRoles(gomock.Any(), userID).Return(nil, errors.New("database error"))

		result, err := roleService.GetUserRoles(context.Background(), userID)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "error getting user roles")
	})
}