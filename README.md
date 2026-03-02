# 太空态势系统后端
## 启动命令
`go run cmd/server/main.go`


## 注意：
Enable Go modules integration 的作用
当你勾选了 Enable Go modules integration,GoLand 会:
1. 识别 go.mod 文件 - 将你的项目识别为 Go Modules 项目
2. 到 GOPATH/pkg/mod 查找依赖 - 即使你的项目不在 GOPATH 下,GoLand 也会去:
   - D:\yby\app\person\gopath\pkg\mod\ 目录下查找所有依赖包
   - 这是 Go Modules 的依赖缓存目录
   正确解析导入路径 - 根据 go.mod 中的模块名 SituationBak 来解析项目内部的导入


```text
你的项目: D:\yby\app\person\project\SituationBak\SituationBak\
├── go.mod  (module SituationBak)
├── internal/
└── cmd/

依赖缓存: D:\yby\app\person\gopath\pkg\mod\
├── github.com/gofiber/fiber/v3@v3.0.0-beta.3/
├── github.com/redis/go-redis/v9@v9.5.0/
├── go.uber.org/zap@v1.27.0/
└── ... (所有其他依赖)
```

## 项目目录结构
```text
SituationBak/
├── api/proto/                     # gRPC Proto定义
│   ├── auth/v1/auth.proto
│   ├── user/v1/user.proto
│   ├── satellite/v1/satellite.proto
│   └── favorite/v1/favorite.proto
├── services/                      # 微服务
│   ├── gateway/                   # API网关 (HTTP:4000)
│   ├── auth/                      # 认证服务 (gRPC:50051)
│   ├── user/                      # 用户服务 (gRPC:50052)
│   ├── satellite/                 # 卫星服务 (gRPC:50053)
│   └── favorite/                  # 收藏服务 (gRPC:50054)
├── pkg/                           # 共享包
│   ├── config/                    # 配置管理
│   ├── database/                  # MySQL/Redis连接
│   ├── logger/                    # 日志
│   ├── errors/                    # 错误码
│   ├── model/                     # 数据模型
│   └── utils/                     # 工具函数
├── deployments/docker/            # Docker部署
│   └── docker-compose.yaml
├── go.work                        # Go Workspace
└── Makefile                       # 构建脚本
```

## 主要内容
1. Proto 接口定义 - 4个服务的gRPC接口
2. 共享包迁移 - config, logger, errors, database, model, utils
3. 认证服务(auth-svc) - 完整实现注册/登录/Token刷新/验证
4. API网关(gateway) - HTTP路由 + gRPC客户端 + JWT中间件
5. 其他服务框架 - user-svc, satellite-svc, favorite-svc 基础框架
6. 基础设施 - docker-compose (MySQL, Redis, Consul)
7. 构建工具 - Makefile, go.work

## 后续步骤
1. 运行 make proto 生成gRPC代码（需安装protoc）
2. 运行 make infra-up 启动基础设施
3. 完善 user-svc, satellite-svc, favorite-svc 的业务逻辑
4. 运行 make docker-up 启动所有服务