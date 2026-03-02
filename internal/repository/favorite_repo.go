package repository

import (
	"errors"

	"SituationBak/internal/model"
	"gorm.io/gorm"
)

// FavoriteRepository 收藏数据访问层
type FavoriteRepository struct {
	db *gorm.DB
}

// NewFavoriteRepository 创建收藏Repository实例
func NewFavoriteRepository() *FavoriteRepository {
	return &FavoriteRepository{db: DB}
}

// Create 创建收藏
func (r *FavoriteRepository) Create(favorite *model.Favorite) error {
	return r.db.Create(favorite).Error
}

// FindByID 根据ID查找收藏
func (r *FavoriteRepository) FindByID(id uint) (*model.Favorite, error) {
	var favorite model.Favorite
	err := r.db.First(&favorite, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &favorite, err
}

// FindByUserID 根据用户ID查找所有收藏
func (r *FavoriteRepository) FindByUserID(userID uint, page, pageSize int) ([]model.Favorite, int64, error) {
	var favorites []model.Favorite
	var total int64

	r.db.Model(&model.Favorite{}).Where("user_id = ?", userID).Count(&total)

	offset := (page - 1) * pageSize
	err := r.db.Where("user_id = ?", userID).
		Order("created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&favorites).Error

	return favorites, total, err
}

// FindByUserIDAndNoradID 根据用户ID和NORAD ID查找收藏
func (r *FavoriteRepository) FindByUserIDAndNoradID(userID uint, noradID int) (*model.Favorite, error) {
	var favorite model.Favorite
	err := r.db.Where("user_id = ? AND norad_id = ?", userID, noradID).First(&favorite).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &favorite, err
}

// Delete 删除收藏
func (r *FavoriteRepository) Delete(id uint) error {
	return r.db.Delete(&model.Favorite{}, id).Error
}

// DeleteByUserIDAndNoradID 根据用户ID和NORAD ID删除收藏
func (r *FavoriteRepository) DeleteByUserIDAndNoradID(userID uint, noradID int) error {
	return r.db.Where("user_id = ? AND norad_id = ?", userID, noradID).Delete(&model.Favorite{}).Error
}

// Exists 检查是否已收藏
func (r *FavoriteRepository) Exists(userID uint, noradID int) (bool, error) {
	var count int64
	err := r.db.Model(&model.Favorite{}).
		Where("user_id = ? AND norad_id = ?", userID, noradID).
		Count(&count).Error
	return count > 0, err
}

// CountByUserID 获取用户收藏数量
func (r *FavoriteRepository) CountByUserID(userID uint) (int64, error) {
	var count int64
	err := r.db.Model(&model.Favorite{}).Where("user_id = ?", userID).Count(&count).Error
	return count, err
}

// GetAllByUserID 获取用户所有收藏（不分页）
func (r *FavoriteRepository) GetAllByUserID(userID uint) ([]model.Favorite, error) {
	var favorites []model.Favorite
	err := r.db.Where("user_id = ?", userID).Order("created_at DESC").Find(&favorites).Error
	return favorites, err
}

// Update 更新收藏
func (r *FavoriteRepository) Update(favorite *model.Favorite) error {
	return r.db.Save(favorite).Error
}

// UpdateNotes 更新收藏备注
func (r *FavoriteRepository) UpdateNotes(id uint, notes string) error {
	return r.db.Model(&model.Favorite{}).Where("id = ?", id).Update("notes", notes).Error
}
