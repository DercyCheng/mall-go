package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"mall-go/pkg/errors"
	"mall-go/pkg/response"
	"mall-go/pkg/validation"
	"mall-go/services/user-service/application/dto"
	"mall-go/services/user-service/application/service"
)

// UserHandler handles HTTP requests related to users
type UserHandler struct {
	userService service.UserService
	roleService service.RoleService
}

// NewUserHandler creates a new UserHandler
func NewUserHandler(userService service.UserService, roleService service.RoleService) *UserHandler {
	return &UserHandler{
		userService: userService,
		roleService: roleService,
	}
}

// Register handles user registration
func (h *UserHandler) Register(c *gin.Context) {
	var req dto.UserCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "无效的请求参数: "+err.Error())
		return
	}

	// 执行额外的自定义验证
	if err := validation.ValidateStruct(req); err != nil {
		response.Error(c, err)
		return
	}

	// 验证用户名是否包含特殊字符
	if !validation.ValidateUsername(req.Username) {
		response.BadRequest(c, "用户名只能包含字母、数字、下划线、短横线和点，长度为3-20个字符")
		return
	}

	// 清理输入
	req.NickName = validation.SanitizeString(req.NickName)
	req.Note = validation.SanitizeString(req.Note)

	id, err := h.userService.Register(c.Request.Context(), req)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, gin.H{"id": id})
}

// Login handles user login
func (h *UserHandler) Login(c *gin.Context) {
	var req dto.UserLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "无效的登录信息: "+err.Error())
		return
	}

	// 防止XSS攻击
	req.Username = validation.SanitizeString(req.Username)

	resp, err := h.userService.Login(c.Request.Context(), req)
	if err != nil {
		// 检查是否为凭证错误
		if _, ok := errors.As(err); ok {
			response.Error(c, err)
		} else {
			response.Unauthorized(c, "用户名或密码错误")
		}
		return
	}

	response.Success(c, resp)
}

// GetUserInfo gets the current user's information
func (h *UserHandler) GetUserInfo(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		response.Unauthorized(c, "用户未认证")
		return
	}

	userInfo, err := h.userService.GetUserInfo(c.Request.Context(), userID.(string))
	if err != nil {
		if appErr, ok := errors.As(err); ok && appErr.Type == errors.TypeNotFound {
			response.NotFound(c, "用户不存在")
		} else {
			response.Error(c, err)
		}
		return
	}

	response.Success(c, userInfo)
}

// GetUserByID gets a user by ID
func (h *UserHandler) GetUserByID(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.BadRequest(c, "用户ID不能为空")
		return
	}

	userInfo, err := h.userService.GetUserInfo(c.Request.Context(), id)
	if err != nil {
		if appErr, ok := errors.As(err); ok && appErr.Type == errors.TypeNotFound {
			response.NotFound(c, "用户不存在")
		} else {
			response.Error(c, err)
		}
		return
	}

	response.Success(c, userInfo)
}

// ListUsers lists all users with pagination
func (h *UserHandler) ListUsers(c *gin.Context) {
	pageStr := c.DefaultQuery("page", "1")
	pageSizeStr := c.DefaultQuery("pageSize", "10")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil || pageSize < 1 || pageSize > 100 {
		// 限制最大页大小为100，防止资源耗尽攻击
		pageSize = 10
	}

	result, err := h.userService.ListUsers(c.Request.Context(), page, pageSize)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, result)
}

// SearchUsers searches for users
func (h *UserHandler) SearchUsers(c *gin.Context) {
	query := validation.SanitizeString(c.Query("keyword"))
	pageStr := c.DefaultQuery("page", "1")
	pageSizeStr := c.DefaultQuery("pageSize", "10")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil || pageSize < 1 || pageSize > 100 {
		// 限制最大页大小为100，防止资源耗尽攻击
		pageSize = 10
	}

	result, err := h.userService.SearchUsers(c.Request.Context(), query, page, pageSize)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, result)
}

// UpdateUser updates a user
func (h *UserHandler) UpdateUser(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.BadRequest(c, "用户ID不能为空")
		return
	}

	var req dto.UserUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "无效的请求参数: "+err.Error())
		return
	}

	// 执行额外的验证
	if err := validation.ValidateStruct(req); err != nil {
		response.Error(c, err)
		return
	}

	// 清理输入
	if req.NickName != "" {
		req.NickName = validation.SanitizeString(req.NickName)
	}
	if req.Note != "" {
		req.Note = validation.SanitizeString(req.Note)
	}

	if err := h.userService.UpdateUser(c.Request.Context(), id, req); err != nil {
		if appErr, ok := errors.As(err); ok && appErr.Type == errors.TypeNotFound {
			response.NotFound(c, "用户不存在")
		} else {
			response.Error(c, err)
		}
		return
	}

	response.Success(c, nil)
}

// DeleteUser deletes a user
func (h *UserHandler) DeleteUser(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.BadRequest(c, "用户ID不能为空")
		return
	}

	if err := h.userService.DeleteUser(c.Request.Context(), id); err != nil {
		if appErr, ok := errors.As(err); ok && appErr.Type == errors.TypeNotFound {
			response.NotFound(c, "用户不存在")
		} else {
			response.Error(c, err)
		}
		return
	}

	response.Success(c, nil)
}

// ChangePassword changes a user's password
func (h *UserHandler) ChangePassword(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		response.Unauthorized(c, "用户未认证")
		return
	}

	var req dto.UserChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "无效的请求参数: "+err.Error())
		return
	}

	// 验证新密码强度
	if !validation.ValidatePassword(req.NewPassword) {
		response.BadRequest(c, "新密码必须包含数字、大小写字母和特殊字符，长度至少为8位")
		return
	}

	if err := h.userService.ChangePassword(c.Request.Context(), userID.(string), req); err != nil {
		if appErr, ok := errors.As(err); ok {
			if appErr.Code == errors.CodePasswordMismatch {
				response.BadRequest(c, "当前密码不正确")
			} else {
				response.Error(c, err)
			}
		} else {
			response.Error(c, err)
		}
		return
	}

	response.Success(c, nil)
}

// ListRoles lists all roles
func (h *UserHandler) ListRoles(c *gin.Context) {
	roles, err := h.roleService.ListRoles(c.Request.Context())
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, roles)
}

// CreateRole creates a new role
func (h *UserHandler) CreateRole(c *gin.Context) {
	var req dto.RoleCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "无效的请求参数: "+err.Error())
		return
	}

	// 执行额外的验证
	if err := validation.ValidateStruct(req); err != nil {
		response.Error(c, err)
		return
	}

	// 清理输入
	req.Name = validation.SanitizeString(req.Name)
	req.Description = validation.SanitizeString(req.Description)

	id, err := h.roleService.CreateRole(c.Request.Context(), req)
	if err != nil {
		if appErr, ok := errors.As(err); ok && appErr.Code == errors.CodeRoleAlreadyExists {
			response.ErrorWithCode(c, 409, "角色名称已存在")
		} else {
			response.Error(c, err)
		}
		return
	}

	response.Success(c, gin.H{"id": id})
}

// GetRole gets a role by ID
func (h *UserHandler) GetRole(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.BadRequest(c, "角色ID不能为空")
		return
	}

	role, err := h.roleService.GetRole(c.Request.Context(), id)
	if err != nil {
		if appErr, ok := errors.As(err); ok && appErr.Type == errors.TypeNotFound {
			response.NotFound(c, "角色不存在")
		} else {
			response.Error(c, err)
		}
		return
	}

	response.Success(c, role)
}

// UpdateRole updates a role
func (h *UserHandler) UpdateRole(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.BadRequest(c, "角色ID不能为空")
		return
	}

	var req dto.RoleUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "无效的请求参数: "+err.Error())
		return
	}

	// 执行额外的验证
	if err := validation.ValidateStruct(req); err != nil {
		response.Error(c, err)
		return
	}

	// 清理输入
	if req.Name != "" {
		req.Name = validation.SanitizeString(req.Name)
	}
	if req.Description != "" {
		req.Description = validation.SanitizeString(req.Description)
	}

	if err := h.roleService.UpdateRole(c.Request.Context(), id, req); err != nil {
		if appErr, ok := errors.As(err); ok {
			if appErr.Type == errors.TypeNotFound {
				response.NotFound(c, "角色不存在")
			} else if appErr.Code == errors.CodeRoleAlreadyExists {
				response.ErrorWithCode(c, 409, "角色名称已存在")
			} else {
				response.Error(c, err)
			}
		} else {
			response.Error(c, err)
		}
		return
	}

	response.Success(c, nil)
}

// DeleteRole deletes a role
func (h *UserHandler) DeleteRole(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.BadRequest(c, "角色ID不能为空")
		return
	}

	if err := h.roleService.DeleteRole(c.Request.Context(), id); err != nil {
		if appErr, ok := errors.As(err); ok && appErr.Type == errors.TypeNotFound {
			response.NotFound(c, "角色不存在")
		} else {
			response.Error(c, err)
		}
		return
	}

	response.Success(c, nil)
}

// AssignRolesToUser assigns roles to a user
func (h *UserHandler) AssignRolesToUser(c *gin.Context) {
	var req dto.AssignRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "无效的请求参数: "+err.Error())
		return
	}

	// 执行额外的验证
	if err := validation.ValidateStruct(req); err != nil {
		response.Error(c, err)
		return
	}

	if err := h.roleService.AssignRolesToUser(c.Request.Context(), req); err != nil {
		if appErr, ok := errors.As(err); ok {
			if appErr.Type == errors.TypeNotFound {
				if appErr.Code == errors.CodeUserNotFound {
					response.NotFound(c, "用户不存在")
				} else if appErr.Code == errors.CodeRoleNotFound {
					response.NotFound(c, "角色不存在")
				} else {
					response.Error(c, err)
				}
			} else {
				response.Error(c, err)
			}
		} else {
			response.Error(c, err)
		}
		return
	}

	response.Success(c, nil)
}

// GetUserRoles gets roles for a user
func (h *UserHandler) GetUserRoles(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.BadRequest(c, "用户ID不能为空")
		return
	}

	roles, err := h.roleService.GetUserRoles(c.Request.Context(), id)
	if err != nil {
		if appErr, ok := errors.As(err); ok && appErr.Type == errors.TypeNotFound {
			response.NotFound(c, "用户不存在")
		} else {
			response.Error(c, err)
		}
		return
	}

	response.Success(c, roles)
}
