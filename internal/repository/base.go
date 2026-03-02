package repository

import (
	"context"
	"fmt"
	"time"

	"SituationBak/internal/config"
	"SituationBak/internal/model"
	"SituationBak/shared/database"

	"github.com/redis/go-redis/v9"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// ==================== 全局实例 ====================

var (
	// DB 全局数据库实例
	DB *gorm.DB
	// Redis 全局Redis客户端
	Redis *redis.Client
)

// ==================== 数据库初始化 ====================

// InitMySQL 初始化MySQL连接
func InitMySQL(cfg *config.DatabaseConfig) error {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=Local",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DBName, cfg.Charset)

	var logLevel logger.LogLevel
	if config.IsDevelopment() {
		logLevel = logger.Info
	} else {
		logLevel = logger.Silent
	}

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logLevel),
	})
	if err != nil {
		return fmt.Errorf("连接MySQL失败: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("获取数据库实例失败: %w", err)
	}

	// 设置连接池
	sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
	sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
	sqlDB.SetConnMaxLifetime(time.Duration(cfg.ConnMaxLifetime) * time.Second)

	// 设置全局实例
	DB = db
	database.DB = db

	return nil
}

// InitRedis 初始化Redis连接
func InitRedis(cfg *config.RedisConfig) error {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Password: cfg.Password,
		DB:       cfg.DB,
		PoolSize: cfg.PoolSize,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return fmt.Errorf("连接Redis失败: %w", err)
	}

	// 设置全局实例
	Redis = client
	database.RedisClient = client

	return nil
}

// ==================== 数据库迁移 ====================

// AutoMigrate 自动迁移数据库表
func AutoMigrate() error {
	return DB.AutoMigrate(model.AllModels()...)
}

// MigrateModel 迁移指定模型
func MigrateModel(models ...interface{}) error {
	return DB.AutoMigrate(models...)
}

// ==================== 连接关闭 ====================

// Close 关闭数据库连接
func Close() error {
	if Redis != nil {
		if err := Redis.Close(); err != nil {
			return err
		}
	}

	if DB != nil {
		sqlDB, err := DB.DB()
		if err != nil {
			return err
		}
		return sqlDB.Close()
	}

	return nil
}

// CloseMySQL 仅关闭MySQL连接
func CloseMySQL() error {
	if DB != nil {
		sqlDB, err := DB.DB()
		if err != nil {
			return err
		}
		return sqlDB.Close()
	}
	return nil
}

// CloseRedis 仅关闭Redis连接
func CloseRedis() error {
	if Redis != nil {
		return Redis.Close()
	}
	return nil
}

// ==================== 事务处理 ====================

// Transaction 事务处理
func Transaction(fn func(tx *gorm.DB) error) error {
	return DB.Transaction(fn)
}

// TransactionWithContext 带上下文的事务处理
func TransactionWithContext(ctx context.Context, fn func(tx *gorm.DB) error) error {
	return DB.WithContext(ctx).Transaction(fn)
}

// ==================== 数据库状态检查 ====================

// IsDBConnected 检查数据库是否已连接
func IsDBConnected() bool {
	if DB == nil {
		return false
	}
	sqlDB, err := DB.DB()
	if err != nil {
		return false
	}
	return sqlDB.Ping() == nil
}

// IsRedisConnected 检查Redis是否已连接
func IsRedisConnected() bool {
	if Redis == nil {
		return false
	}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	return Redis.Ping(ctx).Err() == nil
}

// GetDB 获取数据库实例
func GetDB() *gorm.DB {
	return DB
}

// GetRedis 获取Redis实例
func GetRedis() *redis.Client {
	return Redis
}
