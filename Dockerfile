# Orbital Tracker API

FROM golang:1.22-alpine AS builder

# 设置工作目录
WORKDIR /app

# 安装必要的系统依赖
RUN apk add --no-cache git ca-certificates tzdata

# 复制 go.mod 和 go.sum
COPY go.mod go.sum ./

# 下载依赖
RUN go mod download

# 复制源代码
COPY . .

# 构建应用
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o server ./cmd/server

# 最终镜像
FROM alpine:3.19

# 设置时区
RUN apk --no-cache add ca-certificates tzdata

WORKDIR /app

# 从构建阶段复制二进制文件
COPY --from=builder /app/server .

# 复制配置文件
COPY configs/ ./configs/

# 暴露端口
EXPOSE 4000

# 设置环境变量
ENV APP_ENV=production

# 启动命令
CMD ["./server"]
