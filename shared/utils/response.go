package utils

import (
	"strconv"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"

	"SituationBak/shared/errors"
)

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

// PagedData 分页数据
type PagedData struct {
	Items      interface{} `json:"items"`
	Pagination Pagination  `json:"pagination"`
}

// Success 成功响应
func Success(c fiber.Ctx, data interface{}) error {
	return c.Status(fiber.StatusOK).JSON(Response{
		Code:    0,
		Message: "success",
		Data:    data,
		Success: true,
	})
}

// SuccessWithMessage 带消息的成功响应
func SuccessWithMessage(c fiber.Ctx, message string, data interface{}) error {
	return c.Status(fiber.StatusOK).JSON(Response{
		Code:    0,
		Message: message,
		Data:    data,
		Success: true,
	})
}

// Created 创建成功响应
func Created(c fiber.Ctx, data interface{}) error {
	return c.Status(fiber.StatusCreated).JSON(Response{
		Code:    0,
		Message: "created",
		Data:    data,
		Success: true,
	})
}

// Fail 失败响应
func Fail(c fiber.Ctx, code int, message string) error {
	httpStatus := errors.GetHTTPStatus(code)
	return c.Status(httpStatus).JSON(Response{
		Code:    code,
		Message: message,
		Data:    nil,
		Success: false,
	})
}

// FailWithCode 根据错误码返回失败响应
func FailWithCode(c fiber.Ctx, code int) error {
	httpStatus := errors.GetHTTPStatus(code)
	return c.Status(httpStatus).JSON(Response{
		Code:    code,
		Message: errors.GetMessage(code),
		Data:    nil,
		Success: false,
	})
}

// FailWithError 根据错误返回失败响应
func FailWithError(c fiber.Ctx, err error) error {
	if appErr, ok := errors.AsAppError(err); ok {
		httpStatus := errors.GetHTTPStatus(appErr.Code)
		return c.Status(httpStatus).JSON(Response{
			Code:    appErr.Code,
			Message: appErr.Message,
			Data:    nil,
			Success: false,
		})
	}
	return c.Status(fiber.StatusInternalServerError).JSON(Response{
		Code:    errors.CodeInternalError,
		Message: "服务器内部错误",
		Data:    nil,
		Success: false,
	})
}

// PagedResponse 分页响应
func PagedResponse(c fiber.Ctx, items interface{}, page, pageSize int, total int64) error {
	totalPages := int(total) / pageSize
	if int(total)%pageSize > 0 {
		totalPages++
	}

	return c.Status(fiber.StatusOK).JSON(Response{
		Code:    0,
		Message: "success",
		Data: PagedData{
			Items: items,
			Pagination: Pagination{
				Page:       page,
				PageSize:   pageSize,
				Total:      total,
				TotalPages: totalPages,
			},
		},
		Success: true,
	})
}

// GetPagination 从请求中获取分页参数
func GetPagination(c fiber.Ctx) (int, int) {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	pageSize, _ := strconv.Atoi(c.Query("page_size", "20"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 20
	}
	if pageSize > 100 {
		pageSize = 100
	}

	return page, pageSize
}

// GenerateRequestID 生成请求ID
func GenerateRequestID() string {
	return uuid.New().String()
}
