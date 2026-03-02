package service

import (
	"SituationBak/internal/dto/request"
	"SituationBak/internal/dto/response"
	"SituationBak/internal/pkg/errors"
	"SituationBak/internal/pkg/utils"
	"SituationBak/internal/repository"
)

// UserService 用户服务
type UserService struct {
	userRepo     *repository.UserRepository
	settingsRepo *repository.SettingsRepository
}

// NewUserService 创建用户服务实例
func NewUserService() *UserService {
	return &UserService{
		userRepo:     repository.NewUserRepository(),
		settingsRepo: repository.NewSettingsRepository(),
	}
}

// GetProfile 获取用户资料
func (s *UserService) GetProfile(userID uint) (*response.UserProfileResponse, error) {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, errors.ErrInternal(err)
	}
	if user == nil {
		return nil, errors.ErrUserNotFound()
	}

	return &response.UserProfileResponse{
		ID:          user.ID,
		Username:    user.Username,
		Email:       user.Email,
		Role:        user.Role,
		AvatarURL:   user.AvatarURL,
		IsActive:    user.IsActive,
		LastLoginAt: user.LastLoginAt,
		CreatedAt:   user.CreatedAt,
		UpdatedAt:   user.UpdatedAt,
	}, nil
}

// UpdateProfile 更新用户资料
func (s *UserService) UpdateProfile(userID uint, req *request.UpdateProfileRequest) (*response.UserProfileResponse, error) {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, errors.ErrInternal(err)
	}
	if user == nil {
		return nil, errors.ErrUserNotFound()
	}

	// 检查用户名是否被占用
	if req.Username != "" && req.Username != user.Username {
		exists, err := s.userRepo.ExistsByUsername(req.Username)
		if err != nil {
			return nil, errors.ErrInternal(err)
		}
		if exists {
			return nil, errors.WithCode(errors.CodeUsernameExist)
		}
		user.Username = req.Username
	}

	// 检查邮箱是否被占用
	if req.Email != "" && req.Email != user.Email {
		exists, err := s.userRepo.ExistsByEmail(req.Email)
		if err != nil {
			return nil, errors.ErrInternal(err)
		}
		if exists {
			return nil, errors.WithCode(errors.CodeEmailExists)
		}
		user.Email = req.Email
	}

	// 更新头像
	if req.AvatarURL != "" {
		user.AvatarURL = req.AvatarURL
	}

	if err := s.userRepo.Update(user); err != nil {
		return nil, errors.ErrInternal(err)
	}

	return &response.UserProfileResponse{
		ID:          user.ID,
		Username:    user.Username,
		Email:       user.Email,
		Role:        user.Role,
		AvatarURL:   user.AvatarURL,
		IsActive:    user.IsActive,
		LastLoginAt: user.LastLoginAt,
		CreatedAt:   user.CreatedAt,
		UpdatedAt:   user.UpdatedAt,
	}, nil
}

// ChangePassword 修改密码
func (s *UserService) ChangePassword(userID uint, req *request.ChangePasswordRequest) error {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return errors.ErrInternal(err)
	}
	if user == nil {
		return errors.ErrUserNotFound()
	}

	// 验证旧密码
	if !utils.CheckPassword(req.OldPassword, user.PasswordHash) {
		return errors.ErrPasswordWrong()
	}

	// 加密新密码
	hashedPassword, err := utils.HashPassword(req.NewPassword)
	if err != nil {
		return errors.ErrInternal(err)
	}

	user.PasswordHash = hashedPassword
	if err := s.userRepo.Update(user); err != nil {
		return errors.ErrInternal(err)
	}

	return nil
}

// GetSettings 获取用户设置
func (s *UserService) GetSettings(userID uint) (*response.UserSettingsResponse, error) {
	settings, err := s.settingsRepo.GetOrCreate(userID)
	if err != nil {
		return nil, errors.ErrInternal(err)
	}

	return &response.UserSettingsResponse{
		SatelliteLimit: settings.SatelliteLimit,
		ShowDebris:     settings.ShowDebris,
		Theme:          settings.Theme,
		Language:       settings.Language,
		SettingsJSON:   settings.SettingsJSON,
	}, nil
}

// UpdateSettings 更新用户设置
func (s *UserService) UpdateSettings(userID uint, req *request.UpdateSettingsRequest) (*response.UserSettingsResponse, error) {
	settings, err := s.settingsRepo.GetOrCreate(userID)
	if err != nil {
		return nil, errors.ErrInternal(err)
	}

	// 更新字段
	if req.SatelliteLimit != nil {
		settings.SatelliteLimit = *req.SatelliteLimit
	}
	if req.ShowDebris != nil {
		settings.ShowDebris = *req.ShowDebris
	}
	if req.Theme != "" {
		settings.Theme = req.Theme
	}
	if req.Language != "" {
		settings.Language = req.Language
	}

	if err := s.settingsRepo.Update(settings); err != nil {
		return nil, errors.ErrInternal(err)
	}

	return &response.UserSettingsResponse{
		SatelliteLimit: settings.SatelliteLimit,
		ShowDebris:     settings.ShowDebris,
		Theme:          settings.Theme,
		Language:       settings.Language,
		SettingsJSON:   settings.SettingsJSON,
	}, nil
}
