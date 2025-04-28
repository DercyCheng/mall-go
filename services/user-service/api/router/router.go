package router

import (
	"github.com/gin-gonic/gin"

	"mall-go/services/user-service/api/handler"
	"mall-go/services/user-service/api/middleware"
)

// SetupRouter 设置路由
func SetupRouter(r *gin.Engine, userHandler *handler.UserHandler, roleHandler *handler.RoleHandler, permissionHandler *handler.PermissionHandler) {
	// 健康检查接口
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
		})
	})

	// API分组
	api := r.Group("/api")
	{
		// 公开接口
		api.POST("/users/register", userHandler.Register)
		api.POST("/users/login", userHandler.Login)

		// 受保护接口
		authenticated := api.Group("/")
		authenticated.Use(middleware.JWTAuth())
		{
			// 用户相关接口
			users := authenticated.Group("/users")
			{
				users.GET("/info", userHandler.GetUserInfo)       // 获取当前用户信息
				users.PUT("/info", userHandler.UpdateUserInfo)    // 更新当前用户信息
				users.PUT("/password", userHandler.UpdatePassword) // 修改密码
				users.GET("", userHandler.ListUsers)              // 获取用户列表
				users.GET("/search", userHandler.SearchUsers)     // 搜索用户
				users.GET("/:id", userHandler.GetUser)            // 获取指定用户信息
				users.PATCH("/:id/status", userHandler.UpdateStatus) // 更新用户状态
				users.POST("/:id/roles", userHandler.AssignRoles) // 分配角色
				users.DELETE("/:id", userHandler.DeleteUser)      // 删除用户
			}

			// 角色相关接口
			roles := authenticated.Group("/roles")
			{
				roles.POST("", roleHandler.CreateRole)              // 创建角色
				roles.GET("", roleHandler.ListRoles)                // 获取角色列表
				roles.GET("/:id", roleHandler.GetRole)              // 获取角色详情
				roles.PUT("/:id", roleHandler.UpdateRole)           // 更新角色
				roles.DELETE("/:id", roleHandler.DeleteRole)        // 删除角色
				roles.PATCH("/:id/status", roleHandler.UpdateRoleStatus) // 更新角色状态
				roles.POST("/:id/permissions", roleHandler.AssignPermissions) // 分配权限
				roles.GET("/:id/permissions", roleHandler.GetRolePermissions) // 获取角色权限
			}
			
			// 权限相关接口
			permissions := authenticated.Group("/permissions")
			{
				permissions.POST("", permissionHandler.CreatePermission)              // 创建权限
				permissions.GET("", permissionHandler.ListPermissions)                // 获取权限列表
				permissions.GET("/:id", permissionHandler.GetPermission)              // 获取权限详情
				permissions.PUT("/:id", permissionHandler.UpdatePermission)           // 更新权限
				permissions.DELETE("/:id", permissionHandler.DeletePermission)        // 删除权限
				permissions.PATCH("/:id/status", permissionHandler.UpdatePermissionStatus) // 更新权限状态
				permissions.POST("/roles/:id/assign", permissionHandler.AssignPermissionsToRole) // 分配权限到角色
			}
		}
	}
}