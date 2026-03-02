package model

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"

	"gorm.io/gorm"
)

// ==================== 用户模型 ====================

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

// ==================== 用户设置模型 ====================

// UserSettings 用户设置模型
type UserSettings struct {
	ID             uint           `gorm:"primaryKey" json:"id"`
	UserID         uint           `gorm:"uniqueIndex;not null" json:"user_id"`
	SatelliteLimit int            `gorm:"default:5000" json:"satellite_limit"`
	ShowDebris     bool           `gorm:"default:false" json:"show_debris"`
	Theme          string         `gorm:"size:20;default:dark" json:"theme"`
	Language       string         `gorm:"size:10;default:zh-CN" json:"language"`
	SettingsJSON   SettingsJSON   `gorm:"type:json" json:"settings_json,omitempty"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`

	// 关联
	User *User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

// TableName 指定表名
func (UserSettings) TableName() string {
	return "user_settings"
}

// SettingsJSON 自定义JSON设置
type SettingsJSON map[string]interface{}

// Value 实现 driver.Valuer 接口
func (s SettingsJSON) Value() (driver.Value, error) {
	if s == nil {
		return nil, nil
	}
	return json.Marshal(s)
}

// Scan 实现 sql.Scanner 接口
func (s *SettingsJSON) Scan(value interface{}) error {
	if value == nil {
		*s = nil
		return nil
	}

	var bytes []byte
	switch v := value.(type) {
	case []byte:
		bytes = v
	case string:
		bytes = []byte(v)
	default:
		return errors.New("无法扫描 SettingsJSON 类型")
	}

	return json.Unmarshal(bytes, s)
}

// BeforeUpdate 更新前钩子
func (u *UserSettings) BeforeUpdate(tx *gorm.DB) error {
	u.UpdatedAt = time.Now()
	return nil
}

// 主题常量
const (
	ThemeDark  = "dark"
	ThemeLight = "light"
	ThemeAuto  = "auto"
)

// 语言常量
const (
	LangZhCN = "zh-CN"
	LangEnUS = "en-US"
)

// ==================== 收藏模型 ====================

// Favorite 收藏模型
type Favorite struct {
	ID            uint           `gorm:"primaryKey" json:"id"`
	UserID        uint           `gorm:"not null;index" json:"user_id"`
	NoradID       int            `gorm:"not null" json:"norad_id"`
	SatelliteID   string         `gorm:"size:20" json:"satellite_id,omitempty"` // 可选的字符串ID
	SatelliteName string         `gorm:"size:100" json:"satellite_name"`
	Name          string         `gorm:"size:100" json:"name,omitempty"` // 别名，用于兼容
	Notes         string         `gorm:"type:text" json:"notes,omitempty"`
	CreatedAt     time.Time      `json:"created_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`

	// 关联
	User *User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

// TableName 指定表名
func (Favorite) TableName() string {
	return "favorites"
}

// BeforeCreate 创建前钩子
func (f *Favorite) BeforeCreate(tx *gorm.DB) error {
	f.CreatedAt = time.Now()
	return nil
}

// ==================== TLE 缓存模型 ====================

// TLECache TLE缓存模型
type TLECache struct {
	NoradID   int       `gorm:"primaryKey" json:"norad_id"`
	Name      string    `gorm:"size:100" json:"name"`
	TLELine1  string    `gorm:"size:70" json:"tle_line1"`
	TLELine2  string    `gorm:"size:70" json:"tle_line2"`
	Epoch     time.Time `gorm:"index" json:"epoch"`
	UpdatedAt time.Time `json:"updated_at"`
}

// TableName 指定表名
func (TLECache) TableName() string {
	return "tle_cache"
}

// TLE 返回完整的TLE数据
type TLE struct {
	NoradID int    `json:"norad_id"`
	Name    string `json:"name"`
	Line1   string `json:"line1"`
	Line2   string `json:"line2"`
}

// ToTLE 转换为TLE结构
func (t *TLECache) ToTLE() *TLE {
	return &TLE{
		NoradID: t.NoradID,
		Name:    t.Name,
		Line1:   t.TLELine1,
		Line2:   t.TLELine2,
	}
}

// IsExpired 判断TLE数据是否过期（超过24小时）
func (t *TLECache) IsExpired() bool {
	return time.Since(t.UpdatedAt) > 24*time.Hour
}

// ==================== 模型列表（用于自动迁移）====================

// AllModels 返回所有模型列表（用于数据库迁移）
func AllModels() []interface{} {
	return []interface{}{
		&User{},
		&UserSettings{},
		&Favorite{},
		&TLECache{},
	}
}
