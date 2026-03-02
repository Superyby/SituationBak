package model

import (
	"time"
)

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

// TLE 返回完整�?TLE 数据
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

// IsExpired 判断TLE数据是否过期（超�?4小时�?
func (t *TLECache) IsExpired() bool {
	return time.Since(t.UpdatedAt) > 24*time.Hour
}
