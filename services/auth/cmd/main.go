package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	"SituationBak/pkg/config"
	"SituationBak/pkg/database"
	"SituationBak/pkg/logger"
	"SituationBak/pkg/model"
	"SituationBak/services/auth/internal/repository"
	"SituationBak/services/auth/internal/server"
	"SituationBak/services/auth/internal/service"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var (
	configPath = flag.String("config", "./configs/config.yaml", "配置文件路径")
	grpcPort   = flag.Int("port", 50051, "gRPC服务端口")
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

	logger.Info("启动认证服务",
		logger.String("service", "auth-svc"),
		logger.Int("port", *grpcPort),
	)

	// 初始化数据库
	if err := database.InitMySQL(&cfg.Database); err != nil {
		logger.Fatal("初始化MySQL失败", logger.Err(err))
	}
	defer database.Close()

	// 自动迁移
	if err := database.DB.AutoMigrate(&model.User{}, &model.UserSettings{}); err != nil {
		logger.Fatal("数据库迁移失败", logger.Err(err))
	}

	// 创建Repository
	userRepo := repository.NewUserRepository(database.DB)
	settingsRepo := repository.NewSettingsRepository(database.DB)

	// 创建Service
	authService := service.NewAuthService(userRepo, settingsRepo, &cfg.JWT)

	// 创建gRPC服务器
	grpcServer := grpc.NewServer()

	// 注册服务
	authServer := server.NewAuthServer(authService)
	server.RegisterAuthServiceServer(grpcServer, authServer)

	// 开启反射服务（用于grpcurl调试）
	reflection.Register(grpcServer)

	// 监听端口
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *grpcPort))
	if err != nil {
		logger.Fatal("监听端口失败", logger.Err(err), logger.Int("port", *grpcPort))
	}

	// 启动服务
	go func() {
		logger.Info("gRPC服务已启动", logger.Int("port", *grpcPort))
		if err := grpcServer.Serve(lis); err != nil {
			logger.Fatal("gRPC服务启动失败", logger.Err(err))
		}
	}()

	// 等待退出信号
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("正在关闭认证服务...")
	grpcServer.GracefulStop()
	logger.Info("认证服务已关闭")
}
