package model

import (
	"time"

	"gorm.io/gorm"
)

// Favorite 收藏模型
type Favorite struct {
	ID            uint           `gorm:"primaryKey" json:"id"`
	UserID        uint           `gorm:"not null;index" json:"user_id"`
	NoradID       int            `gorm:"not null" json:"norad_id"`
	SatelliteName string         `gorm:"size:100" json:"satellite_name"`
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

// BeforeCreate 创建前钩�?
func (f *Favorite) BeforeCreate(tx *gorm.DB) error {
	f.CreatedAt = time.Now()
	return nil
}
