// Package model 重导出 shared/model 中的模型定义
// 为了保持向后兼容性，internal/model 现在是 shared/model 的别名
package model

import (
	sharedModel "SituationBak/shared/model"
)

// ==================== 类型别名 ====================

type (
	// User 用户模型
	User = sharedModel.User
	// UserSettings 用户设置模型
	UserSettings = sharedModel.UserSettings
	// SettingsJSON 自定义JSON设置
	SettingsJSON = sharedModel.SettingsJSON
	// Favorite 收藏模型
	Favorite = sharedModel.Favorite
	// TLECache TLE缓存模型
	TLECache = sharedModel.TLECache
	// TLE TLE数据结构
	TLE = sharedModel.TLE
)

// ==================== 常量重导出 ====================

const (
	// 用户角色
	RoleUser  = sharedModel.RoleUser
	RoleAdmin = sharedModel.RoleAdmin

	// 主题
	ThemeDark  = sharedModel.ThemeDark
	ThemeLight = sharedModel.ThemeLight
	ThemeAuto  = sharedModel.ThemeAuto

	// 语言
	LangZhCN = sharedModel.LangZhCN
	LangEnUS = sharedModel.LangEnUS
)

// ==================== 函数重导出 ====================

// AllModels 返回所有模型列表（用于数据库迁移）
func AllModels() []interface{} {
	return sharedModel.AllModels()
}
