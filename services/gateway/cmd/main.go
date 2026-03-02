package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"SituationBak/pkg/config"
	"SituationBak/pkg/logger"
	"SituationBak/services/gateway/internal/client"
	"SituationBak/services/gateway/internal/handler"
	"SituationBak/services/gateway/internal/middleware"

	"github.com/gofiber/fiber/v3"
)

var (
	configPath = flag.String("config", "./configs/config.yaml", "配置文件路径")
	httpPort   = flag.Int("port", 4000, "HTTP服务端口")
)

func main() {
	flag.Parse()

	// 加载配置
	cfg, err := config.LoadConfig(*configPath)
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

	logger.Info("启动API网关",
		logger.String("service", "gateway"),
		logger.Int("port", *httpPort),
	)

	// 连接gRPC服务
	authClient, err := client.NewAuthClient(cfg.GRPC.AuthAddr)
	if err != nil {
		logger.Fatal("连接认证服务失败", logger.Err(err), logger.String("addr", cfg.GRPC.AuthAddr))
	}
	defer authClient.Close()

	logger.Info("已连接认证服务", logger.String("addr", cfg.GRPC.AuthAddr))

	// 创建Fiber应用
	app := fiber.New(fiber.Config{
		AppName: "SituationBak Gateway",
	})

	// 全局中间件
	app.Use(middleware.RecoveryMiddleware())
	app.Use(middleware.CORSMiddleware())

	// 健康检查
	app.Get("/health", func(c fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "ok",
			"service": "gateway",
		})
	})

	// API v1
	api := app.Group("/api/v1")

	// 初始化Handler
	authHandler := handler.NewAuthHandler(authClient)

	// 认证路由（无需认证）
	auth := api.Group("/auth")
	auth.Post("/register", authHandler.Register)
	auth.Post("/login", authHandler.Login)
	auth.Post("/refresh", authHandler.RefreshToken)

	// 认证路由（需要认证）
	authProtected := api.Group("/auth", middleware.AuthMiddleware(authClient))
	authProtected.Post("/logout", authHandler.Logout)
	authProtected.Get("/me", authHandler.GetMe)

	// TODO: 添加其他服务路由
	// user := api.Group("/user", middleware.AuthMiddleware(authClient))
	// satellites := api.Group("/satellites")
	// favorites := api.Group("/favorites", middleware.AuthMiddleware(authClient))

	// 启动HTTP服务
	go func() {
		addr := fmt.Sprintf(":%d", *httpPort)
		logger.Info("HTTP服务已启动", logger.String("addr", addr))
		if err := app.Listen(addr); err != nil {
			logger.Fatal("HTTP服务启动失败", logger.Err(err))
		}
	}()

	// 等待退出信号
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("正在关闭API网关...")
	if err := app.Shutdown(); err != nil {
		logger.Error("关闭HTTP服务失败", logger.Err(err))
	}
	logger.Info("API网关已关闭")
}
