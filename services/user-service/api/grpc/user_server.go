package grpc

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"mall-go/services/user-service/application/dto"
	"mall-go/services/user-service/application/service"
	userpb "mall-go/services/user-service/proto"
)

// UserServer 实现UserService gRPC接口
type UserServer struct {
	userpb.UnimplementedUserServiceServer
	userService service.UserService
}

// NewUserServer 创建新的UserServer
func NewUserServer(userService service.UserService) *UserServer {
	return &UserServer{
		userService: userService,
	}
}

// Register 注册新用户
func (s *UserServer) Register(ctx context.Context, req *userpb.RegisterRequest) (*userpb.RegisterResponse, error) {
	// 转换请求
	dtoReq := dto.UserCreateRequest{
		Username: req.Username,
		Password: req.Password,
		Email:    req.Email,
		Phone:    req.Phone,
		NickName: req.Username, // Using username as default nickname
	}

	// 调用应用服务
	userId, err := s.userService.Register(ctx, dtoReq)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "用户注册失败: %v", err)
	}

	return &userpb.RegisterResponse{
		Success: true,
		Message: "用户注册成功",
		UserId:  userId,
	}, nil
}

// Login 用户登录
func (s *UserServer) Login(ctx context.Context, req *userpb.LoginRequest) (*userpb.LoginResponse, error) {
	// 转换请求
	dtoReq := dto.UserLoginRequest{
		Username: req.Username,
		Password: req.Password,
	}

	// 调用应用服务
	result, err := s.userService.Login(ctx, dtoReq)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "用户登录失败: %v", err)
	}

	// 构建用户响应
	user := &userpb.User{
		Id:        result.User.ID,
		Username:  result.User.Username,
		Email:     result.User.Email,
		Phone:     result.User.Phone,
		Avatar:    result.User.Icon, // Map Icon to Avatar
		Status:    int32(result.User.Status),
		CreatedAt: result.User.CreatedAt, // CreatedAt is already a string
	}

	return &userpb.LoginResponse{
		Success: true,
		Message: "用户登录成功",
		Token:   result.Token,
		User:    user,
	}, nil
}

// GetUserInfo 获取用户信息
func (s *UserServer) GetUserInfo(ctx context.Context, req *userpb.GetUserInfoRequest) (*userpb.GetUserInfoResponse, error) {
	// 调用应用服务
	result, err := s.userService.GetUserInfo(ctx, req.UserId)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "获取用户信息失败: %v", err)
	}

	// 构建用户响应
	user := &userpb.User{
		Id:        result.ID,
		Username:  result.Username,
		Email:     result.Email,
		Phone:     result.Phone,
		Avatar:    result.Icon, // Map Icon to Avatar
		Status:    int32(result.Status),
		CreatedAt: result.CreatedAt, // CreatedAt is already a string
	}

	return &userpb.GetUserInfoResponse{
		Success: true,
		Message: "获取用户信息成功",
		User:    user,
	}, nil
}

// UpdateUserInfo 更新用户信息
func (s *UserServer) UpdateUserInfo(ctx context.Context, req *userpb.UpdateUserInfoRequest) (*userpb.UpdateUserInfoResponse, error) {
	// 转换请求
	dtoReq := dto.UserUpdateRequest{
		Email:    req.Email,
		Phone:    req.Phone,
		Icon:     req.Avatar,   // Map Avatar to Icon
		NickName: req.Username, // Map Username to NickName
	}

	// 调用应用服务
	err := s.userService.UpdateUser(ctx, req.UserId, dtoReq)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "更新用户信息失败: %v", err)
	}

	return &userpb.UpdateUserInfoResponse{
		Success: true,
		Message: "更新用户信息成功",
	}, nil
}
