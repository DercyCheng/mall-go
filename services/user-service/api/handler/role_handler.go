package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"mall-go/pkg/response"
	"mall-go/services/user-service/application/service"
	"mall-go/services/user-service/domain/model"
)

// RoleHandler 角色API处理器
type RoleHandler struct {
	roleService *service.RoleService
}

// NewRoleHandler 创建角色API处理器
func NewRoleHandler(roleService *service.RoleService) *RoleHandler {
	return &RoleHandler{
		roleService: roleService,
	}
}

// CreateRole 创建角色
// @Summary 创建角色
// @Description 创建新角色
// @Tags 角色管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param role body dto.RoleDTO true "角色信息"
// @Success 200 {object} response.Response{data=dto.RoleDTO}
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/roles [post]
func (h *RoleHandler) CreateRole(c *gin.Context) {
	var req struct {
		Name        string `json:"name" binding:"required"`
		Description string `json:"description"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "无效的请求参数: "+err.Error())
		return
	}

	role, err := h.roleService.CreateRole(c.Request.Context(), req.Name, req.Description)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "创建角色失败: "+err.Error())
		return
	}

	response.Success(c, role)
}

// GetRole 获取角色详情
// @Summary 获取角色详情
// @Description 根据ID获取角色详情
// @Tags 角色管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path string true "角色ID"
// @Success 200 {object} response.Response{data=dto.RoleDTO}
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 404 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/roles/{id} [get]
func (h *RoleHandler) GetRole(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.BadRequest(c, "角色ID不能为空")
		return
	}

	role, err := h.roleService.GetRole(c.Request.Context(), id)
	if err != nil {
		response.Error(c, http.StatusNotFound, "获取角色失败: "+err.Error())
		return
	}

	response.Success(c, role)
}

// UpdateRole 更新角色
// @Summary 更新角色
// @Description 更新角色信息
// @Tags 角色管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path string true "角色ID"
// @Param role body dto.RoleDTO true "角色信息"
// @Success 200 {object} response.Response{data=dto.RoleDTO}
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 404 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/roles/{id} [put]
func (h *RoleHandler) UpdateRole(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.BadRequest(c, "角色ID不能为空")
		return
	}

	var req struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "无效的请求参数: "+err.Error())
		return
	}

	role, err := h.roleService.UpdateRole(c.Request.Context(), id, req.Name, req.Description)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "更新角色失败: "+err.Error())
		return
	}

	response.Success(c, role)
}

// DeleteRole 删除角色
// @Summary 删除角色
// @Description 删除指定角色
// @Tags 角色管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path string true "角色ID"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/roles/{id} [delete]
func (h *RoleHandler) DeleteRole(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.BadRequest(c, "角色ID不能为空")
		return
	}

	if err := h.roleService.DeleteRole(c.Request.Context(), id); err != nil {
		response.Error(c, http.StatusInternalServerError, "删除角色失败: "+err.Error())
		return
	}

	response.Success(c, nil)
}

// UpdateRoleStatus 更新角色状态
// @Summary 更新角色状态
// @Description 更新角色的状态
// @Tags 角色管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path string true "角色ID"
// @Param status body struct{Status string} true "状态信息"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/roles/{id}/status [patch]
func (h *RoleHandler) UpdateRoleStatus(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.BadRequest(c, "角色ID不能为空")
		return
	}

	var req struct {
		Status string `json:"status" binding:"required,oneof=active inactive"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "无效的请求参数: "+err.Error())
		return
	}

	if err := h.roleService.UpdateRoleStatus(c.Request.Context(), id, model.RoleStatus(req.Status)); err != nil {
		response.Error(c, http.StatusInternalServerError, "更新角色状态失败: "+err.Error())
		return
	}

	response.Success(c, nil)
}

// ListRoles 获取角色列表
// @Summary 获取角色列表
// @Description 分页获取所有角色
// @Tags 角色管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param page query int false "页码" default(1)
// @Param pageSize query int false "每页数量" default(10)
// @Success 200 {object} response.Response{data=dto.RoleListResponse}
// @Failure 401 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/roles [get]
func (h *RoleHandler) ListRoles(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))

	// 校验分页参数
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	roles, err := h.roleService.ListRoles(c.Request.Context(), page, pageSize)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "获取角色列表失败: "+err.Error())
		return
	}

	response.Success(c, roles)
}

// AssignPermissions 分配权限
// @Summary 分配权限
// @Description 为角色分配权限
// @Tags 角色管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path string true "角色ID"
// @Param permissions body struct{PermissionIDs []string} true "权限信息"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/roles/{id}/permissions [post]
func (h *RoleHandler) AssignPermissions(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.BadRequest(c, "角色ID不能为空")
		return
	}

	var req struct {
		PermissionIDs []string `json:"permissionIds" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "无效的请求参数: "+err.Error())
		return
	}

	if err := h.roleService.AssignPermissions(c.Request.Context(), id, req.PermissionIDs); err != nil {
		response.Error(c, http.StatusInternalServerError, "分配权限失败: "+err.Error())
		return
	}

	response.Success(c, nil)
}

// GetRolePermissions 获取角色权限
// @Summary 获取角色权限
// @Description 获取角色的所有权限
// @Tags 角色管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path string true "角色ID"
// @Success 200 {object} response.Response{data=[]dto.PermissionDTO}
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/roles/{id}/permissions [get]
func (h *RoleHandler) GetRolePermissions(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.BadRequest(c, "角色ID不能为空")
		return
	}

	permissions, err := h.roleService.GetRolePermissions(c.Request.Context(), id)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "获取角色权限失败: "+err.Error())
		return
	}

	response.Success(c, permissions)
}
