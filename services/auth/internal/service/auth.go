package service

import (
	"time"

	"SituationBak/services/auth/internal/repository"
	"SituationBak/shared/config"
	"SituationBak/shared/errors"
	"SituationBak/shared/model"
	"SituationBak/shared/utils"

	"github.com/golang-jwt/jwt/v5"
)

// AuthService 认证服务
type AuthService struct {
	userRepo     *repository.UserRepository
	settingsRepo *repository.SettingsRepository
	jwtSecret    string
	jwtExpire    time.Duration
	jwtRefresh   time.Duration
}

// NewAuthService 创建认证服务实例
func NewAuthService(
	userRepo *repository.UserRepository,
	settingsRepo *repository.SettingsRepository,
	cfg *config.JWTConfig,
) *AuthService {
	return &AuthService{
		userRepo:     userRepo,
		settingsRepo: settingsRepo,
		jwtSecret:    cfg.Secret,
		jwtExpire:    cfg.ExpireHours,
		jwtRefresh:   cfg.RefreshExpireHours,
	}
}

// Claims JWT Claims结构
type Claims struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

// RegisterRequest 注册请求
type RegisterRequest struct {
	Username string
	Email    string
	Password string
}

// LoginRequest 登录请求
type LoginRequest struct {
	Username string
	Password string
}

// LoginResponse 登录响应
type LoginResponse struct {
	AccessToken  string
	RefreshToken string
	TokenType    string
	ExpiresIn    int64
	User         *UserInfo
}

// UserInfo 用户信息
type UserInfo struct {
	ID        uint
	Username  string
	Email     string
	Role      string
	AvatarURL string
	CreatedAt time.Time
}

// TokenResponse Token响应
type TokenResponse struct {
	AccessToken  string
	RefreshToken string
	TokenType    string
	ExpiresIn    int64
}

// ValidateTokenResponse 验证Token响应
type ValidateTokenResponse struct {
	Valid    bool
	UserID   uint
	Username string
	Role     string
}

// Register 用户注册
func (s *AuthService) Register(req *RegisterRequest) (*LoginResponse, error) {
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
func (s *AuthService) Login(req *LoginRequest) (*LoginResponse, error) {
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
func (s *AuthService) RefreshToken(refreshToken string) (*TokenResponse, error) {
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

	return &TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    int64(s.jwtExpire * 3600),
	}, nil
}

// ValidateToken 验证Token并返回用户信息
func (s *AuthService) ValidateToken(tokenString string) (*ValidateTokenResponse, error) {
	claims, err := s.parseToken(tokenString)
	if err != nil {
		return &ValidateTokenResponse{Valid: false}, nil
	}

	return &ValidateTokenResponse{
		Valid:    true,
		UserID:   claims.UserID,
		Username: claims.Username,
		Role:     claims.Role,
	}, nil
}

// GetCurrentUser 获取当前用户信息
func (s *AuthService) GetCurrentUser(userID uint) (*UserInfo, error) {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, errors.ErrInternal(err)
	}
	if user == nil {
		return nil, errors.ErrUserNotFound()
	}

	return &UserInfo{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		Role:      user.Role,
		AvatarURL: user.AvatarURL,
		CreatedAt: user.CreatedAt,
	}, nil
}

// generateLoginResponse 生成登录响应
func (s *AuthService) generateLoginResponse(user *model.User) (*LoginResponse, error) {
	accessToken, err := s.generateToken(user, false)
	if err != nil {
		return nil, errors.ErrInternal(err)
	}

	refreshToken, err := s.generateToken(user, true)
	if err != nil {
		return nil, errors.ErrInternal(err)
	}

	return &LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    int64(s.jwtExpire * 3600),
		User: &UserInfo{
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
		expireHours = s.jwtRefresh
	} else {
		expireHours = s.jwtExpire
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
	return token.SignedString([]byte(s.jwtSecret))
}

// parseToken 解析JWT Token
func (s *AuthService) parseToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.jwtSecret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.ErrTokenInvalid()
}
