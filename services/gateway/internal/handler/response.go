package handler

import (
	"SituationBak/pkg/errors"

	"github.com/gofiber/fiber/v3"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Response 统一响应结构
type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

// Success 成功响应
func Success(c fiber.Ctx, data interface{}) error {
	return c.Status(fiber.StatusOK).JSON(Response{
		Code:    0,
		Message: "success",
		Data:    data,
	})
}

// SuccessWithMessage 成功响应（带消息）
func SuccessWithMessage(c fiber.Ctx, message string, data interface{}) error {
	return c.Status(fiber.StatusOK).JSON(Response{
		Code:    0,
		Message: message,
		Data:    data,
	})
}

// Created 创建成功响应
func Created(c fiber.Ctx, data interface{}) error {
	return c.Status(fiber.StatusCreated).JSON(Response{
		Code:    0,
		Message: "success",
		Data:    data,
	})
}

// Fail 失败响应
func Fail(c fiber.Ctx, code int, message string) error {
	httpStatus := errors.GetHTTPStatus(code)
	return c.Status(httpStatus).JSON(Response{
		Code:    code,
		Message: message,
		Data:    nil,
	})
}

// FailWithCode 失败响应（使用预定义错误码）
func FailWithCode(c fiber.Ctx, code int) error {
	httpStatus := errors.GetHTTPStatus(code)
	return c.Status(httpStatus).JSON(Response{
		Code:    code,
		Message: errors.GetMessage(code),
		Data:    nil,
	})
}

// FailWithError 失败响应（使用错误对象）
func FailWithError(c fiber.Ctx, err error) error {
	if appErr, ok := errors.AsAppError(err); ok {
		httpStatus := errors.GetHTTPStatus(appErr.Code)
		return c.Status(httpStatus).JSON(Response{
			Code:    appErr.Code,
			Message: appErr.Message,
			Data:    nil,
		})
	}
	return c.Status(fiber.StatusInternalServerError).JSON(Response{
		Code:    errors.CodeInternalError,
		Message: "服务器内部错误",
		Data:    nil,
	})
}

// FailWithGRPCError 处理gRPC错误响应
func FailWithGRPCError(c fiber.Ctx, err error) error {
	st, ok := status.FromError(err)
	if !ok {
		return c.Status(fiber.StatusInternalServerError).JSON(Response{
			Code:    errors.CodeInternalError,
			Message: "服务调用失败",
			Data:    nil,
		})
	}

	httpStatus := grpcCodeToHTTPStatus(st.Code())
	code := grpcCodeToAppCode(st.Code())

	return c.Status(httpStatus).JSON(Response{
		Code:    code,
		Message: st.Message(),
		Data:    nil,
	})
}

// grpcCodeToHTTPStatus gRPC状态码转HTTP状态码
func grpcCodeToHTTPStatus(code codes.Code) int {
	switch code {
	case codes.OK:
		return fiber.StatusOK
	case codes.InvalidArgument:
		return fiber.StatusBadRequest
	case codes.NotFound:
		return fiber.StatusNotFound
	case codes.AlreadyExists:
		return fiber.StatusConflict
	case codes.Unauthenticated:
		return fiber.StatusUnauthorized
	case codes.PermissionDenied:
		return fiber.StatusForbidden
	case codes.Unavailable:
		return fiber.StatusServiceUnavailable
	default:
		return fiber.StatusInternalServerError
	}
}

// grpcCodeToAppCode gRPC状态码转应用错误码
func grpcCodeToAppCode(code codes.Code) int {
	switch code {
	case codes.OK:
		return errors.CodeSuccess
	case codes.InvalidArgument:
		return errors.CodeInvalidParams
	case codes.NotFound:
		return errors.CodeNotFound
	case codes.AlreadyExists:
		return errors.CodeAlreadyExists
	case codes.Unauthenticated:
		return errors.CodeUnauthorized
	case codes.PermissionDenied:
		return errors.CodeForbidden
	case codes.Unavailable:
		return errors.CodeServiceUnavailable
	default:
		return errors.CodeInternalError
	}
}
