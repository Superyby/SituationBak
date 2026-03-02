package config

import (
	"fmt"
	"time"

	sharedConfig "SituationBak/shared/config"
)

// ==================== 类型别名（重导出 shared/config） ====================

type (
	// Config 应用配置结构体
	Config = sharedConfig.Config
	// AppConfig 应用配置
	AppConfig = sharedConfig.AppConfig
	// DatabaseConfig MySQL数据库配置
	DatabaseConfig = sharedConfig.DatabaseConfig
	// RedisConfig Redis配置
	RedisConfig = sharedConfig.RedisConfig
	// ClickHouseConfig ClickHouse配置
	ClickHouseConfig = sharedConfig.ClickHouseConfig
	// JWTConfig JWT配置
	JWTConfig = sharedConfig.JWTConfig
	// LogConfig 日志配置
	LogConfig = sharedConfig.LogConfig
	// RateLimitConfig 限流配置
	RateLimitConfig = sharedConfig.RateLimitConfig
	// ExternalConfig 第三方API配置
	ExternalConfig = sharedConfig.ExternalConfig
	// KeepTrackConfig KeepTrack API配置
	KeepTrackConfig = sharedConfig.KeepTrackConfig
	// SpaceTrackConfig SpaceTrack API配置
	SpaceTrackConfig = sharedConfig.SpaceTrackConfig
	// CORSConfig CORS配置
	CORSConfig = sharedConfig.CORSConfig
	// GRPCConfig gRPC服务地址配置
	GRPCConfig = sharedConfig.GRPCConfig
)

// ==================== 全局配置 ====================

// GlobalConfig 全局配置实例指针
// 指向 shared/config 中的 GlobalConfig
var GlobalConfig *Config

// GetGlobalConfig 获取全局配置
func GetGlobalConfig() *Config {
	if GlobalConfig != nil {
		return GlobalConfig
	}
	return sharedConfig.GlobalConfig
}

// ==================== 重导出辅助函数 ====================

// IsDevelopment 判断是否为开发环境
func IsDevelopment() bool {
	cfg := GetGlobalConfig()
	if cfg != nil {
		return cfg.App.Env == "development"
	}
	return sharedConfig.IsDevelopment()
}

// IsProduction 判断是否为生产环境
func IsProduction() bool {
	cfg := GetGlobalConfig()
	if cfg != nil {
		return cfg.App.Env == "production"
	}
	return sharedConfig.IsProduction()
}

// IsTesting 判断是否为测试环境
func IsTesting() bool {
	cfg := GetGlobalConfig()
	if cfg != nil {
		return cfg.App.Env == "testing"
	}
	return sharedConfig.IsTesting()
}

// GetEnv 获取环境变量，如果不存在则返回默认值
func GetEnv(key, defaultValue string) string {
	return sharedConfig.GetEnv(key, defaultValue)
}

// MustGetEnv 获取环境变量，如果不存在则 panic
func MustGetEnv(key string) string {
	return sharedConfig.MustGetEnv(key)
}

// GetEnvInt 获取整数类型的环境变量
func GetEnvInt(key string, defaultValue int) int {
	return sharedConfig.GetEnvInt(key, defaultValue)
}

// GetEnvBool 获取布尔类型的环境变量
func GetEnvBool(key string, defaultValue bool) bool {
	return sharedConfig.GetEnvBool(key, defaultValue)
}

// ==================== 本地辅助函数 ====================

// GetJWTExpireHours 获取 JWT 过期时间（小时）
func GetJWTExpireHours() time.Duration {
	cfg := GetGlobalConfig()
	if cfg != nil {
		return cfg.JWT.ExpireHours
	}
	return 24
}

// GetJWTRefreshExpireHours 获取刷新 Token 过期时间（小时）
func GetJWTRefreshExpireHours() time.Duration {
	cfg := GetGlobalConfig()
	if cfg != nil {
		return cfg.JWT.RefreshExpireHours
	}
	return 168
}

// GetDBDSN 获取数据库 DSN
func GetDBDSN() string {
	cfg := GetGlobalConfig()
	if cfg != nil {
		return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=Local",
			cfg.Database.User,
			cfg.Database.Password,
			cfg.Database.Host,
			cfg.Database.Port,
			cfg.Database.DBName,
			cfg.Database.Charset)
	}
	return ""
}
