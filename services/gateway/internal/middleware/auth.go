package middleware

import (
	"context"
	"strings"

	"SituationBak/pkg/errors"
	"SituationBak/services/gateway/internal/client"

	"github.com/gofiber/fiber/v3"
)

// UserContextKey 用户信息context key
const UserContextKey = "user"

// UserContext 用户上下文信息
type UserContext struct {
	UserID   uint64
	Username string
	Role     string
}

// AuthMiddleware 认证中间件
func AuthMiddleware(authClient *client.AuthClient) fiber.Handler {
	return func(c fiber.Ctx) error {
		// 获取Authorization头
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return respondUnauthorized(c, "缺少认证信息")
		}

		// 解析Bearer Token
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			return respondUnauthorized(c, "无效的认证格式")
		}

		token := parts[1]
		if token == "" {
			return respondUnauthorized(c, "Token不能为空")
		}

		// 调用认证服务验证Token
		resp, err := authClient.ValidateToken(context.Background(), token)
		if err != nil {
			return respondUnauthorized(c, "Token验证失败")
		}

		if !resp.Valid {
			return respondUnauthorized(c, "Token无效或已过期")
		}

		// 将用户信息存入Context
		c.Locals(UserContextKey, &UserContext{
			UserID:   resp.UserId,
			Username: resp.Username,
			Role:     resp.Role,
		})

		return c.Next()
	}
}

// GetUserContext 从Context获取用户信息
func GetUserContext(c fiber.Ctx) *UserContext {
	user, ok := c.Locals(UserContextKey).(*UserContext)
	if !ok {
		return nil
	}
	return user
}

// GetUserID 获取当前用户ID
func GetUserID(c fiber.Ctx) uint64 {
	user := GetUserContext(c)
	if user == nil {
		return 0
	}
	return user.UserID
}

// respondUnauthorized 返回未授权响应
func respondUnauthorized(c fiber.Ctx, message string) error {
	return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
		"code":    errors.CodeUnauthorized,
		"message": message,
		"data":    nil,
	})
}

// RecoveryMiddleware 恢复中间件
func RecoveryMiddleware() fiber.Handler {
	return func(c fiber.Ctx) error {
		defer func() {
			if r := recover(); r != nil {
				c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"code":    errors.CodeInternalError,
					"message": "服务器内部错误",
					"data":    nil,
				})
			}
		}()
		return c.Next()
	}
}

// CORSMiddleware CORS中间件
func CORSMiddleware() fiber.Handler {
	return func(c fiber.Ctx) error {
		c.Set("Access-Control-Allow-Origin", "*")
		c.Set("Access-Control-Allow-Methods", "GET,POST,PUT,DELETE,OPTIONS")
		c.Set("Access-Control-Allow-Headers", "Content-Type,Authorization")
		c.Set("Access-Control-Max-Age", "86400")

		if c.Method() == "OPTIONS" {
			return c.SendStatus(fiber.StatusNoContent)
		}

		return c.Next()
	}
}
