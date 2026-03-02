package config

import (
	"fmt"
	"os"
	"strings"

	sharedConfig "SituationBak/shared/config"

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

	// 设置全局配置（同时设置 internal 和 shared）
	GlobalConfig = &config
	sharedConfig.GlobalConfig = &config

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

// LoadConfigWithEnv 加载配置并应用环境变量覆盖
// 这是 LoadConfig 的别名，保持向后兼容
func LoadConfigWithEnv(configPath string) (*Config, error) {
	return LoadConfig(configPath)
}

// LoadConfigFromEnv 仅从环境变量加载配置（用于容器化部署）
func LoadConfigFromEnv() *Config {
	return &Config{
		App: AppConfig{
			Name:  GetEnv("APP_NAME", "SituationBak"),
			Env:   GetEnv("APP_ENV", "development"),
			Port:  GetEnvInt("APP_PORT", 4000),
			Debug: GetEnvBool("APP_DEBUG", false),
		},
		Database: DatabaseConfig{
			Host:            GetEnv("DB_HOST", "localhost"),
			Port:            GetEnvInt("DB_PORT", 3306),
			User:            GetEnv("DB_USER", "root"),
			Password:        GetEnv("DB_PASSWORD", ""),
			DBName:          GetEnv("DB_NAME", "situationbak"),
			Charset:         GetEnv("DB_CHARSET", "utf8mb4"),
			MaxOpenConns:    GetEnvInt("DB_MAX_OPEN_CONNS", 100),
			MaxIdleConns:    GetEnvInt("DB_MAX_IDLE_CONNS", 10),
			ConnMaxLifetime: GetEnvInt("DB_CONN_MAX_LIFETIME", 3600),
		},
		Redis: RedisConfig{
			Host:     GetEnv("REDIS_HOST", "localhost"),
			Port:     GetEnvInt("REDIS_PORT", 6379),
			Password: GetEnv("REDIS_PASSWORD", ""),
			DB:       GetEnvInt("REDIS_DB", 0),
			PoolSize: GetEnvInt("REDIS_POOL_SIZE", 10),
		},
		JWT: JWTConfig{
			Secret: GetEnv("JWT_SECRET", "your-secret-key"),
		},
	}
}

// 移除 os 包的显式使用，因为已在 config.go 中通过 shared/config 重导出
var _ = os.Getenv // 保持 import 不报错
