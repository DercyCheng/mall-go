package handler

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"

	"mall-go/pkg/response"
	"mall-go/services/user-service/application/dto"
	"mall-go/services/user-service/application/service"
)

// PermissionHandler 权限处理器
type PermissionHandler struct {
	permissionService *service.PermissionService
}

// NewPermissionHandler 创建权限处理器
func NewPermissionHandler(permissionService *service.PermissionService) *PermissionHandler {
	return &PermissionHandler{
		permissionService: permissionService,
	}
}

// CreatePermission 创建权限
// @Summary 创建权限
// @Description 创建新的权限
// @Tags 权限管理
// @Accept json
// @Produce json
// @Param request body dto.CreatePermissionRequest true "创建权限请求"
// @Success 200 {object} response.Response{data=dto.PermissionResponse} "成功"
// @Failure 400 {object} response.Response "请求错误"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /api/v1/permissions [post]
func (h *PermissionHandler) CreatePermission(c *gin.Context) {
	var req dto.CreatePermissionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, "Invalid request: "+err.Error())
		return
	}

	// 验证请求
	if err := req.Validate(); err != nil {
		response.Fail(c, http.StatusBadRequest, err.Error())
		return
	}

	// 创建权限
	perm, err := h.permissionService.CreatePermission(c.Request.Context(), &req)
	if err != nil {
		switch err {
		case service.ErrPermissionValueAlreadyExists:
			response.Fail(c, http.StatusBadRequest, "Permission value already exists")
		default:
			response.Fail(c, http.StatusInternalServerError, "Failed to create permission: "+err.Error())
		}
		return
	}

	response.Success(c, "Permission created successfully", perm)
}

// GetPermission 获取权限
// @Summary 获取权限
// @Description 根据ID获取权限详情
// @Tags 权限管理
// @Accept json
// @Produce json
// @Param id path string true "权限ID"
// @Success 200 {object} response.Response{data=dto.PermissionResponse} "成功"
// @Failure 400 {object} response.Response "请求错误"
// @Failure 404 {object} response.Response "权限不存在"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /api/v1/permissions/{id} [get]
func (h *PermissionHandler) GetPermission(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.Fail(c, http.StatusBadRequest, "Missing permission ID")
		return
	}

	// 获取权限
	perm, err := h.permissionService.GetPermission(c.Request.Context(), id)
	if err != nil {
		response.Fail(c, http.StatusNotFound, "Permission not found: "+err.Error())
		return
	}

	response.Success(c, "Permission retrieved successfully", perm)
}

// UpdatePermission 更新权限
// @Summary 更新权限
// @Description 更新权限信息
// @Tags 权限管理
// @Accept json
// @Produce json
// @Param id path string true "权限ID"
// @Param request body dto.UpdatePermissionRequest true "更新权限请求"
// @Success 200 {object} response.Response{data=dto.PermissionResponse} "成功"
// @Failure 400 {object} response.Response "请求错误"
// @Failure 404 {object} response.Response "权限不存在"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /api/v1/permissions/{id} [put]
func (h *PermissionHandler) UpdatePermission(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.Fail(c, http.StatusBadRequest, "Missing permission ID")
		return
	}

	var req dto.UpdatePermissionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, "Invalid request: "+err.Error())
		return
	}

	// 验证请求
	if err := req.Validate(); err != nil {
		response.Fail(c, http.StatusBadRequest, err.Error())
		return
	}

	// 更新权限
	perm, err := h.permissionService.UpdatePermission(c.Request.Context(), id, &req)
	if err != nil {
		switch err {
		case service.ErrPermissionValueAlreadyExists:
			response.Fail(c, http.StatusBadRequest, "Permission value already exists")
		default:
			response.Fail(c, http.StatusInternalServerError, "Failed to update permission: "+err.Error())
		}
		return
	}

	response.Success(c, "Permission updated successfully", perm)
}

// DeletePermission 删除权限
// @Summary 删除权限
// @Description 删除权限
// @Tags 权限管理
// @Accept json
// @Produce json
// @Param id path string true "权限ID"
// @Success 200 {object} response.Response "成功"
// @Failure 400 {object} response.Response "请求错误"
// @Failure 404 {object} response.Response "权限不存在"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /api/v1/permissions/{id} [delete]
func (h *PermissionHandler) DeletePermission(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.Fail(c, http.StatusBadRequest, "Missing permission ID")
		return
	}

	// 删除权限
	if err := h.permissionService.DeletePermission(c.Request.Context(), id); err != nil {
		response.Fail(c, http.StatusInternalServerError, "Failed to delete permission: "+err.Error())
		return
	}

	response.Success(c, "Permission deleted successfully", nil)
}

// ListPermissions 获取权限列表
// @Summary 获取权限列表
// @Description 分页获取权限列表
// @Tags 权限管理
// @Accept json
// @Produce json
// @Param page query int false "页码，默认1"
// @Param size query int false "每页数量，默认10"
// @Param type query string false "权限类型，可选值：api, menu, button"
// @Success 200 {object} response.Response{data=dto.PermissionListResponse} "成功"
// @Failure 400 {object} response.Response "请求错误"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /api/v1/permissions [get]
func (h *PermissionHandler) ListPermissions(c *gin.Context) {
	// 解析分页参数
	pageStr := c.DefaultQuery("page", "1")
	sizeStr := c.DefaultQuery("size", "10")
	permType := c.Query("type")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}
	size, err := strconv.Atoi(sizeStr)
	if err != nil || size < 1 || size > 100 {
		size = 10
	}

	var permList *dto.PermissionListResponse
	// 根据类型查询权限列表
	if permType != "" {
		permList, err = h.permissionService.ListPermissionsByType(c.Request.Context(), permType, page, size)
	} else {
		permList, err = h.permissionService.ListPermissions(c.Request.Context(), page, size)
	}

	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "Failed to get permissions: "+err.Error())
		return
	}

	response.Success(c, "Permissions retrieved successfully", permList)
}

// UpdatePermissionStatus 更新权限状态
// @Summary 更新权限状态
// @Description 更新权限的启用/禁用状态
// @Tags 权限管理
// @Accept json
// @Produce json
// @Param id path string true "权限ID"
// @Param request body dto.UpdateStatusRequest true "更新状态请求"
// @Success 200 {object} response.Response "成功"
// @Failure 400 {object} response.Response "请求错误"
// @Failure 404 {object} response.Response "权限不存在"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /api/v1/permissions/{id}/status [patch]
func (h *PermissionHandler) UpdatePermissionStatus(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.Fail(c, http.StatusBadRequest, "Missing permission ID")
		return
	}

	var req dto.UpdateStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, "Invalid request: "+err.Error())
		return
	}

	// 验证请求
	if err := req.Validate(); err != nil {
		response.Fail(c, http.StatusBadRequest, err.Error())
		return
	}

	// 更新权限状态
	if err := h.permissionService.UpdatePermissionStatus(c.Request.Context(), id, req.Status); err != nil {
		switch err {
		case service.ErrInvalidPermissionStatus:
			response.Fail(c, http.StatusBadRequest, "Invalid permission status")
		default:
			response.Fail(c, http.StatusInternalServerError, "Failed to update permission status: "+err.Error())
		}
		return
	}

	response.Success(c, "Permission status updated successfully", nil)
}

// AssignPermissionsToRole 分配权限到角色
// @Summary 分配权限到角色
// @Description 为指定角色分配一组权限
// @Tags 权限管理
// @Accept json
// @Produce json
// @Param id path string true "角色ID"
// @Param request body dto.AssignPermissionsToRoleRequest true "权限分配请求"
// @Success 200 {object} response.Response "成功"
// @Failure 400 {object} response.Response "请求错误"
// @Failure 404 {object} response.Response "角色不存在"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /api/v1/permissions/roles/{id}/assign [post]
func (h *PermissionHandler) AssignPermissionsToRole(c *gin.Context) {
	// 获取角色ID
	roleID := c.Param("id")
	if roleID == "" {
		response.Fail(c, http.StatusBadRequest, "Missing role ID")
		return
	}

	// 解析请求
	var req dto.AssignPermissionsToRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, "Invalid request: "+err.Error())
		return
	}

	// 验证请求
	if err := req.Validate(); err != nil {
		response.Fail(c, http.StatusBadRequest, err.Error())
		return
	}

	// 调用服务层分配权限
	roleService := service.NewRoleService(h.permissionService.GetRoleRepository(), h.permissionService.GetPermissionRepository())
	err := roleService.AssignPermissions(c.Request.Context(), roleID, req.PermissionIDs)
	if err != nil {
		if err.Error() == "record not found" {
			response.Fail(c, http.StatusNotFound, "Role not found")
			return
		}
		if strings.HasPrefix(err.Error(), "invalid permission ID:") {
			response.Fail(c, http.StatusBadRequest, err.Error())
			return
		}
		response.Fail(c, http.StatusInternalServerError, "Failed to assign permissions: "+err.Error())
		return
	}

	response.Success(c, "Permissions assigned successfully", nil)
}