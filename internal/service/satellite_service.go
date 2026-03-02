package service

import (
	"strings"
	"time"

	"SituationBak/internal/dto/response"
	"SituationBak/internal/model"
	"SituationBak/internal/pkg/errors"
	"SituationBak/internal/repository"
)

// SatelliteService 卫星服务
type SatelliteService struct {
	favoriteRepo *repository.FavoriteRepository
}

// NewSatelliteService 创建卫星服务实例
func NewSatelliteService() *SatelliteService {
	return &SatelliteService{
		favoriteRepo: repository.NewFavoriteRepository(),
	}
}

// Mock数据 - 卫星列表
var mockSatellites = []response.SatelliteInfo{
	{NoradID: 25544, Name: "ISS (ZARYA)", Category: "stations", Country: "USA/RUS", LaunchDate: "1998-11-20", Period: 92.9, Inclination: 51.6, Apogee: 422, Perigee: 418, ObjectType: "PAYLOAD"},
	{NoradID: 48274, Name: "CSS (TIANHE)", Category: "stations", Country: "CHN", LaunchDate: "2021-04-29", Period: 92.2, Inclination: 41.5, Apogee: 390, Perigee: 385, ObjectType: "PAYLOAD"},
	{NoradID: 20580, Name: "HUBBLE SPACE TELESCOPE", Category: "science", Country: "USA", LaunchDate: "1990-04-24", Period: 95.4, Inclination: 28.5, Apogee: 540, Perigee: 535, ObjectType: "PAYLOAD"},
	{NoradID: 43013, Name: "STARLINK-1", Category: "communication", Country: "USA", LaunchDate: "2018-02-22", Period: 95.0, Inclination: 53.0, Apogee: 550, Perigee: 540, ObjectType: "PAYLOAD"},
	{NoradID: 27424, Name: "ENVISAT", Category: "earth-observation", Country: "ESA", LaunchDate: "2002-03-01", Period: 100.6, Inclination: 98.5, Apogee: 790, Perigee: 785, ObjectType: "PAYLOAD"},
	{NoradID: 28654, Name: "NOAA 18", Category: "weather", Country: "USA", LaunchDate: "2005-05-20", Period: 102.1, Inclination: 99.0, Apogee: 870, Perigee: 850, ObjectType: "PAYLOAD"},
	{NoradID: 37849, Name: "TIANGONG 1", Category: "stations", Country: "CHN", LaunchDate: "2011-09-29", Period: 88.6, Inclination: 42.8, Apogee: 355, Perigee: 330, ObjectType: "PAYLOAD"},
	{NoradID: 33591, Name: "GOES 14", Category: "weather", Country: "USA", LaunchDate: "2009-06-27", Period: 1436.1, Inclination: 0.1, Apogee: 35800, Perigee: 35780, ObjectType: "PAYLOAD"},
	{NoradID: 39084, Name: "LANDSAT 8", Category: "earth-observation", Country: "USA", LaunchDate: "2013-02-11", Period: 98.9, Inclination: 98.2, Apogee: 710, Perigee: 705, ObjectType: "PAYLOAD"},
	{NoradID: 41866, Name: "JASON 3", Category: "science", Country: "USA", LaunchDate: "2016-01-17", Period: 112.4, Inclination: 66.0, Apogee: 1340, Perigee: 1330, ObjectType: "PAYLOAD"},
}

// Mock数据 - TLE
var mockTLEs = map[int]*response.TLEData{
	25544: {NoradID: 25544, Name: "ISS (ZARYA)", Line1: "1 25544U 98067A   24061.50000000  .00016717  00000-0  10270-3 0  9993", Line2: "2 25544  51.6400 208.9163 0006703  35.0752 325.0645 15.49560978437867"},
	48274: {NoradID: 48274, Name: "CSS (TIANHE)", Line1: "1 48274U 21035A   24061.50000000  .00012000  00000-0  90000-4 0  9999", Line2: "2 48274  41.4700 180.0000 0007000  10.0000 350.0000 15.60000000 10000"},
	20580: {NoradID: 20580, Name: "HUBBLE SPACE TELESCOPE", Line1: "1 20580U 90037B   24061.50000000  .00001000  00000-0  50000-4 0  9999", Line2: "2 20580  28.4700  50.0000 0002500 200.0000 160.0000 15.09000000200000"},
}

// Mock数据 - 分类
var mockCategories = []response.CategoryInfo{
	{ID: "stations", Name: "空间站", Description: "包括国际空间站、中国空间站等", Count: 3},
	{ID: "communication", Name: "通信卫星", Description: "用于通信服务的卫星", Count: 5000},
	{ID: "weather", Name: "气象卫星", Description: "用于天气预报的卫星", Count: 50},
	{ID: "earth-observation", Name: "地球观测", Description: "用于观测地球的卫星", Count: 200},
	{ID: "science", Name: "科学卫星", Description: "用于科学研究的卫星", Count: 100},
	{ID: "navigation", Name: "导航卫星", Description: "GPS、北斗等导航卫星", Count: 120},
	{ID: "military", Name: "军事卫星", Description: "军事用途卫星", Count: 500},
	{ID: "debris", Name: "太空碎片", Description: "太空碎片和废弃卫星", Count: 30000},
}

// GetSatellites 获取卫星列表
func (s *SatelliteService) GetSatellites(page, pageSize int, category string) ([]response.SatelliteInfo, int64, error) {
	// 过滤
	var filtered []response.SatelliteInfo
	for _, sat := range mockSatellites {
		if category == "" || sat.Category == category {
			filtered = append(filtered, sat)
		}
	}

	total := int64(len(filtered))

	// 分页
	start := (page - 1) * pageSize
	if start >= len(filtered) {
		return []response.SatelliteInfo{}, total, nil
	}
	end := start + pageSize
	if end > len(filtered) {
		end = len(filtered)
	}

	return filtered[start:end], total, nil
}

// GetSatelliteByID 根据NORAD ID获取卫星详情
func (s *SatelliteService) GetSatelliteByID(noradID int) (*response.SatelliteDetail, error) {
	for _, sat := range mockSatellites {
		if sat.NoradID == noradID {
			detail := &response.SatelliteDetail{
				SatelliteInfo: sat,
				Description:   "这是 " + sat.Name + " 卫星的详细描述信息。",
			}
			if tle, ok := mockTLEs[noradID]; ok {
				detail.TLE = tle
			}
			return detail, nil
		}
	}
	return nil, errors.WithCode(errors.CodeSatelliteNotFound)
}

// GetSatelliteTLE 获取卫星TLE数据
func (s *SatelliteService) GetSatelliteTLE(noradID int) (*response.TLEData, error) {
	if tle, ok := mockTLEs[noradID]; ok {
		tle.Epoch = time.Now().UTC()
		return tle, nil
	}
	return nil, errors.WithCode(errors.CodeSatelliteNotFound)
}

// SearchSatellites 搜索卫星
func (s *SatelliteService) SearchSatellites(query string, page, pageSize int) ([]response.SatelliteInfo, int64, error) {
	query = strings.ToLower(query)
	var results []response.SatelliteInfo

	for _, sat := range mockSatellites {
		if strings.Contains(strings.ToLower(sat.Name), query) {
			results = append(results, sat)
		}
	}

	total := int64(len(results))

	// 分页
	start := (page - 1) * pageSize
	if start >= len(results) {
		return []response.SatelliteInfo{}, total, nil
	}
	end := start + pageSize
	if end > len(results) {
		end = len(results)
	}

	return results[start:end], total, nil
}

// GetCategories 获取卫星分类列表
func (s *SatelliteService) GetCategories() []response.CategoryInfo {
	return mockCategories
}

// GetFavorites 获取收藏列表
func (s *SatelliteService) GetFavorites(userID uint, page, pageSize int) ([]response.FavoriteInfo, int64, error) {
	favorites, total, err := s.favoriteRepo.FindByUserID(userID, page, pageSize)
	if err != nil {
		return nil, 0, errors.ErrInternal(err)
	}

	var result []response.FavoriteInfo
	for _, f := range favorites {
		result = append(result, response.FavoriteInfo{
			ID:            f.ID,
			NoradID:       f.NoradID,
			SatelliteName: f.SatelliteName,
			Notes:         f.Notes,
			CreatedAt:     f.CreatedAt,
		})
	}

	return result, total, nil
}

// AddFavorite 添加收藏
func (s *SatelliteService) AddFavorite(userID uint, noradID int, name, notes string) (*response.FavoriteInfo, error) {
	// 检查是否已收藏
	exists, err := s.favoriteRepo.Exists(userID, noradID)
	if err != nil {
		return nil, errors.ErrInternal(err)
	}
	if exists {
		return nil, errors.WithCode(errors.CodeFavoriteExists)
	}

	favorite := &model.Favorite{
		UserID:        userID,
		NoradID:       noradID,
		SatelliteName: name,
		Notes:         notes,
	}

	if err := s.favoriteRepo.Create(favorite); err != nil {
		return nil, errors.ErrInternal(err)
	}

	return &response.FavoriteInfo{
		ID:            favorite.ID,
		NoradID:       favorite.NoradID,
		SatelliteName: favorite.SatelliteName,
		Notes:         favorite.Notes,
		CreatedAt:     favorite.CreatedAt,
	}, nil
}

// DeleteFavorite 删除收藏
func (s *SatelliteService) DeleteFavorite(userID uint, favoriteID uint) error {
	// 检查收藏是否存在且属于该用户
	favorite, err := s.favoriteRepo.FindByID(favoriteID)
	if err != nil {
		return errors.ErrInternal(err)
	}
	if favorite == nil {
		return errors.WithCode(errors.CodeFavoriteNotFound)
	}
	if favorite.UserID != userID {
		return errors.ErrForbidden()
	}

	if err := s.favoriteRepo.Delete(favoriteID); err != nil {
		return errors.ErrInternal(err)
	}

	return nil
}
