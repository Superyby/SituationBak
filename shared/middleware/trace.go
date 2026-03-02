package middleware

import (
	"SituationBak/shared/constants"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

// TraceMiddleware 链路追踪中间件
// 为每个请求生成或传递唯一的 TraceID，便于日志追踪和问题排查
func TraceMiddleware() fiber.Handler {
	return func(c fiber.Ctx) error {
		// 优先从请求头获取 TraceID（支持分布式追踪）
		traceID := c.Get(constants.HeaderTraceID)
		if traceID == "" {
			// 也尝试获取 X-Request-ID
			traceID = c.Get(constants.HeaderRequestID)
		}
		if traceID == "" {
			// 生成新的 TraceID
			traceID = uuid.New().String()
		}

		// 存入上下文
		c.Locals(constants.CtxKeyTraceID, traceID)

		// 设置响应头，便于客户端追踪
		c.Set(constants.HeaderTraceID, traceID)

		return c.Next()
	}
}

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
