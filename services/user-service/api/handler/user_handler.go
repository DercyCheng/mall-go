package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"mall-go/pkg/response"
	"mall-go/services/user-service/application/dto"
	"mall-go/services/user-service/application/service"
)

// UserHandler 用户API处理器
type UserHandler struct {
	userService *service.UserService
}

// NewUserHandler 创建用户API处理器
func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

// Register 用户注册
// @Summary 用户注册
// @Description 注册新用户
// @Tags 用户管理
// @Accept json
// @Produce json
// @Param user body dto.RegisterUserRequest true "用户注册信息"
// @Success 200 {object} response.Response{data=dto.UserResponse}
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/users/register [post]
func (h *UserHandler) Register(c *gin.Context) {
	var req dto.RegisterUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "无效的请求参数: "+err.Error())
		return
	}

	user, err := h.userService.Register(c.Request.Context(), req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "注册失败: "+err.Error())
		return
	}

	response.Success(c, user)
}

// Login 用户登录
// @Summary 用户登录
// @Description 用户登录获取token
// @Tags 用户管理
// @Accept json
// @Produce json
// @Param login body dto.LoginRequest true "登录信息"
// @Success 200 {object} response.Response{data=dto.LoginResponse}
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/users/login [post]
func (h *UserHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "无效的请求参数: "+err.Error())
		return
	}

	loginResp, err := h.userService.Login(c.Request.Context(), req)
	if err != nil {
		response.Error(c, http.StatusUnauthorized, "登录失败: "+err.Error())
		return
	}

	response.Success(c, loginResp)
}

// GetUserInfo 获取用户信息
// @Summary 获取用户信息
// @Description 获取当前登录用户信息
// @Tags 用户管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} response.Response{data=dto.UserResponse}
// @Failure 401 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/users/info [get]
func (h *UserHandler) GetUserInfo(c *gin.Context) {
	// 从JWT token中获取用户ID
	userID, exists := c.Get("userId")
	if !exists {
		response.Unauthorized(c)
		return
	}

	user, err := h.userService.GetUserInfo(c.Request.Context(), userID.(string))
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "获取用户信息失败: "+err.Error())
		return
	}

	response.Success(c, user)
}

// GetUser 获取指定用户信息
// @Summary 获取指定用户信息
// @Description 根据用户ID获取用户信息
// @Tags 用户管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path string true "用户ID"
// @Success 200 {object} response.Response{data=dto.UserResponse}
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 404 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/users/{id} [get]
func (h *UserHandler) GetUser(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.BadRequest(c, "用户ID不能为空")
		return
	}

	user, err := h.userService.GetUserInfo(c.Request.Context(), id)
	if err != nil {
		response.Error(c, http.StatusNotFound, "获取用户信息失败: "+err.Error())
		return
	}

	response.Success(c, user)
}

// UpdateUserInfo 更新用户信息
// @Summary 更新用户信息
// @Description 更新当前登录用户的信息
// @Tags 用户管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param user body dto.UpdateUserRequest true "用户信息"
// @Success 200 {object} response.Response{data=dto.UserResponse}
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/users/info [put]
func (h *UserHandler) UpdateUserInfo(c *gin.Context) {
	// 从JWT token中获取用户ID
	userID, exists := c.Get("userId")
	if !exists {
		response.Unauthorized(c)
		return
	}

	var req dto.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "无效的请求参数: "+err.Error())
		return
	}

	user, err := h.userService.UpdateUserInfo(c.Request.Context(), userID.(string), req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "更新用户信息失败: "+err.Error())
		return
	}

	response.Success(c, user)
}

// UpdatePassword 修改密码
// @Summary 修改密码
// @Description 修改当前登录用户的密码
// @Tags 用户管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param password body dto.UpdatePasswordRequest true "密码信息"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/users/password [put]
func (h *UserHandler) UpdatePassword(c *gin.Context) {
	// 从JWT token中获取用户ID
	userID, exists := c.Get("userId")
	if !exists {
		response.Unauthorized(c)
		return
	}

	var req dto.UpdatePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "无效的请求参数: "+err.Error())
		return
	}

	if err := h.userService.UpdatePassword(c.Request.Context(), userID.(string), req); err != nil {
		response.Error(c, http.StatusInternalServerError, "修改密码失败: "+err.Error())
		return
	}

	response.Success(c, nil)
}

// UpdateStatus 更新用户状态
// @Summary 更新用户状态
// @Description 更新指定用户的状态
// @Tags 用户管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path string true "用户ID"
// @Param status body dto.UpdateStatusRequest true "状态信息"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/users/{id}/status [patch]
func (h *UserHandler) UpdateStatus(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.BadRequest(c, "用户ID不能为空")
		return
	}

	var req dto.UpdateStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "无效的请求参数: "+err.Error())
		return
	}

	if err := h.userService.UpdateStatus(c.Request.Context(), id, req); err != nil {
		response.Error(c, http.StatusInternalServerError, "更新用户状态失败: "+err.Error())
		return
	}

	response.Success(c, nil)
}

// AssignRoles 分配角色
// @Summary 分配角色
// @Description 为指定用户分配角色
// @Tags 用户管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path string true "用户ID"
// @Param roles body dto.AssignRoleRequest true "角色信息"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/users/{id}/roles [post]
func (h *UserHandler) AssignRoles(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.BadRequest(c, "用户ID不能为空")
		return
	}

	var req dto.AssignRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "无效的请求参数: "+err.Error())
		return
	}

	if err := h.userService.AssignRoles(c.Request.Context(), id, req); err != nil {
		response.Error(c, http.StatusInternalServerError, "分配角色失败: "+err.Error())
		return
	}

	response.Success(c, nil)
}

// ListUsers 获取用户列表
// @Summary 获取用户列表
// @Description 分页获取所有用户
// @Tags 用户管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param page query int false "页码" default(1)
// @Param pageSize query int false "每页数量" default(10)
// @Success 200 {object} response.Response{data=dto.UserListResponse}
// @Failure 401 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/users [get]
func (h *UserHandler) ListUsers(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))

	// 校验分页参数
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	users, err := h.userService.GetUsers(c.Request.Context(), page, pageSize)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "获取用户列表失败: "+err.Error())
		return
	}

	response.Success(c, users)
}

// SearchUsers 搜索用户
// @Summary 搜索用户
// @Description 根据关键词搜索用户
// @Tags 用户管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param keyword query string false "搜索关键词"
// @Param page query int false "页码" default(1)
// @Param pageSize query int false "每页数量" default(10)
// @Success 200 {object} response.Response{data=dto.UserListResponse}
// @Failure 401 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/users/search [get]
func (h *UserHandler) SearchUsers(c *gin.Context) {
	var req dto.SearchUserRequest
	req.Keyword = c.Query("keyword")
	req.Page, _ = strconv.Atoi(c.DefaultQuery("page", "1"))
	req.PageSize, _ = strconv.Atoi(c.DefaultQuery("pageSize", "10"))

	// 校验分页参数
	if req.Page < 1 {
		req.Page = 1
	}
	if req.PageSize < 1 || req.PageSize > 100 {
		req.PageSize = 10
	}

	users, err := h.userService.SearchUsers(c.Request.Context(), req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "搜索用户失败: "+err.Error())
		return
	}

	response.Success(c, users)
}

// DeleteUser 删除用户
// @Summary 删除用户
// @Description 删除指定用户
// @Tags 用户管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path string true "用户ID"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/users/{id} [delete]
func (h *UserHandler) DeleteUser(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.BadRequest(c, "用户ID不能为空")
		return
	}

	if err := h.userService.DeleteUser(c.Request.Context(), id); err != nil {
		response.Error(c, http.StatusInternalServerError, "删除用户失败: "+err.Error())
		return
	}

	response.Success(c, nil)
}