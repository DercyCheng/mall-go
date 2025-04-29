package router

import (
	"github.com/gin-gonic/gin"

	"mall-go/services/user-service/api/handler"
	"mall-go/services/user-service/api/middleware"
)

// Router defines the HTTP router interface
type Router interface {
	// Register registers all routes
	Register(engine *gin.Engine)
}

// UserRouter implements Router for user service
type UserRouter struct {
	userHandler *handler.UserHandler
	jwtSecret   string
}

// NewUserRouter creates a new UserRouter
func NewUserRouter(userHandler *handler.UserHandler, jwtSecret string) Router {
	return &UserRouter{
		userHandler: userHandler,
		jwtSecret:   jwtSecret,
	}
}

// Register registers all routes
func (r *UserRouter) Register(engine *gin.Engine) {
	// Apply global middlewares
	engine.Use(middleware.Cors())

	// Health check endpoint
	// swagger:operation GET /health system healthCheck
	// ---
	// summary: 健康检查
	// description: 用于检查服务是否正常运行的API
	// responses:
	//   "200":
	//     description: 服务正常运行
	//     schema:
	//       type: object
	//       properties:
	//         status:
	//           type: string
	//           example: "UP"
	engine.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "UP"})
	})

	// Public routes (no authentication required)
	publicGroup := engine.Group("/api/v1")
	{
		// swagger:operation POST /api/v1/register users registerUser
		// ---
		// summary: 用户注册
		// description: 创建新用户账号
		// parameters:
		// - name: body
		//   in: body
		//   description: 用户注册信息
		//   required: true
		//   schema:
		//     "$ref": "#/definitions/UserCreateRequest"
		// responses:
		//   "200":
		//     "$ref": "#/responses/registerUserResponse"
		//   "400":
		//     "$ref": "#/responses/badRequestResponse"
		//   "500":
		//     "$ref": "#/responses/internalServerErrorResponse"
		publicGroup.POST("/register", r.userHandler.Register)

		// swagger:operation POST /api/v1/login users loginUser
		// ---
		// summary: 用户登录
		// description: 使用用户名和密码登录系统
		// parameters:
		// - name: body
		//   in: body
		//   description: 登录凭证
		//   required: true
		//   schema:
		//     "$ref": "#/definitions/UserLoginRequest"
		// responses:
		//   "200":
		//     "$ref": "#/responses/loginUserResponse"
		//   "400":
		//     "$ref": "#/responses/badRequestResponse"
		//   "401":
		//     "$ref": "#/responses/unauthorizedResponse"
		//   "500":
		//     "$ref": "#/responses/internalServerErrorResponse"
		publicGroup.POST("/login", r.userHandler.Login)
	}

	// Protected routes (authentication required)
	authGroup := engine.Group("/api/v1")
	authGroup.Use(middleware.JWTAuthMiddleware(r.jwtSecret))
	{
		// User endpoints
		// swagger:operation GET /api/v1/users/me users getUserInfo
		// ---
		// summary: 获取当前用户信息
		// description: 获取当前已登录用户的详细信息
		// security:
		// - Bearer: []
		// responses:
		//   "200":
		//     "$ref": "#/responses/getUserInfoResponse"
		//   "401":
		//     "$ref": "#/responses/unauthorizedResponse"
		//   "500":
		//     "$ref": "#/responses/internalServerErrorResponse"
		authGroup.GET("/users/me", r.userHandler.GetUserInfo)

		// swagger:operation GET /api/v1/users users listUsers
		// ---
		// summary: 获取用户列表
		// description: 分页获取所有用户的列表
		// security:
		// - Bearer: []
		// parameters:
		// - name: page
		//   in: query
		//   description: 页码（从1开始）
		//   type: integer
		//   default: 1
		// - name: pageSize
		//   in: query
		//   description: 每页记录数
		//   type: integer
		//   default: 10
		// responses:
		//   "200":
		//     description: 成功获取用户列表
		//   "401":
		//     "$ref": "#/responses/unauthorizedResponse"
		//   "500":
		//     "$ref": "#/responses/internalServerErrorResponse"
		authGroup.GET("/users", r.userHandler.ListUsers)

		// swagger:operation GET /api/v1/users/search users searchUsers
		// ---
		// summary: 搜索用户
		// description: 根据关键词搜索用户
		// security:
		// - Bearer: []
		// parameters:
		// - name: keyword
		//   in: query
		//   description: 搜索关键词
		//   type: string
		//   required: true
		// - name: page
		//   in: query
		//   description: 页码（从1开始）
		//   type: integer
		//   default: 1
		// - name: pageSize
		//   in: query
		//   description: 每页记录数
		//   type: integer
		//   default: 10
		// responses:
		//   "200":
		//     description: 成功获取搜索结果
		//   "401":
		//     "$ref": "#/responses/unauthorizedResponse"
		//   "500":
		//     "$ref": "#/responses/internalServerErrorResponse"
		authGroup.GET("/users/search", r.userHandler.SearchUsers)

		// swagger:operation GET /api/v1/users/{id} users getUserByID
		// ---
		// summary: 获取指定用户信息
		// description: 根据用户ID获取用户详细信息
		// security:
		// - Bearer: []
		// parameters:
		// - name: id
		//   in: path
		//   description: 用户ID
		//   type: string
		//   required: true
		// responses:
		//   "200":
		//     "$ref": "#/responses/getUserInfoResponse"
		//   "400":
		//     "$ref": "#/responses/badRequestResponse"
		//   "401":
		//     "$ref": "#/responses/unauthorizedResponse"
		//   "404":
		//     description: 用户不存在
		//   "500":
		//     "$ref": "#/responses/internalServerErrorResponse"
		authGroup.GET("/users/:id", r.userHandler.GetUserByID)

		// swagger:operation PUT /api/v1/users/{id} users updateUser
		// ---
		// summary: 更新用户信息
		// description: 更新指定用户的信息
		// security:
		// - Bearer: []
		// parameters:
		// - name: id
		//   in: path
		//   description: 用户ID
		//   type: string
		//   required: true
		// - name: body
		//   in: body
		//   description: 用户更新信息
		//   required: true
		//   schema:
		//     "$ref": "#/definitions/UserUpdateRequest"
		// responses:
		//   "200":
		//     "$ref": "#/responses/updateUserResponse"
		//   "400":
		//     "$ref": "#/responses/badRequestResponse"
		//   "401":
		//     "$ref": "#/responses/unauthorizedResponse"
		//   "404":
		//     description: 用户不存在
		//   "500":
		//     "$ref": "#/responses/internalServerErrorResponse"
		authGroup.PUT("/users/:id", r.userHandler.UpdateUser)

		// swagger:operation DELETE /api/v1/users/{id} users deleteUser
		// ---
		// summary: 删除用户
		// description: 删除指定的用户
		// security:
		// - Bearer: []
		// parameters:
		// - name: id
		//   in: path
		//   description: 用户ID
		//   type: string
		//   required: true
		// responses:
		//   "200":
		//     "$ref": "#/responses/deleteUserResponse"
		//   "400":
		//     "$ref": "#/responses/badRequestResponse"
		//   "401":
		//     "$ref": "#/responses/unauthorizedResponse"
		//   "404":
		//     description: 用户不存在
		//   "500":
		//     "$ref": "#/responses/internalServerErrorResponse"
		authGroup.DELETE("/users/:id", r.userHandler.DeleteUser)

		// 其他API端点的Swagger文档
		// ...省略其余路由的Swagger注释，保持相同的模式...
		authGroup.POST("/users/password", r.userHandler.ChangePassword)
		authGroup.GET("/users/:id/roles", r.userHandler.GetUserRoles)

		// Role endpoints
		authGroup.GET("/roles", r.userHandler.ListRoles)
		authGroup.POST("/roles", r.userHandler.CreateRole)
		authGroup.GET("/roles/:id", r.userHandler.GetRole)
		authGroup.PUT("/roles/:id", r.userHandler.UpdateRole)
		authGroup.DELETE("/roles/:id", r.userHandler.DeleteRole)
		authGroup.POST("/roles/assign", r.userHandler.AssignRolesToUser)
	}
}
