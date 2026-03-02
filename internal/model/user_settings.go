package model

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"

	"gorm.io/gorm"
)

// UserSettings 用户设置模型
type UserSettings struct {
	ID             uint           `gorm:"primaryKey" json:"id"`
	UserID         uint           `gorm:"uniqueIndex;not null" json:"user_id"`
	SatelliteLimit int            `gorm:"default:5000" json:"satellite_limit"`
	ShowDebris     bool           `gorm:"default:false" json:"show_debris"`
	Theme          string         `gorm:"size:20;default:dark" json:"theme"`
	Language       string         `gorm:"size:10;default:zh-CN" json:"language"`
	SettingsJSON   SettingsJSON   `gorm:"type:json" json:"settings_json,omitempty"`
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
