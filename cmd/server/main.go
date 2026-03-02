package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"SituationBak/internal/config"
	"SituationBak/internal/repository"
	"SituationBak/internal/router"
	"SituationBak/shared/constants"
	"SituationBak/shared/graceful"
	"SituationBak/shared/logger"

	"github.com/gofiber/fiber/v3"

	_ "SituationBak/docs" // Swagger docs
)

// @title Orbital Tracker API
// @version 1.0.0
// @description 轨道追踪器后端API 服务，提供卫星数据查询、用户认证、收藏管理等功能
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.email support@orbital-tracker.com

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:4000
// @BasePath /api/v1

// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization
// @description JWT 认证令牌，格式：Bearer {token}

func main() {
	// 加载配置
	cfg, err := config.LoadConfig("")
	if err != nil {
		fmt.Printf("加载配置失败: %v\n", err)
		os.Exit(1)
	}

	// 初始化日志
	if err := logger.Init(&cfg.Log); err != nil {
		fmt.Printf("初始化日志失败: %v\n", err)
		os.Exit(1)
	}
	defer logger.Sync()

	logger.Info("Starting Orbital Tracker API",
		logger.String("app", constants.AppName),
		logger.String("version", constants.AppVersion),
		logger.String("env", cfg.App.Env),
		logger.Int("port", cfg.App.Port),
	)

	// 初始化MySQL
	if err := repository.InitMySQL(&cfg.Database); err != nil {
		logger.Fatal("初始化MySQL失败", logger.Err(err))
	}
	logger.Info("MySQL connected successfully")

	// 注册数据库到优雅关闭管理器
	graceful.Register("mysql", &dbCloser{})

	// 自动迁移数据库表
	if err := repository.AutoMigrate(); err != nil {
		logger.Fatal("数据库迁移失败", logger.Err(err))
	}
	logger.Info("Database migration completed")

	// 初始化Redis（可选，连接失败不影响启动）
	if err := repository.InitRedis(&cfg.Redis); err != nil {
		logger.Warn("初始化Redis失败，部分缓存功能将不可用", logger.Err(err))
	} else {
		logger.Info("Redis connected successfully")
		graceful.Register("redis", &redisCloser{})
	}

	// 创建Fiber应用
	app := fiber.New(fiber.Config{
		AppName:      cfg.App.Name,
		ErrorHandler: customErrorHandler,
	})

	// 配置路由
	router.SetupRoutes(app)

	// 优雅关闭
	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		<-sigChan

		logger.Info("Shutting down server...")
		if err := app.Shutdown(); err != nil {
			logger.Error("Server shutdown error", logger.Err(err))
		}
	}()

	// 启动服务器
	addr := fmt.Sprintf(":%d", cfg.App.Port)
	logger.Info("Server starting", logger.String("address", addr))

	if err := app.Listen(addr); err != nil {
		logger.Fatal("Server failed to start", logger.Err(err))
	}

	// 优雅关闭所有资源
	if err := graceful.Close(); err != nil {
		logger.Error("Error closing resources", logger.Err(err))
	}

	logger.Info("Server stopped")
}

// dbCloser 数据库关闭器
type dbCloser struct{}

func (c *dbCloser) Close() error {
	return repository.CloseMySQL()
}

// redisCloser Redis关闭器
type redisCloser struct{}

func (c *redisCloser) Close() error {
	return repository.CloseRedis()
}

// customErrorHandler 自定义错误处理器
func customErrorHandler(c fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError
	message := "服务器内部错误"

	if e, ok := err.(*fiber.Error); ok {
		code = e.Code
		message = e.Message
	}

	logger.Error("HTTP Error",
		logger.Int("status", code),
		logger.String("message", message),
		logger.String("path", c.Path()),
	)

	return c.Status(code).JSON(fiber.Map{
		"code":    code,
		"message": message,
		"success": false,
	})
}
