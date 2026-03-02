package middleware

import (
	"runtime/debug"

	"SituationBak/internal/pkg/errors"
	"SituationBak/internal/pkg/logger"
	"SituationBak/internal/pkg/utils"
	"github.com/gofiber/fiber/v3"
)

// RecoveryMiddleware 异常恢复中间件
func RecoveryMiddleware() fiber.Handler {
	return func(c fiber.Ctx) error {
		defer func() {
			if r := recover(); r != nil {
				// 记录堆栈信息
				stack := string(debug.Stack())
				logger.Error("Panic recovered",
					logger.Any("error", r),
					logger.String("stack", stack),
					logger.String("path", c.Path()),
					logger.String("method", c.Method()),
				)

				// 返回错误响应
				_ = utils.Fail(c, errors.CodeInternalError, "服务器内部错误")
			}
		}()

		return c.Next()
	}
}
