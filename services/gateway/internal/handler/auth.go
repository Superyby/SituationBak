package handler

import (
	"context"

	"SituationBak/services/gateway/internal/client"
	"SituationBak/services/gateway/internal/middleware"
	"SituationBak/shared/errors"

	"github.com/gofiber/fiber/v3"
)

// AuthHandler 认证处理器
type AuthHandler struct {
	authClient *client.AuthClient
}

// NewAuthHandler 创建认证处理器
func NewAuthHandler(authClient *client.AuthClient) *AuthHandler {
	return &AuthHandler{authClient: authClient}
}

// RegisterRequest 注册请求
type RegisterRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// LoginRequest 登录请求
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// RefreshTokenRequest 刷新Token请求
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token"`
}

// Register 用户注册
func (h *AuthHandler) Register(c fiber.Ctx) error {
	var req RegisterRequest
	if err := c.Bind().Body(&req); err != nil {
		return Fail(c, errors.CodeInvalidParams, "请求参数格式错误")
	}

	// 参数验证
	if req.Username == "" || req.Email == "" || req.Password == "" {
		return Fail(c, errors.CodeInvalidParams, "用户名、邮箱和密码不能为空")
	}

	if len(req.Username) < 3 || len(req.Username) > 50 {
		return Fail(c, errors.CodeInvalidParams, "用户名长度应在3-50个字符之间")
	}

	if len(req.Password) < 6 {
		return Fail(c, errors.CodeInvalidParams, "密码长度至少6个字符")
	}

	// 调用认证服务
	result, err := h.authClient.Register(context.Background(), req.Username, req.Email, req.Password)
	if err != nil {
		return FailWithGRPCError(c, err)
	}

	return Created(c, map[string]interface{}{
		"access_token":  result.AccessToken,
		"refresh_token": result.RefreshToken,
		"token_type":    result.TokenType,
		"expires_in":    result.ExpiresIn,
		"user": map[string]interface{}{
			"id":         result.User.Id,
			"username":   result.User.Username,
			"email":      result.User.Email,
			"role":       result.User.Role,
			"avatar_url": result.User.AvatarUrl,
			"created_at": result.User.CreatedAt,
		},
	})
}

// Login 用户登录
func (h *AuthHandler) Login(c fiber.Ctx) error {
	var req LoginRequest
	if err := c.Bind().Body(&req); err != nil {
		return Fail(c, errors.CodeInvalidParams, "请求参数格式错误")
	}

	if req.Username == "" || req.Password == "" {
		return Fail(c, errors.CodeInvalidParams, "用户名和密码不能为空")
	}

	// 调用认证服务
	result, err := h.authClient.Login(context.Background(), req.Username, req.Password)
	if err != nil {
		return FailWithGRPCError(c, err)
	}

	return Success(c, map[string]interface{}{
		"access_token":  result.AccessToken,
		"refresh_token": result.RefreshToken,
		"token_type":    result.TokenType,
		"expires_in":    result.ExpiresIn,
		"user": map[string]interface{}{
			"id":         result.User.Id,
			"username":   result.User.Username,
			"email":      result.User.Email,
			"role":       result.User.Role,
			"avatar_url": result.User.AvatarUrl,
			"created_at": result.User.CreatedAt,
		},
	})
}

// Logout 用户登出
func (h *AuthHandler) Logout(c fiber.Ctx) error {
	return SuccessWithMessage(c, "登出成功", nil)
}

// RefreshToken 刷新Token
func (h *AuthHandler) RefreshToken(c fiber.Ctx) error {
	var req RefreshTokenRequest
	if err := c.Bind().Body(&req); err != nil {
		return Fail(c, errors.CodeInvalidParams, "请求参数格式错误")
	}

	if req.RefreshToken == "" {
		return Fail(c, errors.CodeInvalidParams, "刷新Token不能为空")
	}

	// 调用认证服务
	result, err := h.authClient.RefreshToken(context.Background(), req.RefreshToken)
	if err != nil {
		return FailWithGRPCError(c, err)
	}

	return Success(c, map[string]interface{}{
		"access_token":  result.AccessToken,
		"refresh_token": result.RefreshToken,
		"token_type":    result.TokenType,
		"expires_in":    result.ExpiresIn,
	})
}

// GetMe 获取当前用户信息
func (h *AuthHandler) GetMe(c fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	if userID == 0 {
		return Fail(c, errors.CodeUnauthorized, "未登录")
	}

	// 调用认证服务
	result, err := h.authClient.GetCurrentUser(context.Background(), userID)
	if err != nil {
		return FailWithGRPCError(c, err)
	}

	return Success(c, map[string]interface{}{
		"id":         result.Id,
		"username":   result.Username,
		"email":      result.Email,
		"role":       result.Role,
		"avatar_url": result.AvatarUrl,
		"created_at": result.CreatedAt,
	})
}
