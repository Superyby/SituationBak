package service

import (
	"time"

	"SituationBak/internal/config"
	"SituationBak/internal/dto/request"
	"SituationBak/internal/dto/response"
	"SituationBak/internal/model"
	"SituationBak/shared/errors"
	"SituationBak/shared/utils"
	"SituationBak/internal/repository"
	"github.com/golang-jwt/jwt/v5"
)

// AuthService 认证服务
type AuthService struct {
	userRepo     *repository.UserRepository
	settingsRepo *repository.SettingsRepository
}

// NewAuthService 创建认证服务实例
func NewAuthService() *AuthService {
	return &AuthService{
		userRepo:     repository.NewUserRepository(),
		settingsRepo: repository.NewSettingsRepository(),
	}
}

// Claims JWT Claims结构
type Claims struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

// Register 用户注册
func (s *AuthService) Register(req *request.RegisterRequest) (*response.LoginResponse, error) {
	// 检查用户名是否已存在
	exists, err := s.userRepo.ExistsByUsername(req.Username)
	if err != nil {
		return nil, errors.ErrInternal(err)
	}
	if exists {
		return nil, errors.WithCode(errors.CodeUsernameExist)
	}

	// 检查邮箱是否已存在
	exists, err = s.userRepo.ExistsByEmail(req.Email)
	if err != nil {
		return nil, errors.ErrInternal(err)
	}
	if exists {
		return nil, errors.WithCode(errors.CodeEmailExists)
	}

	// 密码加密
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, errors.ErrInternal(err)
	}

	// 创建用户
	user := &model.User{
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: hashedPassword,
		Role:         model.RoleUser,
		IsActive:     true,
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, errors.ErrInternal(err)
	}

	// 创建默认设置
	_, _ = s.settingsRepo.GetOrCreate(user.ID)

	// 生成Token
	return s.generateLoginResponse(user)
}

// Login 用户登录
func (s *AuthService) Login(req *request.LoginRequest) (*response.LoginResponse, error) {
	// 根据用户名或邮箱查找用户
	user, err := s.userRepo.FindByUsernameOrEmail(req.Username)
	if err != nil {
		return nil, errors.ErrInternal(err)
	}
	if user == nil {
		return nil, errors.WithCode(errors.CodeLoginFailed)
	}

	// 验证密码
	if !utils.CheckPassword(req.Password, user.PasswordHash) {
		return nil, errors.WithCode(errors.CodePasswordWrong)
	}

	// 检查用户状态
	if !user.IsActive {
		return nil, errors.New(errors.CodeForbidden, "账号已被禁用")
	}

	// 更新最后登录时间
	_ = s.userRepo.UpdateLastLogin(user.ID)

	// 生成Token
	return s.generateLoginResponse(user)
}

// RefreshToken 刷新Token
func (s *AuthService) RefreshToken(refreshToken string) (*response.RefreshTokenResponse, error) {
	// 解析刷新Token
	claims, err := s.parseToken(refreshToken)
	if err != nil {
		return nil, errors.WithCode(errors.CodeTokenInvalid)
	}

	// 验证用户是否存在
	user, err := s.userRepo.FindByID(claims.UserID)
	if err != nil {
		return nil, errors.ErrInternal(err)
	}
	if user == nil || !user.IsActive {
		return nil, errors.WithCode(errors.CodeTokenInvalid)
	}

	// 生成新的Token
	accessToken, err := s.generateToken(user, false)
	if err != nil {
		return nil, errors.ErrInternal(err)
	}

	newRefreshToken, err := s.generateToken(user, true)
	if err != nil {
		return nil, errors.ErrInternal(err)
	}

	return &response.RefreshTokenResponse{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    int64(config.GlobalConfig.JWT.ExpireHours * 3600),
	}, nil
}

// GetCurrentUser 获取当前用户信息
func (s *AuthService) GetCurrentUser(userID uint) (*response.UserInfo, error) {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, errors.ErrInternal(err)
	}
	if user == nil {
		return nil, errors.ErrUserNotFound()
	}

	return &response.UserInfo{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		Role:      user.Role,
		AvatarURL: user.AvatarURL,
		CreatedAt: user.CreatedAt,
	}, nil
}

// generateLoginResponse 生成登录响应
func (s *AuthService) generateLoginResponse(user *model.User) (*response.LoginResponse, error) {
	accessToken, err := s.generateToken(user, false)
	if err != nil {
		return nil, errors.ErrInternal(err)
	}

	refreshToken, err := s.generateToken(user, true)
	if err != nil {
		return nil, errors.ErrInternal(err)
	}

	return &response.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    int64(config.GlobalConfig.JWT.ExpireHours * 3600),
		User: &response.UserInfo{
			ID:        user.ID,
			Username:  user.Username,
			Email:     user.Email,
			Role:      user.Role,
			AvatarURL: user.AvatarURL,
			CreatedAt: user.CreatedAt,
		},
	}, nil
}

// generateToken 生成JWT Token
func (s *AuthService) generateToken(user *model.User, isRefresh bool) (string, error) {
	var expireHours time.Duration
	if isRefresh {
		expireHours = config.GlobalConfig.JWT.RefreshExpireHours
	} else {
		expireHours = config.GlobalConfig.JWT.ExpireHours
	}

	claims := &Claims{
		UserID:   user.ID,
		Username: user.Username,
		Role:     user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expireHours * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "SituationBak",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(config.GlobalConfig.JWT.Secret))
}

// parseToken 解析JWT Token
func (s *AuthService) parseToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.GlobalConfig.JWT.Secret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.ErrTokenInvalid()
}

// ValidateToken 验证Token并返回Claims
func (s *AuthService) ValidateToken(tokenString string) (*Claims, error) {
	return s.parseToken(tokenString)
}
