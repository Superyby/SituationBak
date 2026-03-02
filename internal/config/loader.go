package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/viper"
)

// LoadConfig 加载配置文件
func LoadConfig(configPath string) (*Config, error) {
	v := viper.New()

	// 设置配置文件路径
	if configPath != "" {
		v.SetConfigFile(configPath)
	} else {
		// 默认配置文件路径
		v.SetConfigName("config")
		v.SetConfigType("yaml")
		v.AddConfigPath("./configs")
		v.AddConfigPath(".")
	}

	// 读取配置文件
	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("读取配置文件失败: %w", err)
	}

	// 设置环境变量前缀
	v.SetEnvPrefix("ORBITAL")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	// 绑定环境变量覆盖
	bindEnvVariables(v)

	// 解析配置到结构体
	var config Config
	if err := v.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("解析配置失败: %w", err)
	}

	// 设置全局配置
	GlobalConfig = &config

	return &config, nil
}

// bindEnvVariables 绑定环境变量
func bindEnvVariables(v *viper.Viper) {
	// 应用配置
	_ = v.BindEnv("app.env", "APP_ENV")
	_ = v.BindEnv("app.port", "APP_PORT")

	// 数据库配置
	_ = v.BindEnv("database.host", "DB_HOST")
	_ = v.BindEnv("database.port", "DB_PORT")
	_ = v.BindEnv("database.user", "DB_USER")
	_ = v.BindEnv("database.password", "DB_PASSWORD")
	_ = v.BindEnv("database.dbname", "DB_NAME")

	// Redis配置
	_ = v.BindEnv("redis.host", "REDIS_HOST")
	_ = v.BindEnv("redis.port", "REDIS_PORT")
	_ = v.BindEnv("redis.password", "REDIS_PASSWORD")

	// ClickHouse配置
	_ = v.BindEnv("clickhouse.host", "CLICKHOUSE_HOST")
	_ = v.BindEnv("clickhouse.port", "CLICKHOUSE_PORT")
	_ = v.BindEnv("clickhouse.username", "CLICKHOUSE_USERNAME")
	_ = v.BindEnv("clickhouse.password", "CLICKHOUSE_PASSWORD")

	// JWT配置
	_ = v.BindEnv("jwt.secret", "JWT_SECRET")

	// Space-Track配置
	_ = v.BindEnv("external.spacetrack.username", "SPACETRACK_USERNAME")
	_ = v.BindEnv("external.spacetrack.password", "SPACETRACK_PASSWORD")
}

// GetEnv 获取环境变量，如果不存在则返回默认值
func GetEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// IsDevelopment 判断是否为开发环境
func IsDevelopment() bool {
	return GlobalConfig != nil && GlobalConfig.App.Env == "development"
}

// IsProduction 判断是否为生产环境
func IsProduction() bool {
	return GlobalConfig != nil && GlobalConfig.App.Env == "production"
}
