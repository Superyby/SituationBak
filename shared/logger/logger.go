package logger

import (
	"os"
	"time"

	"SituationBak/shared/config"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var log *zap.Logger

// Init 初始化日志系统
func Init(cfg *config.LogConfig) error {
	// 解析日志级别
	level, err := zapcore.ParseLevel(cfg.Level)
	if err != nil {
		level = zapcore.InfoLevel
	}

	// 配置编码器
	var encoderConfig zapcore.EncoderConfig
	if cfg.Format == "json" {
		encoderConfig = zap.NewProductionEncoderConfig()
	} else {
		encoderConfig = zap.NewDevelopmentEncoderConfig()
	}
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder

	// 配置输出
	var cores []zapcore.Core

	// 控制台输出
	if cfg.Output == "stdout" || cfg.Output == "both" {
		consoleEncoder := zapcore.NewConsoleEncoder(encoderConfig)
		cores = append(cores, zapcore.NewCore(
			consoleEncoder,
			zapcore.AddSync(os.Stdout),
			level,
		))
	}

	// 文件输出
	if cfg.Output == "file" || cfg.Output == "both" {
		if cfg.FilePath != "" {
			file, err := os.OpenFile(cfg.FilePath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
			if err != nil {
				return err
			}

			fileEncoder := zapcore.NewJSONEncoder(encoderConfig)
			cores = append(cores, zapcore.NewCore(
				fileEncoder,
				zapcore.AddSync(file),
				level,
			))
		}
	}

	// 创建logger
	core := zapcore.NewTee(cores...)
	log = zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))

	return nil
}

// Sync 同步日志缓冲
func Sync() {
	if log != nil {
		_ = log.Sync()
	}
}

// Debug 调试日志
func Debug(msg string, fields ...zap.Field) {
	log.Debug(msg, fields...)
}

// Info 信息日志
func Info(msg string, fields ...zap.Field) {
	log.Info(msg, fields...)
}

// Warn 警告日志
func Warn(msg string, fields ...zap.Field) {
	log.Warn(msg, fields...)
}

// Error 错误日志
func Error(msg string, fields ...zap.Field) {
	log.Error(msg, fields...)
}

// Fatal 致命错误日志
func Fatal(msg string, fields ...zap.Field) {
	log.Fatal(msg, fields...)
}

// 便捷字段构造函数

// String 字符串字段
func String(key, val string) zap.Field {
	return zap.String(key, val)
}

// Int 整数字段
func Int(key string, val int) zap.Field {
	return zap.Int(key, val)
}

// Int64 64位整数字段
func Int64(key string, val int64) zap.Field {
	return zap.Int64(key, val)
}

// Uint 无符号整数字段
func Uint(key string, val uint) zap.Field {
	return zap.Uint(key, val)
}

// Float64 浮点数字段
func Float64(key string, val float64) zap.Field {
	return zap.Float64(key, val)
}

// Bool 布尔字段
func Bool(key string, val bool) zap.Field {
	return zap.Bool(key, val)
}

// Err 错误字段
func Err(err error) zap.Field {
	return zap.Error(err)
}

// Any 任意类型字段
func Any(key string, val interface{}) zap.Field {
	return zap.Any(key, val)
}

// Duration 时间间隔字段
func Duration(key string, val time.Duration) zap.Field {
	return zap.Duration(key, val)
}

// ==================== 链路追踪相关字段 ====================

// TraceID 链路追踪ID字段
func TraceID(traceID string) zap.Field {
	return zap.String("trace_id", traceID)
}

// UserID 用户ID字段
func UserID(userID uint) zap.Field {
	return zap.Uint("user_id", userID)
}

// Username 用户名字段
func Username(username string) zap.Field {
	return zap.String("username", username)
}

// RequestID 请求ID字段
func RequestID(requestID string) zap.Field {
	return zap.String("request_id", requestID)
}

// ==================== HTTP 相关字段 ====================

// Method HTTP 方法字段
func Method(method string) zap.Field {
	return zap.String("method", method)
}

// Path 请求路径字段
func Path(path string) zap.Field {
	return zap.String("path", path)
}

// Status HTTP 状态码字段
func Status(status int) zap.Field {
	return zap.Int("status", status)
}

// Latency 请求延迟字段
func Latency(latency time.Duration) zap.Field {
	return zap.Duration("latency", latency)
}

// IP 客户端IP字段
func IP(ip string) zap.Field {
	return zap.String("ip", ip)
}

// UserAgent 用户代理字段
func UserAgent(ua string) zap.Field {
	return zap.String("user_agent", ua)
}

// ==================== 业务相关字段 ====================

// Module 模块名称字段
func Module(module string) zap.Field {
	return zap.String("module", module)
}

// Action 操作名称字段
func Action(action string) zap.Field {
	return zap.String("action", action)
}

// Resource 资源名称字段
func Resource(resource string) zap.Field {
	return zap.String("resource", resource)
}

// ResourceID 资源ID字段
func ResourceID(id string) zap.Field {
	return zap.String("resource_id", id)
}

// ==================== 带上下文的日志函数 ====================

// WithFields 返回一个带字段的日志器
func WithFields(fields ...zap.Field) *zap.Logger {
	return log.With(fields...)
}

// WithTraceID 返回一个带 TraceID 的日志器
func WithTraceID(traceID string) *zap.Logger {
	return log.With(TraceID(traceID))
}

// WithUser 返回一个带用户信息的日志器
func WithUser(userID uint, username string) *zap.Logger {
	return log.With(UserID(userID), Username(username))
}

// WithContext 返回一个带完整上下文的日志器
func WithContext(traceID string, userID uint, username string) *zap.Logger {
	return log.With(TraceID(traceID), UserID(userID), Username(username))
}

// ==================== 时间相关字段 ====================

// Time 时间字段
func Time(key string, val time.Time) zap.Field {
	return zap.Time(key, val)
}

// Timestamp 时间戳字段（Unix秒）
func Timestamp(key string, val int64) zap.Field {
	return zap.Int64(key, val)
}

// ==================== 数组和对象字段 ====================

// Strings 字符串数组字段
func Strings(key string, val []string) zap.Field {
	return zap.Strings(key, val)
}

// Ints 整数数组字段
func Ints(key string, val []int) zap.Field {
	return zap.Ints(key, val)
}

// Object 对象字段（使用 JSON 序列化）
func Object(key string, val zapcore.ObjectMarshaler) zap.Field {
	return zap.Object(key, val)
}

// ==================== 辅助函数 ====================

// GetLogger 获取底层 zap.Logger
func GetLogger() *zap.Logger {
	return log
}

// SetLogger 设置底层 zap.Logger（用于测试）
func SetLogger(l *zap.Logger) {
	log = l
}

// NewNop 创建一个空的日志器（用于测试）
func NewNop() *zap.Logger {
	return zap.NewNop()
}

