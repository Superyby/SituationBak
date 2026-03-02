package handler

import (
	"strconv"

	"SituationBak/shared/errors"
	"SituationBak/shared/utils"
	"SituationBak/internal/service"
	"github.com/gofiber/fiber/v3"
)

// SatelliteHandler 卫星处理器
type SatelliteHandler struct {
	satelliteService *service.SatelliteService
}

// NewSatelliteHandler 创建卫星处理器实例
func NewSatelliteHandler() *SatelliteHandler {
	return &SatelliteHandler{
		satelliteService: service.NewSatelliteService(),
	}
}

// GetSatellites 获取卫星列表
// @Summary 获取卫星列表
// @Tags 卫星
// @Produce json
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(20)
// @Param category query string false "分类"
// @Success 200 {object} utils.Response
// @Router /api/v1/satellites [get]
func (h *SatelliteHandler) GetSatellites(c fiber.Ctx) error {
	page, pageSize := utils.GetPagination(c)
	category := c.Query("category")

	satellites, total, err := h.satelliteService.GetSatellites(page, pageSize, category)
	if err != nil {
		return utils.FailWithError(c, err)
	}

	return utils.PagedResponse(c, satellites, page, pageSize, total)
}

// GetSatelliteByID 获取卫星详情
// @Summary 获取卫星详情
// @Tags 卫星
// @Produce json
// @Param id path int true "NORAD ID"
// @Success 200 {object} utils.Response
// @Router /api/v1/satellites/{id} [get]
func (h *SatelliteHandler) GetSatelliteByID(c fiber.Ctx) error {
	idStr := c.Params("id")
	noradID, err := strconv.Atoi(idStr)
	if err != nil {
		return utils.Fail(c, errors.CodeInvalidParams, "无效的卫星ID")
	}

	satellite, err := h.satelliteService.GetSatelliteByID(noradID)
	if err != nil {
		return utils.FailWithError(c, err)
	}

	return utils.Success(c, satellite)
}

// GetSatelliteTLE 获取卫星TLE数据
// @Summary 获取卫星TLE数据
// @Tags 卫星
// @Produce json
// @Param id path int true "NORAD ID"
// @Success 200 {object} utils.Response
// @Router /api/v1/satellites/{id}/tle [get]
func (h *SatelliteHandler) GetSatelliteTLE(c fiber.Ctx) error {
	idStr := c.Params("id")
	noradID, err := strconv.Atoi(idStr)
	if err != nil {
		return utils.Fail(c, errors.CodeInvalidParams, "无效的卫星ID")
	}

	tle, err := h.satelliteService.GetSatelliteTLE(noradID)
	if err != nil {
		return utils.FailWithError(c, err)
	}

	return utils.Success(c, tle)
}

// SearchSatellites 搜索卫星
// @Summary 搜索卫星
// @Tags 卫星
// @Produce json
// @Param q query string true "搜索关键词"
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(20)
// @Success 200 {object} utils.Response
// @Router /api/v1/satellites/search [get]
func (h *SatelliteHandler) SearchSatellites(c fiber.Ctx) error {
	query := c.Query("q")
	if query == "" {
		return utils.Fail(c, errors.CodeInvalidParams, "搜索关键词不能为空")
	}

	page, pageSize := utils.GetPagination(c)

	satellites, total, err := h.satelliteService.SearchSatellites(query, page, pageSize)
	if err != nil {
		return utils.FailWithError(c, err)
	}

	return utils.PagedResponse(c, satellites, page, pageSize, total)
}

// GetCategories 获取卫星分类列表
// @Summary 获取卫星分类列表
// @Tags 卫星
// @Produce json
// @Success 200 {object} utils.Response
// @Router /api/v1/satellites/categories [get]
func (h *SatelliteHandler) GetCategories(c fiber.Ctx) error {
	categories := h.satelliteService.GetCategories()
	return utils.Success(c, categories)
}
