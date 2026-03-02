package model

import (
	"time"

	"gorm.io/gorm"
)

// User 用户模型
type User struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	Username     string         `gorm:"size:50;uniqueIndex;not null" json:"username"`
	Email        string         `gorm:"size:100;uniqueIndex;not null" json:"email"`
	PasswordHash string         `gorm:"size:255;not null" json:"-"`
	Role         string         `gorm:"size:20;default:user" json:"role"`
	AvatarURL    string         `gorm:"size:255" json:"avatar_url,omitempty"`
	IsActive     bool           `gorm:"default:true" json:"is_active"`
	LastLoginAt  *time.Time     `json:"last_login_at,omitempty"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`

	// 关联
	Settings  *UserSettings `gorm:"foreignKey:UserID" json:"settings,omitempty"`
	Favorites []Favorite    `gorm:"foreignKey:UserID" json:"favorites,omitempty"`
}

// TableName 指定表名
func (User) TableName() string {
	return "users"
}

// UserRole 用户角色常量
const (
	RoleUser  = "user"
	RoleAdmin = "admin"
)

// IsAdmin 判断是否为管理员
func (u *User) IsAdmin() bool {
	return u.Role == RoleAdmin
}

// BeforeCreate 创建前钩子
func (u *User) BeforeCreate(tx *gorm.DB) error {
	u.CreatedAt = time.Now()
	u.UpdatedAt = time.Now()
	return nil
}

// BeforeUpdate 更新前钩子
func (u *User) BeforeUpdate(tx *gorm.DB) error {
	u.UpdatedAt = time.Now()
	return nil
}

// UserSettings 用户设置模型
type UserSettings struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `gorm:"uniqueIndex;not null" json:"user_id"`
	Theme     string    `gorm:"size:20;default:dark" json:"theme"`
	Language  string    `gorm:"size:10;default:zh-CN" json:"language"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// TableName 指定表名
func (UserSettings) TableName() string {
	return "user_settings"
}

// Favorite 收藏模型
type Favorite struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	UserID      uint      `gorm:"index;not null" json:"user_id"`
	SatelliteID string    `gorm:"size:20;not null" json:"satellite_id"`
	Name        string    `gorm:"size:100;not null" json:"name"`
	CreatedAt   time.Time `json:"created_at"`
}

// TableName 指定表名
func (Favorite) TableName() string {
	return "favorites"
}
