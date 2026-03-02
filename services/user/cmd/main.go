package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	"SituationBak/shared/config"
	"SituationBak/shared/database"
	"SituationBak/shared/logger"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var (
	configPath = flag.String("config", "./configs/config.yaml", "配置文件路径")
	grpcPort   = flag.Int("port", 50052, "gRPC服务端口")
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

	logger.Info("启动用户服务",
		logger.String("service", "user-svc"),
		logger.Int("port", *grpcPort),
	)

	// 初始化数据库
	if _, err := database.InitMySQL(&cfg.Database); err != nil {
		logger.Fatal("初始化MySQL失败", logger.Err(err))
	}
	defer database.Close()

	// 创建gRPC服务器
	grpcServer := grpc.NewServer()

	// TODO: 注册用户服务
	// userServer := server.NewUserServer(...)
	// server.RegisterUserServiceServer(grpcServer, userServer)

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

	logger.Info("正在关闭用户服务...")
	grpcServer.GracefulStop()
	logger.Info("用户服务已关闭")
}
