package response

import "time"

// SatelliteInfo 卫星基本信息
type SatelliteInfo struct {
	NoradID     int     `json:"norad_id"`
	Name        string  `json:"name"`
	Category    string  `json:"category,omitempty"`
	Country     string  `json:"country,omitempty"`
	LaunchDate  string  `json:"launch_date,omitempty"`
	Period      float64 `json:"period,omitempty"`      // 周期（分钟）
	Inclination float64 `json:"inclination,omitempty"` // 倾角（度�?
	Apogee      float64 `json:"apogee,omitempty"`      // 远地点（km�?
	Perigee     float64 `json:"perigee,omitempty"`     // 近地点（km�?
	ObjectType  string  `json:"object_type,omitempty"` // PAYLOAD, ROCKET BODY, DEBRIS
}

// SatelliteDetail 卫星详细信息
type SatelliteDetail struct {
	SatelliteInfo
	TLE         *TLEData `json:"tle,omitempty"`
	Description string   `json:"description,omitempty"`
}

// TLEData TLE数据
type TLEData struct {
	NoradID int       `json:"norad_id"`
	Name    string    `json:"name"`
	Line1   string    `json:"line1"`
	Line2   string    `json:"line2"`
	Epoch   time.Time `json:"epoch,omitempty"`
}

// CategoryInfo 卫星分类信息
type CategoryInfo struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	Count       int    `json:"count"`
}

// FavoriteInfo 收藏信息
type FavoriteInfo struct {
	ID            uint      `json:"id"`
	NoradID       int       `json:"norad_id"`
	SatelliteName string    `json:"satellite_name"`
	Notes         string    `json:"notes,omitempty"`
	CreatedAt     time.Time `json:"created_at"`
}

// ProxyKeepTrackResponse KeepTrack代理响应
type ProxyKeepTrackResponse struct {
	Satellites []SatelliteInfo `json:"satellites"`
	Total      int             `json:"total"`
	UpdatedAt  time.Time       `json:"updated_at"`
}

// ProxySpaceTrackLoginResponse SpaceTrack登录响应
type ProxySpaceTrackLoginResponse struct {
	Success   bool   `json:"success"`
	Message   string `json:"message"`
	SessionID string `json:"session_id,omitempty"`
}

// ProxySpaceTrackTLEResponse SpaceTrack TLE响应
type ProxySpaceTrackTLEResponse struct {
	TLEList   []TLEData `json:"tle_list"`
	Total     int       `json:"total"`
	UpdatedAt time.Time `json:"updated_at"`
}
