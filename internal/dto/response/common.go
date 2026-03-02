package response

// Response 统一响应结构
type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Success bool        `json:"success"`
}

// Pagination 分页信息
type Pagination struct {
	Page       int   `json:"page"`
	PageSize   int   `json:"page_size"`
	Total      int64 `json:"total"`
	TotalPages int   `json:"total_pages"`
}

// PagedResponse 分页响应数据
type PagedResponse struct {
	Items      interface{} `json:"items"`
	Pagination Pagination  `json:"pagination"`
}

// NewResponse 创建响应
func NewResponse(code int, message string, data interface{}, success bool) *Response {
	return &Response{
		Code:    code,
		Message: message,
		Data:    data,
		Success: success,
	}
}

// Success 成功响应
func Success(data interface{}) *Response {
	return &Response{
		Code:    0,
		Message: "success",
		Data:    data,
		Success: true,
	}
}

// SuccessWithMessage 带消息的成功响应
func SuccessWithMessage(message string, data interface{}) *Response {
	return &Response{
		Code:    0,
		Message: message,
		Data:    data,
		Success: true,
	}
}

// Fail 失败响应
func Fail(code int, message string) *Response {
	return &Response{
		Code:    code,
		Message: message,
		Data:    nil,
		Success: false,
	}
}

// NewPagedResponse 创建分页响应
func NewPagedResponse(items interface{}, page, pageSize int, total int64) *PagedResponse {
	totalPages := int(total) / pageSize
	if int(total)%pageSize > 0 {
		totalPages++
	}

	return &PagedResponse{
		Items: items,
		Pagination: Pagination{
			Page:       page,
			PageSize:   pageSize,
			Total:      total,
			TotalPages: totalPages,
		},
	}
}
