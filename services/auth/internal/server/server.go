package server

import (
	"context"

	"SituationBak/services/auth/internal/service"
	"SituationBak/shared/errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// AuthServer gRPC 认证服务实现
type AuthServer struct {
	UnimplementedAuthServiceServer
	authService *service.AuthService
}

// NewAuthServer 创建认证服务gRPC服务器
func NewAuthServer(authService *service.AuthService) *AuthServer {
	return &AuthServer{
		authService: authService,
	}
}

// Register 用户注册
func (s *AuthServer) Register(ctx context.Context, req *RegisterRequest) (*LoginResponse, error) {
	result, err := s.authService.Register(&service.RegisterRequest{
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		return nil, toGRPCError(err)
	}

	return &LoginResponse{
		AccessToken:  result.AccessToken,
		RefreshToken: result.RefreshToken,
		TokenType:    result.TokenType,
		ExpiresIn:    result.ExpiresIn,
		User: &UserInfo{
			Id:        uint64(result.User.ID),
			Username:  result.User.Username,
			Email:     result.User.Email,
			Role:      result.User.Role,
			AvatarUrl: result.User.AvatarURL,
			CreatedAt: result.User.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		},
	}, nil
}

// Login 用户登录
func (s *AuthServer) Login(ctx context.Context, req *LoginRequest) (*LoginResponse, error) {
	result, err := s.authService.Login(&service.LoginRequest{
		Username: req.Username,
		Password: req.Password,
	})
	if err != nil {
		return nil, toGRPCError(err)
	}

	return &LoginResponse{
		AccessToken:  result.AccessToken,
		RefreshToken: result.RefreshToken,
		TokenType:    result.TokenType,
		ExpiresIn:    result.ExpiresIn,
		User: &UserInfo{
			Id:        uint64(result.User.ID),
			Username:  result.User.Username,
			Email:     result.User.Email,
			Role:      result.User.Role,
			AvatarUrl: result.User.AvatarURL,
			CreatedAt: result.User.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		},
	}, nil
}

// RefreshToken 刷新Token
func (s *AuthServer) RefreshToken(ctx context.Context, req *RefreshTokenRequest) (*TokenResponse, error) {
	result, err := s.authService.RefreshToken(req.RefreshToken)
	if err != nil {
		return nil, toGRPCError(err)
	}

	return &TokenResponse{
		AccessToken:  result.AccessToken,
		RefreshToken: result.RefreshToken,
		TokenType:    result.TokenType,
		ExpiresIn:    result.ExpiresIn,
	}, nil
}

// ValidateToken 验证Token
func (s *AuthServer) ValidateToken(ctx context.Context, req *ValidateTokenRequest) (*ValidateTokenResponse, error) {
	result, err := s.authService.ValidateToken(req.AccessToken)
	if err != nil {
		return nil, toGRPCError(err)
	}

	return &ValidateTokenResponse{
		Valid:    result.Valid,
		UserId:   uint64(result.UserID),
		Username: result.Username,
		Role:     result.Role,
	}, nil
}

// GetCurrentUser 获取当前用户信息
func (s *AuthServer) GetCurrentUser(ctx context.Context, req *GetCurrentUserRequest) (*UserInfo, error) {
	result, err := s.authService.GetCurrentUser(uint(req.UserId))
	if err != nil {
		return nil, toGRPCError(err)
	}

	return &UserInfo{
		Id:        uint64(result.ID),
		Username:  result.Username,
		Email:     result.Email,
		Role:      result.Role,
		AvatarUrl: result.AvatarURL,
		CreatedAt: result.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}, nil
}

// toGRPCError 将应用错误转换为gRPC错误
func toGRPCError(err error) error {
	if appErr, ok := errors.AsAppError(err); ok {
		code := mapToGRPCCode(appErr.Code)
		return status.Error(code, appErr.Message)
	}
	return status.Error(codes.Internal, err.Error())
}

// mapToGRPCCode 将业务错误码映射为gRPC状态码
func mapToGRPCCode(code int) codes.Code {
	switch code {
	case errors.CodeInvalidParams:
		return codes.InvalidArgument
	case errors.CodeNotFound, errors.CodeUserNotFound:
		return codes.NotFound
	case errors.CodeAlreadyExists, errors.CodeEmailExists, errors.CodeUsernameExist:
		return codes.AlreadyExists
	case errors.CodeUnauthorized, errors.CodeTokenExpired, errors.CodeTokenInvalid, errors.CodeLoginFailed, errors.CodePasswordWrong:
		return codes.Unauthenticated
	case errors.CodeForbidden:
		return codes.PermissionDenied
	case errors.CodeInternalError:
		return codes.Internal
	case errors.CodeServiceUnavailable:
		return codes.Unavailable
	default:
		return codes.Unknown
	}
}
