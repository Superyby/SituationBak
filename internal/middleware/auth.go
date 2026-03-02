package middleware

import (
	"strings"

	"SituationBak/shared/errors"
	"SituationBak/shared/utils"
	"SituationBak/internal/service"
	"github.com/gofiber/fiber/v3"
)

// AuthMiddleware JWT认证中间�?
func AuthMiddleware() fiber.Handler {
	authService := service.NewAuthService()

	return func(c fiber.Ctx) error {
		// 获取Authorization�?
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return utils.FailWithCode(c, errors.CodeUnauthorized)
		}

		// 检查Bearer前缀
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			return utils.FailWithCode(c, errors.CodeTokenInvalid)
		}

		tokenString := parts[1]
		if tokenString == "" {
			return utils.FailWithCode(c, errors.CodeTokenInvalid)
		}

		// 验证Token
		claims, err := authService.ValidateToken(tokenString)
		if err != nil {
			return utils.FailWithCode(c, errors.CodeTokenExpired)
		}

		// 将用户信息存入上下文
		c.Locals("userID", claims.UserID)
		c.Locals("username", claims.Username)
		c.Locals("role", claims.Role)

		return c.Next()
	}
}

// OptionalAuthMiddleware 可选认证中间件（不强制登录�?
func OptionalAuthMiddleware() fiber.Handler {
	authService := service.NewAuthService()

	return func(c fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Next()
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			return c.Next()
		}

		tokenString := parts[1]
		if tokenString == "" {
			return c.Next()
		}

		claims, err := authService.ValidateToken(tokenString)
		if err == nil {
			c.Locals("userID", claims.UserID)
			c.Locals("username", claims.Username)
			c.Locals("role", claims.Role)
		}

		return c.Next()
	}
}

// AdminMiddleware 管理员权限中间件
func AdminMiddleware() fiber.Handler {
	return func(c fiber.Ctx) error {
		role := c.Locals("role")
		if role == nil || role.(string) != "admin" {
			return utils.FailWithCode(c, errors.CodeForbidden)
		}
		return c.Next()
	}
}

// GetUserID 从上下文获取用户ID
func GetUserID(c fiber.Ctx) uint {
	userID := c.Locals("userID")
	if userID == nil {
		return 0
	}
	return userID.(uint)
}

// GetUsername 从上下文获取用户�?
func GetUsername(c fiber.Ctx) string {
	username := c.Locals("username")
	if username == nil {
		return ""
	}
	return username.(string)
}

// GetUserRole 从上下文获取用户角色
func GetUserRole(c fiber.Ctx) string {
	role := c.Locals("role")
	if role == nil {
		return ""
	}
	return role.(string)
}
