package handler

import (
	"SituationBak/internal/dto/request"
	"SituationBak/internal/middleware"
	"SituationBak/internal/pkg/errors"
	"SituationBak/internal/pkg/utils"
	"SituationBak/internal/service"
	"github.com/gofiber/fiber/v3"
)

// UserHandler 用户处理器
type UserHandler struct {
	userService *service.UserService
}

// NewUserHandler 创建用户处理器实例
func NewUserHandler() *UserHandler {
	return &UserHandler{
		userService: service.NewUserService(),
	}
}

// GetProfile 获取用户资料
// @Summary 获取用户资料
// @Tags 用户
// @Security Bearer
// @Success 200 {object} utils.Response
// @Router /api/v1/user/profile [get]
func (h *UserHandler) GetProfile(c fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	if userID == 0 {
		return utils.FailWithCode(c, errors.CodeUnauthorized)
	}

	result, err := h.userService.GetProfile(userID)
	if err != nil {
		return utils.FailWithError(c, err)
	}

	return utils.Success(c, result)
}

// UpdateProfile 更新用户资料
// @Summary 更新用户资料
// @Tags 用户
// @Security Bearer
// @Accept json
// @Produce json
// @Param body body request.UpdateProfileRequest true "更新信息"
// @Success 200 {object} utils.Response
// @Router /api/v1/user/profile [put]
func (h *UserHandler) UpdateProfile(c fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	if userID == 0 {
		return utils.FailWithCode(c, errors.CodeUnauthorized)
	}

	var req request.UpdateProfileRequest
	if err := c.Bind().Body(&req); err != nil {
		return utils.Fail(c, errors.CodeInvalidParams, "请求参数格式错误")
	}

	result, err := h.userService.UpdateProfile(userID, &req)
	if err != nil {
		return utils.FailWithError(c, err)
	}

	return utils.Success(c, result)
}

// ChangePassword 修改密码
// @Summary 修改密码
// @Tags 用户
// @Security Bearer
// @Accept json
// @Produce json
// @Param body body request.ChangePasswordRequest true "密码信息"
// @Success 200 {object} utils.Response
// @Router /api/v1/user/password [put]
func (h *UserHandler) ChangePassword(c fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	if userID == 0 {
		return utils.FailWithCode(c, errors.CodeUnauthorized)
	}

	var req request.ChangePasswordRequest
	if err := c.Bind().Body(&req); err != nil {
		return utils.Fail(c, errors.CodeInvalidParams, "请求参数格式错误")
	}

	if req.OldPassword == "" || req.NewPassword == "" {
		return utils.Fail(c, errors.CodeInvalidParams, "旧密码和新密码不能为空")
	}

	if len(req.NewPassword) < 6 {
		return utils.Fail(c, errors.CodeInvalidParams, "新密码长度至少6个字符")
	}

	err := h.userService.ChangePassword(userID, &req)
	if err != nil {
		return utils.FailWithError(c, err)
	}

	return utils.SuccessWithMessage(c, "密码修改成功", nil)
}

// GetSettings 获取用户设置
// @Summary 获取用户设置
// @Tags 用户
// @Security Bearer
// @Success 200 {object} utils.Response
// @Router /api/v1/user/settings [get]
func (h *UserHandler) GetSettings(c fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	if userID == 0 {
		return utils.FailWithCode(c, errors.CodeUnauthorized)
	}

	result, err := h.userService.GetSettings(userID)
	if err != nil {
		return utils.FailWithError(c, err)
	}

	return utils.Success(c, result)
}

// UpdateSettings 更新用户设置
// @Summary 更新用户设置
// @Tags 用户
// @Security Bearer
// @Accept json
// @Produce json
// @Param body body request.UpdateSettingsRequest true "设置信息"
// @Success 200 {object} utils.Response
// @Router /api/v1/user/settings [put]
func (h *UserHandler) UpdateSettings(c fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	if userID == 0 {
		return utils.FailWithCode(c, errors.CodeUnauthorized)
	}

	var req request.UpdateSettingsRequest
	if err := c.Bind().Body(&req); err != nil {
		return utils.Fail(c, errors.CodeInvalidParams, "请求参数格式错误")
	}

	result, err := h.userService.UpdateSettings(userID, &req)
	if err != nil {
		return utils.FailWithError(c, err)
	}

	return utils.Success(c, result)
}
