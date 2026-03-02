package service

import (
	"time"

	"SituationBak/internal/dto/response"
)

// ProxyService API代理服务（Mock模式）
type ProxyService struct{}

// NewProxyService 创建代理服务实例
func NewProxyService() *ProxyService {
	return &ProxyService{}
}

// GetKeepTrackSatellites 获取KeepTrack卫星数据（Mock）
func (s *ProxyService) GetKeepTrackSatellites() (*response.ProxyKeepTrackResponse, error) {
	// Mock数据
	satellites := []response.SatelliteInfo{
		{NoradID: 25544, Name: "ISS (ZARYA)", Category: "stations", Country: "USA/RUS", ObjectType: "PAYLOAD"},
		{NoradID: 48274, Name: "CSS (TIANHE)", Category: "stations", Country: "CHN", ObjectType: "PAYLOAD"},
		{NoradID: 20580, Name: "HUBBLE SPACE TELESCOPE", Category: "science", Country: "USA", ObjectType: "PAYLOAD"},
		{NoradID: 43013, Name: "STARLINK-1", Category: "communication", Country: "USA", ObjectType: "PAYLOAD"},
		{NoradID: 27424, Name: "ENVISAT", Category: "earth-observation", Country: "ESA", ObjectType: "PAYLOAD"},
	}

	return &response.ProxyKeepTrackResponse{
		Satellites: satellites,
		Total:      len(satellites),
		UpdatedAt:  time.Now().UTC(),
	}, nil
}

// SpaceTrackLogin Space-Track登录（Mock）
func (s *ProxyService) SpaceTrackLogin(username, password string) (*response.ProxySpaceTrackLoginResponse, error) {
	// Mock 登录逻辑 - 简单验证
	if username == "" || password == "" {
		return &response.ProxySpaceTrackLoginResponse{
			Success: false,
			Message: "用户名和密码不能为空",
		}, nil
	}

	// Mock 成功响应
	return &response.ProxySpaceTrackLoginResponse{
		Success:   true,
		Message:   "登录成功（Mock模式）",
		SessionID: "mock-session-" + time.Now().Format("20060102150405"),
	}, nil
}

// GetSpaceTrackTLE 获取Space-Track TLE数据（Mock）
func (s *ProxyService) GetSpaceTrackTLE(noradIDs []int) (*response.ProxySpaceTrackTLEResponse, error) {
	// Mock TLE数据
	tleList := []response.TLEData{
		{
			NoradID: 25544,
			Name:    "ISS (ZARYA)",
			Line1:   "1 25544U 98067A   24061.50000000  .00016717  00000-0  10270-3 0  9993",
			Line2:   "2 25544  51.6400 208.9163 0006703  35.0752 325.0645 15.49560978437867",
			Epoch:   time.Now().UTC(),
		},
		{
			NoradID: 48274,
			Name:    "CSS (TIANHE)",
			Line1:   "1 48274U 21035A   24061.50000000  .00012000  00000-0  90000-4 0  9999",
			Line2:   "2 48274  41.4700 180.0000 0007000  10.0000 350.0000 15.60000000 10000",
			Epoch:   time.Now().UTC(),
		},
		{
			NoradID: 20580,
			Name:    "HUBBLE SPACE TELESCOPE",
			Line1:   "1 20580U 90037B   24061.50000000  .00001000  00000-0  50000-4 0  9999",
			Line2:   "2 20580  28.4700  50.0000 0002500 200.0000 160.0000 15.09000000200000",
			Epoch:   time.Now().UTC(),
		},
	}

	// 如果指定了NORAD ID，则过滤
	if len(noradIDs) > 0 {
		var filtered []response.TLEData
		noradSet := make(map[int]bool)
		for _, id := range noradIDs {
			noradSet[id] = true
		}
		for _, tle := range tleList {
			if noradSet[tle.NoradID] {
				filtered = append(filtered, tle)
			}
		}
		tleList = filtered
	}

	return &response.ProxySpaceTrackTLEResponse{
		TLEList:   tleList,
		Total:     len(tleList),
		UpdatedAt: time.Now().UTC(),
	}, nil
}

// BatchGetTLE 批量获取TLE数据（Mock）
func (s *ProxyService) BatchGetTLE(noradIDs []int) ([]response.TLEData, error) {
	resp, err := s.GetSpaceTrackTLE(noradIDs)
	if err != nil {
		return nil, err
	}
	return resp.TLEList, nil
}
