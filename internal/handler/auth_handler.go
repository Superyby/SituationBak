package handler

import (
	"SituationBak/internal/dto/request"
	"SituationBak/internal/middleware"
	"SituationBak/shared/errors"
	"SituationBak/shared/utils"
	"SituationBak/internal/service"
	"github.com/gofiber/fiber/v3"
)

// AuthHandler 认证处理器
type AuthHandler struct {
	authService *service.AuthService
}

// NewAuthHandler 创建认证处理器实例
func NewAuthHandler() *AuthHandler {
	return &AuthHandler{
		authService: service.NewAuthService(),
	}
}

// Register 用户注册
// @Summary 用户注册
// @Tags 认证
// @Accept json
// @Produce json
// @Param body body request.RegisterRequest true "注册信息"
// @Success 200 {object} utils.Response
// @Router /api/v1/auth/register [post]
func (h *AuthHandler) Register(c fiber.Ctx) error {
	var req request.RegisterRequest
	if err := c.Bind().Body(&req); err != nil {
		return utils.Fail(c, errors.CodeInvalidParams, "请求参数格式错误")
	}

	// 参数验证
	if req.Username == "" || req.Email == "" || req.Password == "" {
		return utils.Fail(c, errors.CodeInvalidParams, "用户名、邮箱和密码不能为空")
	}

	if len(req.Username) < 3 || len(req.Username) > 50 {
		return utils.Fail(c, errors.CodeInvalidParams, "用户名长度应在3-50个字符之间")
	}

	if len(req.Password) < 6 {
		return utils.Fail(c, errors.CodeInvalidParams, "密码长度至少6个字符")
	}

	result, err := h.authService.Register(&req)
	if err != nil {
		return utils.FailWithError(c, err)
	}

	return utils.Created(c, result)
}

// Login 用户登录
// @Summary 用户登录
// @Tags 认证
// @Accept json
// @Produce json
// @Param body body request.LoginRequest true "登录信息"
// @Success 200 {object} utils.Response
// @Router /api/v1/auth/login [post]
func (h *AuthHandler) Login(c fiber.Ctx) error {
	var req request.LoginRequest
	if err := c.Bind().Body(&req); err != nil {
		return utils.Fail(c, errors.CodeInvalidParams, "请求参数格式错误")
	}

	if req.Username == "" || req.Password == "" {
		return utils.Fail(c, errors.CodeInvalidParams, "用户名和密码不能为空")
	}

	result, err := h.authService.Login(&req)
	if err != nil {
		return utils.FailWithError(c, err)
	}

	return utils.Success(c, result)
}

// Logout 用户登出
// @Summary 用户登出
// @Tags 认证
// @Security Bearer
// @Success 200 {object} utils.Response
// @Router /api/v1/auth/logout [post]
func (h *AuthHandler) Logout(c fiber.Ctx) error {
	// JWT无状态，登出只需客户端清除Token
	// 如需实现Token黑名单，可在此处将Token加入黑名单
	return utils.SuccessWithMessage(c, "登出成功", nil)
}

// RefreshToken 刷新Token
// @Summary 刷新Token
// @Tags 认证
// @Accept json
// @Produce json
// @Param body body request.RefreshTokenRequest true "刷新Token"
// @Success 200 {object} utils.Response
// @Router /api/v1/auth/refresh [post]
func (h *AuthHandler) RefreshToken(c fiber.Ctx) error {
	var req request.RefreshTokenRequest
	if err := c.Bind().Body(&req); err != nil {
		return utils.Fail(c, errors.CodeInvalidParams, "请求参数格式错误")
	}

	if req.RefreshToken == "" {
		return utils.Fail(c, errors.CodeInvalidParams, "刷新Token不能为空")
	}

	result, err := h.authService.RefreshToken(req.RefreshToken)
	if err != nil {
		return utils.FailWithError(c, err)
	}

	return utils.Success(c, result)
}

// GetMe 获取当前用户信息
// @Summary 获取当前用户信息
// @Tags 认证
// @Security Bearer
// @Success 200 {object} utils.Response
// @Router /api/v1/auth/me [get]
func (h *AuthHandler) GetMe(c fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	if userID == 0 {
		return utils.FailWithCode(c, errors.CodeUnauthorized)
	}

	result, err := h.authService.GetCurrentUser(userID)
	if err != nil {
		return utils.FailWithError(c, err)
	}

	return utils.Success(c, result)
}
