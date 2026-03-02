package utils

import (
	"SituationBak/shared/constants"

	"github.com/gofiber/fiber/v3"
)

// GetTraceID 从 Fiber 上下文获取 TraceID
func GetTraceID(c fiber.Ctx) string {
	traceID := c.Locals(constants.CtxKeyTraceID)
	if traceID == nil {
		return ""
	}
	if s, ok := traceID.(string); ok {
		return s
	}
	return ""
}

// GetUserID 从 Fiber 上下文获取用户ID
func GetUserID(c fiber.Ctx) uint {
	userID := c.Locals(constants.CtxKeyUserID)
	if userID == nil {
		return 0
	}
	if id, ok := userID.(uint); ok {
		return id
	}
	return 0
}

// GetUsername 从 Fiber 上下文获取用户名
func GetUsername(c fiber.Ctx) string {
	username := c.Locals(constants.CtxKeyUsername)
	if username == nil {
		return ""
	}
	if s, ok := username.(string); ok {
		return s
	}
	return ""
}

// GetUserRole 从 Fiber 上下文获取用户角色
func GetUserRole(c fiber.Ctx) string {
	role := c.Locals(constants.CtxKeyRole)
	if role == nil {
		return ""
	}
	if s, ok := role.(string); ok {
		return s
	}
	return ""
}

// SetUserInfo 设置用户信息到上下文
func SetUserInfo(c fiber.Ctx, userID uint, username, role string) {
	c.Locals(constants.CtxKeyUserID, userID)
	c.Locals(constants.CtxKeyUsername, username)
	c.Locals(constants.CtxKeyRole, role)
}

// SetTraceID 设置 TraceID 到上下文
func SetTraceID(c fiber.Ctx, traceID string) {
	c.Locals(constants.CtxKeyTraceID, traceID)
}

// IsAuthenticated 检查用户是否已认证
func IsAuthenticated(c fiber.Ctx) bool {
	return GetUserID(c) > 0
}

// IsAdmin 检查用户是否为管理员
func IsAdmin(c fiber.Ctx) bool {
	return GetUserRole(c) == constants.RoleAdmin
}
