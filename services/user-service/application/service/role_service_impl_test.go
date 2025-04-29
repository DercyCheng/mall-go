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
			Sort:        1,
			Status:      1,
		}

		// 设置模拟行为
		mockRoleRepo.EXPECT().FindByName(gomock.Any(), req.Name).Return(nil, nil)
		mockRoleRepo.EXPECT().Save(gomock.Any(), gomock.Any()).DoAndReturn(
			func(ctx context.Context, role *model.Role) error {
				// 验证模型转换是否正确
				assert.Equal(t, req.Name, role.Name)
				assert.Equal(t, req.Description, role.Description)
				assert.Equal(t, req.Sort, role.Sort)
				assert.Equal(t, model.RoleStatus(req.Status), role.Status)
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

	// 测试场景：成功获取角色信息
	t.Run("GetRole_Success", func(t *testing.T) {
		// 准备测试数据
		roleID := "role-123"
		role := &model.Role{
			ID:          roleID,
			Name:        "admin",
			Description: "Administrator role",
			Sort:        1,
			Status:      model.RoleStatusEnabled,
			CreatedAt:   time.Now().Add(-24 * time.Hour),
		}

		// 设置模拟行为
		mockRoleRepo.EXPECT().FindByID(gomock.Any(), roleID).Return(role, nil)

		// 执行方法
		roleDTO, err := roleService.GetRole(context.Background(), roleID)

		// 断言结果
		assert.NoError(t, err)
		assert.NotNil(t, roleDTO)
		assert.Equal(t, roleID, roleDTO.ID)
		assert.Equal(t, role.Name, roleDTO.Name)
		assert.Equal(t, role.Description, roleDTO.Description)
		assert.Equal(t, int(role.Status), roleDTO.Status)
	})

	// 测试场景：角色不存在
	t.Run("GetRole_RoleNotFound", func(t *testing.T) {
		// 准备测试数据
		roleID := "nonexistent-role"

		// 设置模拟行为
		mockRoleRepo.EXPECT().FindByID(gomock.Any(), roleID).Return(nil, nil)

		// 执行方法
		roleDTO, err := roleService.GetRole(context.Background(), roleID)

		// 断言结果
		assert.Error(t, err)
		assert.Equal(t, "role not found", err.Error())
		assert.Nil(t, roleDTO)
	})

	// 测试场景：数据库错误
	t.Run("GetRole_DatabaseError", func(t *testing.T) {
		// 准备测试数据
		roleID := "role-123"

		// 设置模拟行为
		mockRoleRepo.EXPECT().FindByID(gomock.Any(), roleID).Return(nil, errors.New("database connection error"))

		// 执行方法
		roleDTO, err := roleService.GetRole(context.Background(), roleID)

		// 断言结果
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "error finding role")
		assert.Nil(t, roleDTO)
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

	// 测试场景：成功更新角色信息
	t.Run("UpdateRole_Success", func(t *testing.T) {
		// 准备测试数据
		roleID := "role-123"
		role := &model.Role{
			ID:          roleID,
			Name:        "old_name",
			Description: "Old description",
			Sort:        1,
			Status:      model.RoleStatusEnabled,
		}

		updateReq := dto.RoleUpdateRequest{
			Name:        "new_name",
			Description: "New description",
			Sort:        2,
			Status:      int(model.RoleStatusDisabled),
		}

		// 设置模拟行为
		mockRoleRepo.EXPECT().FindByID(gomock.Any(), roleID).Return(role, nil)
		mockRoleRepo.EXPECT().Update(gomock.Any(), gomock.Any()).DoAndReturn(
			func(ctx context.Context, updatedRole *model.Role) error {
				// 验证更新是否应用
				assert.Equal(t, updateReq.Name, updatedRole.Name)
				assert.Equal(t, updateReq.Description, updatedRole.Description)
				assert.Equal(t, updateReq.Sort, updatedRole.Sort)
				assert.Equal(t, model.RoleStatus(updateReq.Status), updatedRole.Status)
				return nil
			})

		// 执行方法
		err := roleService.UpdateRole(context.Background(), roleID, updateReq)

		// 断言结果
		assert.NoError(t, err)
	})

	// 测试场景：角色不存在
	t.Run("UpdateRole_RoleNotFound", func(t *testing.T) {
		// 准备测试数据
		roleID := "nonexistent-role"
		updateReq := dto.RoleUpdateRequest{
			Name: "new_name",
		}

		// 设置模拟行为
		mockRoleRepo.EXPECT().FindByID(gomock.Any(), roleID).Return(nil, nil)

		// 执行方法
		err := roleService.UpdateRole(context.Background(), roleID, updateReq)

		// 断言结果
		assert.Error(t, err)
		assert.Equal(t, "role not found", err.Error())
	})

	// 测试场景：更新数据库失败
	t.Run("UpdateRole_DatabaseError", func(t *testing.T) {
		// 准备测试数据
		roleID := "role-123"
		role := &model.Role{
			ID:   roleID,
			Name: "old_name",
		}

		updateReq := dto.RoleUpdateRequest{
			Name: "new_name",
		}

		// 设置模拟行为
		mockRoleRepo.EXPECT().FindByID(gomock.Any(), roleID).Return(role, nil)
		mockRoleRepo.EXPECT().Update(gomock.Any(), gomock.Any()).Return(errors.New("database error"))

		// 执行方法
		err := roleService.UpdateRole(context.Background(), roleID, updateReq)

		// 断言结果
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "error updating role")
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
		// 准备测试数据
		roleID := "role-123"
		role := &model.Role{
			ID:   roleID,
			Name: "admin",
		}

		// 设置模拟行为
		mockRoleRepo.EXPECT().FindByID(gomock.Any(), roleID).Return(role, nil)
		mockRoleRepo.EXPECT().Delete(gomock.Any(), roleID).Return(nil)

		// 执行方法
		err := roleService.DeleteRole(context.Background(), roleID)

		// 断言结果
		assert.NoError(t, err)
	})

	// 测试场景：角色不存在
	t.Run("DeleteRole_RoleNotFound", func(t *testing.T) {
		// 准备测试数据
		roleID := "nonexistent-role"

		// 设置模拟行为
		mockRoleRepo.EXPECT().FindByID(gomock.Any(), roleID).Return(nil, nil)

		// 执行方法
		err := roleService.DeleteRole(context.Background(), roleID)

		// 断言结果
		assert.Error(t, err)
		assert.Equal(t, "role not found", err.Error())
	})

	// 测试场景：删除数据库失败
	t.Run("DeleteRole_DatabaseError", func(t *testing.T) {
		// 准备测试数据
		roleID := "role-123"
		role := &model.Role{
			ID:   roleID,
			Name: "admin",
		}

		// 设置模拟行为
		mockRoleRepo.EXPECT().FindByID(gomock.Any(), roleID).Return(role, nil)
		mockRoleRepo.EXPECT().Delete(gomock.Any(), roleID).Return(errors.New("database error"))

		// 执行方法
		err := roleService.DeleteRole(context.Background(), roleID)

		// 断言结果
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
		// 准备测试数据
		roles := []model.Role{
			{
				ID:          "role-1",
				Name:        "admin",
				Description: "Administrator role",
				Sort:        1,
				Status:      model.RoleStatusEnabled,
			},
			{
				ID:          "role-2",
				Name:        "user",
				Description: "Regular user role",
				Sort:        2,
				Status:      model.RoleStatusEnabled,
			},
		}

		// 设置模拟行为
		mockRoleRepo.EXPECT().List(gomock.Any()).Return(roles, nil)

		// 执行方法
		roleDTOs, err := roleService.ListRoles(context.Background())

		// 断言结果
		assert.NoError(t, err)
		assert.Len(t, roleDTOs, 2)
		assert.Equal(t, roles[0].ID, roleDTOs[0].ID)
		assert.Equal(t, roles[0].Name, roleDTOs[0].Name)
		assert.Equal(t, roles[1].ID, roleDTOs[1].ID)
		assert.Equal(t, roles[1].Name, roleDTOs[1].Name)
	})

	// 测试场景：数据库错误
	t.Run("ListRoles_DatabaseError", func(t *testing.T) {
		// 设置模拟行为
		mockRoleRepo.EXPECT().List(gomock.Any()).Return(nil, errors.New("database error"))

		// 执行方法
		roleDTOs, err := roleService.ListRoles(context.Background())

		// 断言结果
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "error listing roles")
		assert.Nil(t, roleDTOs)
	})

	// 测试场景：角色列表为空
	t.Run("ListRoles_EmptyList", func(t *testing.T) {
		// 设置模拟行为
		mockRoleRepo.EXPECT().List(gomock.Any()).Return([]model.Role{}, nil)

		// 执行方法
		roleDTOs, err := roleService.ListRoles(context.Background())

		// 断言结果
		assert.NoError(t, err)
		assert.Empty(t, roleDTOs)
	})
}

func TestRoleServiceImpl_AssignRolesToUser(t *testing.T) {
	// 创建 gomock 控制器
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// 创建模拟的依赖对象
	mockRoleRepo := mocks.NewMockRoleRepository(ctrl)
	mockUserRepo := mocks.NewMockUserRepository(ctrl)

	// 创建被测试的服务
	roleService := service.NewRoleServiceWithUserRepo(mockRoleRepo, mockUserRepo)

	// 测试场景：成功分配角色给用户
	t.Run("AssignRolesToUser_Success", func(t *testing.T) {
		// 准备测试数据
		userID := "user-123"
		req := dto.AssignRoleRequest{
			UserID:  userID,
			RoleIDs: []string{"role-1", "role-2"},
		}
		user := &model.User{
			ID:       userID,
			Username: "testuser",
		}

		// 设置模拟行为
		mockUserRepo.EXPECT().FindByID(gomock.Any(), userID).Return(user, nil)
		
		// 先清除已有角色
		mockRoleRepo.EXPECT().ClearUserRoles(gomock.Any(), userID).Return(nil)
		
		// 分别分配两个角色
		mockRoleRepo.EXPECT().AssignRoleToUser(gomock.Any(), userID, "role-1").Return(nil)
		mockRoleRepo.EXPECT().AssignRoleToUser(gomock.Any(), userID, "role-2").Return(nil)

		// 执行方法
		err := roleService.AssignRolesToUser(context.Background(), req)

		// 断言结果
		assert.NoError(t, err)
	})

	// 测试场景：用户不存在
	t.Run("AssignRolesToUser_UserNotFound", func(t *testing.T) {
		// 准备测试数据
		userID := "nonexistent-user"
		req := dto.AssignRoleRequest{
			UserID:  userID,
			RoleIDs: []string{"role-1"},
		}

		// 设置模拟行为
		mockUserRepo.EXPECT().FindByID(gomock.Any(), userID).Return(nil, nil)

		// 执行方法
		err := roleService.AssignRolesToUser(context.Background(), req)

		// 断言结果
		assert.Error(t, err)
		assert.Equal(t, "user not found", err.Error())
	})

	// 测试场景：清除角色时数据库错误
	t.Run("AssignRolesToUser_ClearRolesError", func(t *testing.T) {
		// 准备测试数据
		userID := "user-123"
		req := dto.AssignRoleRequest{
			UserID:  userID,
			RoleIDs: []string{"role-1"},
		}
		user := &model.User{
			ID:       userID,
			Username: "testuser",
		}

		// 设置模拟行为
		mockUserRepo.EXPECT().FindByID(gomock.Any(), userID).Return(user, nil)
		mockRoleRepo.EXPECT().ClearUserRoles(gomock.Any(), userID).Return(errors.New("database error"))

		// 执行方法
		err := roleService.AssignRolesToUser(context.Background(), req)

		// 断言结果
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "error clearing user roles")
	})

	// 测试场景：分配角色时数据库错误
	t.Run("AssignRolesToUser_AssignRoleError", func(t *testing.T) {
		// 准备测试数据
		userID := "user-123"
		req := dto.AssignRoleRequest{
			UserID:  userID,
			RoleIDs: []string{"role-1", "role-2"},
		}
		user := &model.User{
			ID:       userID,
			Username: "testuser",
		}

		// 设置模拟行为
		mockUserRepo.EXPECT().FindByID(gomock.Any(), userID).Return(user, nil)
		mockRoleRepo.EXPECT().ClearUserRoles(gomock.Any(), userID).Return(nil)
		mockRoleRepo.EXPECT().AssignRoleToUser(gomock.Any(), userID, "role-1").Return(nil)
		mockRoleRepo.EXPECT().AssignRoleToUser(gomock.Any(), userID, "role-2").Return(errors.New("database error"))

		// 执行方法
		err := roleService.AssignRolesToUser(context.Background(), req)

		// 断言结果
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "error assigning role")
	})
}

func TestRoleServiceImpl_GetUserRoles(t *testing.T) {
	// 创建 gomock 控制器
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// 创建模拟的依赖对象
	mockRoleRepo := mocks.NewMockRoleRepository(ctrl)
	mockUserRepo := mocks.NewMockUserRepository(ctrl)

	// 创建被测试的服务
	roleService := service.NewRoleServiceWithUserRepo(mockRoleRepo, mockUserRepo)

	// 测试场景：成功获取用户角色
	t.Run("GetUserRoles_Success", func(t *testing.T) {
		// 准备测试数据
		userID := "user-123"
		user := &model.User{
			ID:       userID,
			Username: "testuser",
		}
		roles := []model.Role{
			{
				ID:          "role-1",
				Name:        "admin",
				Description: "Administrator role",
			},
			{
				ID:          "role-2",
				Name:        "user",
				Description: "Regular user role",
			},
		}

		// 设置模拟行为
		mockUserRepo.EXPECT().FindByID(gomock.Any(), userID).Return(user, nil)
		mockRoleRepo.EXPECT().GetUserRoles(gomock.Any(), userID).Return(roles, nil)

		// 执行方法
		roleDTOs, err := roleService.GetUserRoles(context.Background(), userID)

		// 断言结果
		assert.NoError(t, err)
		assert.Len(t, roleDTOs, 2)
		assert.Equal(t, roles[0].ID, roleDTOs[0].ID)
		assert.Equal(t, roles[0].Name, roleDTOs[0].Name)
		assert.Equal(t, roles[1].ID, roleDTOs[1].ID)
		assert.Equal(t, roles[1].Name, roleDTOs[1].Name)
	})

	// 测试场景：用户不存在
	t.Run("GetUserRoles_UserNotFound", func(t *testing.T) {
		// 准备测试数据
		userID := "nonexistent-user"

		// 设置模拟行为
		mockUserRepo.EXPECT().FindByID(gomock.Any(), userID).Return(nil, nil)

		// 执行方法
		roleDTOs, err := roleService.GetUserRoles(context.Background(), userID)

		// 断言结果
		assert.Error(t, err)
		assert.Equal(t, "user not found", err.Error())
		assert.Nil(t, roleDTOs)
	})

	// 测试场景：获取角色时数据库错误
	t.Run("GetUserRoles_DatabaseError", func(t *testing.T) {
		// 准备测试数据
		userID := "user-123"
		user := &model.User{
			ID:       userID,
			Username: "testuser",
		}

		// 设置模拟行为
		mockUserRepo.EXPECT().FindByID(gomock.Any(), userID).Return(user, nil)
		mockRoleRepo.EXPECT().GetUserRoles(gomock.Any(), userID).Return(nil, errors.New("database error"))

		// 执行方法
		roleDTOs, err := roleService.GetUserRoles(context.Background(), userID)

		// 断言结果
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "error getting user roles")
		assert.Nil(t, roleDTOs)
	})

	// 测试场景：用户没有角色
	t.Run("GetUserRoles_NoRoles", func(t *testing.T) {
		// 准备测试数据
		userID := "user-123"
		user := &model.User{
			ID:       userID,
			Username: "testuser",
		}

		// 设置模拟行为
		mockUserRepo.EXPECT().FindByID(gomock.Any(), userID).Return(user, nil)
		mockRoleRepo.EXPECT().GetUserRoles(gomock.Any(), userID).Return([]model.Role{}, nil)

		// 执行方法
		roleDTOs, err := roleService.GetUserRoles(context.Background(), userID)

		// 断言结果
		assert.NoError(t, err)
		assert.Empty(t, roleDTOs)
	})
}