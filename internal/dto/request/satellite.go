package request

// SatelliteListRequest 卫星列表请求
type SatelliteListRequest struct {
	Page      int    `query:"page" validate:"omitempty,min=1"`
	PageSize  int    `query:"page_size" validate:"omitempty,min=1,max=100"`
	Category  string `query:"category" validate:"omitempty"`
	SortBy    string `query:"sort_by" validate:"omitempty,oneof=name norad_id launch_date"`
	SortOrder string `query:"sort_order" validate:"omitempty,oneof=asc desc"`
}

// SatelliteSearchRequest 卫星搜索请求
type SatelliteSearchRequest struct {
	Query    string `query:"q" validate:"required,min=1,max=100"`
	Page     int    `query:"page" validate:"omitempty,min=1"`
	PageSize int    `query:"page_size" validate:"omitempty,min=1,max=100"`
}

// AddFavoriteRequest 添加收藏请求
type AddFavoriteRequest struct {
	NoradID       int    `json:"norad_id" validate:"required"`
	SatelliteName string `json:"satellite_name,omitempty" validate:"omitempty,max=100"`
	Notes         string `json:"notes,omitempty" validate:"omitempty,max=500"`
}

// UpdateFavoriteRequest 更新收藏请求
type UpdateFavoriteRequest struct {
	Notes string `json:"notes,omitempty" validate:"omitempty,max=500"`
}
