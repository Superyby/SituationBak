package config

import (
	"fmt"
	"os"
	"time"

	"SituationBak/shared/constants"
)

// GlobalConfig 全局配置实例
var GlobalConfig *Config

// Config 应用配置结构体
type Config struct {
	App        AppConfig        `mapstructure:"app"`
	Database   DatabaseConfig   `mapstructure:"database"`
	Redis      RedisConfig      `mapstructure:"redis"`
	ClickHouse ClickHouseConfig `mapstructure:"clickhouse"`
	JWT        JWTConfig        `mapstructure:"jwt"`
	Log        LogConfig        `mapstructure:"log"`
	RateLimit  RateLimitConfig  `mapstructure:"ratelimit"`
	External   ExternalConfig   `mapstructure:"external"`
	CORS       CORSConfig       `mapstructure:"cors"`
	GRPC       GRPCConfig       `mapstructure:"grpc"`
}

// AppConfig 应用配置
type AppConfig struct {
	Name  string `mapstructure:"name"`
	Env   string `mapstructure:"env"`
	Port  int    `mapstructure:"port"`
	Debug bool   `mapstructure:"debug"`
}

// DatabaseConfig MySQL数据库配置
type DatabaseConfig struct {
	Host            string `mapstructure:"host"`
	Port            int    `mapstructure:"port"`
	User            string `mapstructure:"user"`
	Password        string `mapstructure:"password"`
	DBName          string `mapstructure:"dbname"`
	Charset         string `mapstructure:"charset"`
	MaxOpenConns    int    `mapstructure:"max_open_conns"`
	MaxIdleConns    int    `mapstructure:"max_idle_conns"`
	ConnMaxLifetime int    `mapstructure:"conn_max_lifetime"`
}

// RedisConfig Redis配置
type RedisConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
	PoolSize int    `mapstructure:"pool_size"`
}

// ClickHouseConfig ClickHouse配置
type ClickHouseConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Database string `mapstructure:"database"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
}

// JWTConfig JWT配置
type JWTConfig struct {
	Secret             string        `mapstructure:"secret"`
	ExpireHours        time.Duration `mapstructure:"expire_hours"`
	RefreshExpireHours time.Duration `mapstructure:"refresh_expire_hours"`
}

// LogConfig 日志配置
type LogConfig struct {
	Level    string `mapstructure:"level"`
	Format   string `mapstructure:"format"`
	Output   string `mapstructure:"output"`
	FilePath string `mapstructure:"file_path"`
}

// RateLimitConfig 限流配置
type RateLimitConfig struct {
	RequestsPerSecond int `mapstructure:"requests_per_second"`
	Burst             int `mapstructure:"burst"`
}

// ExternalConfig 第三方API配置
type ExternalConfig struct {
	KeepTrack  KeepTrackConfig  `mapstructure:"keeptrack"`
	SpaceTrack SpaceTrackConfig `mapstructure:"spacetrack"`
}

// KeepTrackConfig KeepTrack API配置
type KeepTrackConfig struct {
	BaseURL string `mapstructure:"base_url"`
	Timeout int    `mapstructure:"timeout"`
}

// SpaceTrackConfig SpaceTrack API配置
type SpaceTrackConfig struct {
	BaseURL  string `mapstructure:"base_url"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	Timeout  int    `mapstructure:"timeout"`
}

// CORSConfig CORS配置
type CORSConfig struct {
	AllowedOrigins   []string `mapstructure:"allowed_origins"`
	AllowedMethods   []string `mapstructure:"allowed_methods"`
	AllowedHeaders   []string `mapstructure:"allowed_headers"`
	ExposeHeaders    []string `mapstructure:"expose_headers"`
	AllowCredentials bool     `mapstructure:"allow_credentials"`
	MaxAge           int      `mapstructure:"max_age"`
}

// GRPCConfig gRPC服务地址配置
type GRPCConfig struct {
	AuthAddr      string `mapstructure:"auth_addr"`
	UserAddr      string `mapstructure:"user_addr"`
	SatelliteAddr string `mapstructure:"satellite_addr"`
	FavoriteAddr  string `mapstructure:"favorite_addr"`
}

// ==================== 辅助方法 ====================

// DSN 返回 MySQL 连接字符串
func (c *DatabaseConfig) DSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=Local",
		c.User, c.Password, c.Host, c.Port, c.DBName, c.Charset)
}

// ==================== 环境判断函数 ====================

// IsDevelopment 判断是否为开发环境
func IsDevelopment() bool {
	if GlobalConfig != nil {
		return GlobalConfig.App.Env == constants.EnvDevelopment
	}
	return GetEnv("APP_ENV", constants.EnvDevelopment) == constants.EnvDevelopment
}

// IsProduction 判断是否为生产环境
func IsProduction() bool {
	if GlobalConfig != nil {
		return GlobalConfig.App.Env == constants.EnvProduction
	}
	return GetEnv("APP_ENV", "") == constants.EnvProduction
}

// IsTesting 判断是否为测试环境
func IsTesting() bool {
	if GlobalConfig != nil {
		return GlobalConfig.App.Env == constants.EnvTesting
	}
	return GetEnv("APP_ENV", "") == constants.EnvTesting
}

// ==================== 环境变量辅助函数 ====================

// GetEnv 获取环境变量，如果不存在则返回默认值
func GetEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// MustGetEnv 获取环境变量，如果不存在则 panic
func MustGetEnv(key string) string {
	value := os.Getenv(key)
	if value == "" {
		panic(fmt.Sprintf("environment variable %s is required", key))
	}
	return value
}

// GetEnvInt 获取整数类型的环境变量
func GetEnvInt(key string, defaultValue int) int {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	var result int
	_, err := fmt.Sscanf(value, "%d", &result)
	if err != nil {
		return defaultValue
	}
	return result
}

// GetEnvBool 获取布尔类型的环境变量
func GetEnvBool(key string, defaultValue bool) bool {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value == "true" || value == "1" || value == "yes"
}

// ==================== 配置获取函数 ====================

// GetConfig 获取全局配置
func GetConfig() *Config {
	return GlobalConfig
}

// SetConfig 设置全局配置
func SetConfig(cfg *Config) {
	GlobalConfig = cfg
}
