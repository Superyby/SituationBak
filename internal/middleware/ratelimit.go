package middleware

import (
	"sync"
	"time"

	"SituationBak/internal/config"
	"SituationBak/shared/errors"
	"SituationBak/shared/utils"
	"github.com/gofiber/fiber/v3"
)

// 简单的令牌桶限流器
type rateLimiter struct {
	tokens     float64
	maxTokens  float64
	refillRate float64 // 每秒补充的token数
	lastRefill time.Time
	mu         sync.Mutex
}

var (
	limiters  = make(map[string]*rateLimiter)
	limiterMu sync.Mutex
)

// RateLimitMiddleware 限流中间件
func RateLimitMiddleware() fiber.Handler {
	cfg := config.GlobalConfig.RateLimit

	return func(c fiber.Ctx) error {
		ip := c.IP()

		limiterMu.Lock()
		limiter, exists := limiters[ip]
		if !exists {
			limiter = &rateLimiter{
				tokens:     float64(cfg.Burst),
				maxTokens:  float64(cfg.Burst),
				refillRate: float64(cfg.RequestsPerSecond),
				lastRefill: time.Now(),
			}
			limiters[ip] = limiter
		}
		limiterMu.Unlock()

		if !limiter.allow() {
			return utils.FailWithCode(c, errors.CodeTooManyRequests)
		}

		return c.Next()
	}
}

// allow 检查是否允许请求
func (r *rateLimiter) allow() bool {
	r.mu.Lock()
	defer r.mu.Unlock()

	now := time.Now()
	elapsed := now.Sub(r.lastRefill).Seconds()
	r.lastRefill = now

	// 补充token
	r.tokens += elapsed * r.refillRate
	if r.tokens > r.maxTokens {
		r.tokens = r.maxTokens
	}

	// 检查是否有足够的token
	if r.tokens >= 1 {
		r.tokens--
		return true
	}

	return false
}

// CleanupLimiters 清理过期的限流器（可在定时任务中调用）
func CleanupLimiters() {
	limiterMu.Lock()
	defer limiterMu.Unlock()

	expireTime := time.Now().Add(-5 * time.Minute)
	for ip, limiter := range limiters {
		limiter.mu.Lock()
		if limiter.lastRefill.Before(expireTime) {
			delete(limiters, ip)
		}
		limiter.mu.Unlock()
	}
}
