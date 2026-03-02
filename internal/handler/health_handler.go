package handler

import (
	"runtime"
	"time"

	"SituationBak/internal/repository"
	"SituationBak/shared/constants"

	"github.com/gofiber/fiber/v3"
)

// HealthHandler 健康检查处理器
type HealthHandler struct {
	startTime time.Time
}

// NewHealthHandler 创建健康检查处理器实例
func NewHealthHandler() *HealthHandler {
	return &HealthHandler{
		startTime: time.Now(),
	}
}

// HealthResponse 健康检查响应
type HealthResponse struct {
	Status    string `json:"status"`
	Timestamp string `json:"timestamp"`
	Uptime    string `json:"uptime"`
	Version   string `json:"version"`
}

// ReadyResponse 就绪检查响应
type ReadyResponse struct {
	Status   string                   `json:"status"`
	Services map[string]ServiceStatus `json:"services"`
}

// ServiceStatus 服务状态
type ServiceStatus struct {
	Status  string `json:"status"`
	Message string `json:"message,omitempty"`
}

// Health 基本健康检查
// @Summary 健康检查
// @Description 返回服务基本健康状态
// @Tags 系统
// @Produce json
// @Success 200 {object} HealthResponse
// @Router /health [get]
func (h *HealthHandler) Health(c fiber.Ctx) error {
	uptime := time.Since(h.startTime)

	return c.JSON(HealthResponse{
		Status:    "healthy",
		Timestamp: time.Now().Format(constants.TimeFormatISO8601),
		Uptime:    uptime.String(),
		Version:   constants.AppVersion,
	})
}

// Ready 就绪检查（检查所有依赖服务）
// @Summary 就绪检查
// @Description 检查所有依赖服务（数据库、Redis等）的连接状态
// @Tags 系统
// @Produce json
// @Success 200 {object} ReadyResponse
// @Failure 503 {object} ReadyResponse
// @Router /ready [get]
func (h *HealthHandler) Ready(c fiber.Ctx) error {
	services := make(map[string]ServiceStatus)
	allHealthy := true

	// 检查 MySQL
	if repository.IsDBConnected() {
		services["mysql"] = ServiceStatus{
			Status: "healthy",
		}
	} else {
		services["mysql"] = ServiceStatus{
			Status:  "unhealthy",
			Message: "数据库连接失败",
		}
		allHealthy = false
	}

	// 检查 Redis
	if repository.IsRedisConnected() {
		services["redis"] = ServiceStatus{
			Status: "healthy",
		}
	} else {
		services["redis"] = ServiceStatus{
			Status:  "degraded",
			Message: "Redis连接失败，部分功能可能受影响",
		}
		// Redis 不是必需的，所以不标记为不健康
	}

	status := "ready"
	httpStatus := fiber.StatusOK
	if !allHealthy {
		status = "not_ready"
		httpStatus = fiber.StatusServiceUnavailable
	}

	return c.Status(httpStatus).JSON(ReadyResponse{
		Status:   status,
		Services: services,
	})
}

// Live 存活检查
// @Summary 存活检查
// @Description 简单的存活探针，用于K8s liveness probe
// @Tags 系统
// @Produce json
// @Success 200 {object} map[string]string
// @Router /live [get]
func (h *HealthHandler) Live(c fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"status": "alive",
	})
}

// Info 系统信息
// @Summary 系统信息
// @Description 返回系统运行信息
// @Tags 系统
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /info [get]
func (h *HealthHandler) Info(c fiber.Ctx) error {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	return c.JSON(fiber.Map{
		"app": fiber.Map{
			"name":    constants.AppName,
			"version": constants.AppVersion,
		},
		"runtime": fiber.Map{
			"go_version":    runtime.Version(),
			"go_os":         runtime.GOOS,
			"go_arch":       runtime.GOARCH,
			"num_cpu":       runtime.NumCPU(),
			"num_goroutine": runtime.NumGoroutine(),
		},
		"memory": fiber.Map{
			"alloc_mb":       memStats.Alloc / 1024 / 1024,
			"total_alloc_mb": memStats.TotalAlloc / 1024 / 1024,
			"sys_mb":         memStats.Sys / 1024 / 1024,
			"num_gc":         memStats.NumGC,
		},
		"uptime": time.Since(h.startTime).String(),
	})
}
