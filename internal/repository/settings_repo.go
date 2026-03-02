package repository

import (
	"errors"

	"SituationBak/internal/model"
	"gorm.io/gorm"
)

// SettingsRepository 用户设置数据访问�?
type SettingsRepository struct {
	db *gorm.DB
}

// NewSettingsRepository 创建设置Repository实例
func NewSettingsRepository() *SettingsRepository {
	return &SettingsRepository{db: DB}
}

// Create 创建用户设置
func (r *SettingsRepository) Create(settings *model.UserSettings) error {
	return r.db.Create(settings).Error
}

// FindByUserID 根据用户ID查找设置
func (r *SettingsRepository) FindByUserID(userID uint) (*model.UserSettings, error) {
	var settings model.UserSettings
	err := r.db.Where("user_id = ?", userID).First(&settings).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &settings, err
}

// Update 更新用户设置
func (r *SettingsRepository) Update(settings *model.UserSettings) error {
	return r.db.Save(settings).Error
}

// UpdateFields 更新指定字段
func (r *SettingsRepository) UpdateFields(userID uint, fields map[string]interface{}) error {
	return r.db.Model(&model.UserSettings{}).Where("user_id = ?", userID).Updates(fields).Error
}

// Delete 删除用户设置
func (r *SettingsRepository) Delete(userID uint) error {
	return r.db.Where("user_id = ?", userID).Delete(&model.UserSettings{}).Error
}

// GetOrCreate 获取或创建默认设�?
func (r *SettingsRepository) GetOrCreate(userID uint) (*model.UserSettings, error) {
	settings, err := r.FindByUserID(userID)
	if err != nil {
		return nil, err
	}

	if settings == nil {
		// 创建默认设置
		settings = &model.UserSettings{
			UserID:         userID,
			SatelliteLimit: 5000,
			ShowDebris:     false,
			Theme:          model.ThemeDark,
			Language:       model.LangZhCN,
		}
		if err := r.Create(settings); err != nil {
			return nil, err
		}
	}

	return settings, nil
}

// UpdateTheme 更新主题设置
func (r *SettingsRepository) UpdateTheme(userID uint, theme string) error {
	return r.db.Model(&model.UserSettings{}).Where("user_id = ?", userID).Update("theme", theme).Error
}

// UpdateLanguage 更新语言设置
func (r *SettingsRepository) UpdateLanguage(userID uint, language string) error {
	return r.db.Model(&model.UserSettings{}).Where("user_id = ?", userID).Update("language", language).Error
}

// UpdateSatelliteLimit 更新卫星数量限制
func (r *SettingsRepository) UpdateSatelliteLimit(userID uint, limit int) error {
	return r.db.Model(&model.UserSettings{}).Where("user_id = ?", userID).Update("satellite_limit", limit).Error
}

// UpdateShowDebris 更新是否显示碎片
func (r *SettingsRepository) UpdateShowDebris(userID uint, show bool) error {
	return r.db.Model(&model.UserSettings{}).Where("user_id = ?", userID).Update("show_debris", show).Error
}
