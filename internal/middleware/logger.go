package middleware

import (
	"time"

	"SituationBak/internal/pkg/logger"
	"SituationBak/internal/pkg/utils"
	"github.com/gofiber/fiber/v3"
)

// LoggerMiddleware 请求日志中间件
func LoggerMiddleware() fiber.Handler {
	return func(c fiber.Ctx) error {
		start := time.Now()
		requestID := utils.GenerateRequestID()

		// 设置请求ID到响应头
		c.Set("X-Request-ID", requestID)
		c.Locals("requestID", requestID)

		// 处理请求
		err := c.Next()

		// 计算耗时
		latency := time.Since(start)

		// 记录日志
		logger.Info("HTTP Request",
			logger.String("request_id", requestID),
			logger.String("method", c.Method()),
			logger.String("path", c.Path()),
			logger.Int("status", c.Response().StatusCode()),
			logger.Duration("latency", latency),
			logger.String("ip", c.IP()),
			logger.String("user_agent", c.Get("User-Agent")),
		)

		return err
	}
}

// GetRequestID 从上下文获取请求ID
func GetRequestID(c fiber.Ctx) string {
	requestID := c.Locals("requestID")
	if requestID == nil {
		return ""
	}
	return requestID.(string)
}
