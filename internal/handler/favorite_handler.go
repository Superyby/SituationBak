package handler

import (
	"strconv"

	"SituationBak/internal/dto/request"
	"SituationBak/internal/middleware"
	"SituationBak/shared/errors"
	"SituationBak/shared/utils"
	"SituationBak/internal/service"
	"github.com/gofiber/fiber/v3"
)

// FavoriteHandler 收藏处理�?
type FavoriteHandler struct {
	satelliteService *service.SatelliteService
}

// NewFavoriteHandler 创建收藏处理器实�?
func NewFavoriteHandler() *FavoriteHandler {
	return &FavoriteHandler{
		satelliteService: service.NewSatelliteService(),
	}
}

// GetFavorites 获取收藏列表
// @Summary 获取收藏列表
// @Tags 收藏
// @Security Bearer
// @Produce json
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(20)
// @Success 200 {object} utils.Response
// @Router /api/v1/favorites [get]
func (h *FavoriteHandler) GetFavorites(c fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	if userID == 0 {
		return utils.FailWithCode(c, errors.CodeUnauthorized)
	}

	page, pageSize := utils.GetPagination(c)

	favorites, total, err := h.satelliteService.GetFavorites(userID, page, pageSize)
	if err != nil {
		return utils.FailWithError(c, err)
	}

	return utils.PagedResponse(c, favorites, page, pageSize, total)
}

// AddFavorite 添加收藏
// @Summary 添加收藏
// @Tags 收藏
// @Security Bearer
// @Accept json
// @Produce json
// @Param body body request.AddFavoriteRequest true "收藏信息"
// @Success 201 {object} utils.Response
// @Router /api/v1/favorites [post]
func (h *FavoriteHandler) AddFavorite(c fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	if userID == 0 {
		return utils.FailWithCode(c, errors.CodeUnauthorized)
	}

	var req request.AddFavoriteRequest
	if err := c.Bind().Body(&req); err != nil {
		return utils.Fail(c, errors.CodeInvalidParams, "请求参数格式错误")
	}

	if req.NoradID <= 0 {
		return utils.Fail(c, errors.CodeInvalidParams, "无效的卫星NORAD ID")
	}

	favorite, err := h.satelliteService.AddFavorite(userID, req.NoradID, req.SatelliteName, req.Notes)
	if err != nil {
		return utils.FailWithError(c, err)
	}

	return utils.Created(c, favorite)
}

// DeleteFavorite 删除收藏
// @Summary 删除收藏
// @Tags 收藏
// @Security Bearer
// @Param id path int true "收藏ID"
// @Success 200 {object} utils.Response
// @Router /api/v1/favorites/{id} [delete]
func (h *FavoriteHandler) DeleteFavorite(c fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	if userID == 0 {
		return utils.FailWithCode(c, errors.CodeUnauthorized)
	}

	idStr := c.Params("id")
	favoriteID, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return utils.Fail(c, errors.CodeInvalidParams, "无效的收藏ID")
	}

	err = h.satelliteService.DeleteFavorite(userID, uint(favoriteID))
	if err != nil {
		return utils.FailWithError(c, err)
	}

	return utils.SuccessWithMessage(c, "删除成功", nil)
}
