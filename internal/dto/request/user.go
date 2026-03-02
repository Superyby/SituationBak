package request

// UpdateProfileRequest 更新用户资料请求
type UpdateProfileRequest struct {
	Username  string `json:"username,omitempty" validate:"omitempty,min=3,max=50"`
	Email     string `json:"email,omitempty" validate:"omitempty,email"`
	AvatarURL string `json:"avatar_url,omitempty" validate:"omitempty,url,max=255"`
}

// UpdateSettingsRequest 更新用户设置请求
type UpdateSettingsRequest struct {
	SatelliteLimit *int   `json:"satellite_limit,omitempty" validate:"omitempty,min=100,max=50000"`
	ShowDebris     *bool  `json:"show_debris,omitempty"`
	Theme          string `json:"theme,omitempty" validate:"omitempty,oneof=dark light auto"`
	Language       string `json:"language,omitempty" validate:"omitempty,oneof=zh-CN en-US"`
}
