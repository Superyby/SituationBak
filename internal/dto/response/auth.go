package response

import "time"

// UserInfo 用户基本信息
type UserInfo struct {
	ID        uint      `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Role      string    `json:"role"`
	AvatarURL string    `json:"avatar_url,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

// LoginResponse 登录响应
type LoginResponse struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	TokenType    string    `json:"token_type"`
	ExpiresIn    int64     `json:"expires_in"` // 秒
	User         *UserInfo `json:"user"`
}

// RefreshTokenResponse 刷新Token响应
type RefreshTokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int64  `json:"expires_in"`
}

// UserProfileResponse 用户资料响应
type UserProfileResponse struct {
	ID          uint       `json:"id"`
	Username    string     `json:"username"`
	Email       string     `json:"email"`
	Role        string     `json:"role"`
	AvatarURL   string     `json:"avatar_url,omitempty"`
	IsActive    bool       `json:"is_active"`
	LastLoginAt *time.Time `json:"last_login_at,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

// UserSettingsResponse 用户设置响应
type UserSettingsResponse struct {
	SatelliteLimit int                    `json:"satellite_limit"`
	ShowDebris     bool                   `json:"show_debris"`
	Theme          string                 `json:"theme"`
	Language       string                 `json:"language"`
	SettingsJSON   map[string]interface{} `json:"settings_json,omitempty"`
}
