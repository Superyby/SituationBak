package handler

import (
	"SituationBak/shared/utils"
	"SituationBak/internal/service"
	"github.com/gofiber/fiber/v3"
)

// ProxyHandler API代理处理�?
type ProxyHandler struct {
	proxyService *service.ProxyService
}

// NewProxyHandler 创建代理处理器实�?
func NewProxyHandler() *ProxyHandler {
	return &ProxyHandler{
		proxyService: service.NewProxyService(),
	}
}

// GetKeepTrackSatellites 获取KeepTrack卫星数据
// @Summary 获取KeepTrack卫星数据（Mock模式�?
// @Tags 代理
// @Produce json
// @Success 200 {object} utils.Response
// @Router /api/v1/proxy/keeptrack/sats [get]
func (h *ProxyHandler) GetKeepTrackSatellites(c fiber.Ctx) error {
	result, err := h.proxyService.GetKeepTrackSatellites()
	if err != nil {
		return utils.FailWithError(c, err)
	}
	return utils.Success(c, result)
}

// SpaceTrackLogin Space-Track登录
// @Summary Space-Track登录（Mock模式�?
// @Tags 代理
// @Security Bearer
// @Accept json
// @Produce json
// @Success 200 {object} utils.Response
// @Router /api/v1/proxy/spacetrack/login [post]
func (h *ProxyHandler) SpaceTrackLogin(c fiber.Ctx) error {
	type loginReq struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	var req loginReq
	if err := c.Bind().Body(&req); err != nil {
		req.Username = ""
		req.Password = ""
	}

	result, err := h.proxyService.SpaceTrackLogin(req.Username, req.Password)
	if err != nil {
		return utils.FailWithError(c, err)
	}

	return utils.Success(c, result)
}

// GetSpaceTrackTLE 获取Space-Track TLE数据
// @Summary 获取Space-Track TLE数据（Mock模式�?
// @Tags 代理
// @Security Bearer
// @Produce json
// @Param norad_ids query string false "NORAD ID列表，逗号分隔"
// @Success 200 {object} utils.Response
// @Router /api/v1/proxy/spacetrack/tle [get]
func (h *ProxyHandler) GetSpaceTrackTLE(c fiber.Ctx) error {
	// 解析NORAD ID参数（可选）
	var noradIDs []int
	// 简化处理，返回所有Mock数据

	result, err := h.proxyService.GetSpaceTrackTLE(noradIDs)
	if err != nil {
		return utils.FailWithError(c, err)
	}

	return utils.Success(c, result)
}
