package constants

// ==================== 上下文键 ====================

const (
	// CtxKeyTraceID 链路追踪ID
	CtxKeyTraceID = "traceID"
	// CtxKeyUserID 用户ID
	CtxKeyUserID = "userID"
	// CtxKeyUsername 用户名
	CtxKeyUsername = "username"
	// CtxKeyRole 用户角色
	CtxKeyRole = "role"
	// CtxKeyRequestID 请求ID (别名)
	CtxKeyRequestID = "requestID"
)

// ==================== HTTP 头 ====================

const (
	// HeaderTraceID 链路追踪ID请求头
	HeaderTraceID = "X-Trace-ID"
	// HeaderRequestID 请求ID请求头
	HeaderRequestID = "X-Request-ID"
	// HeaderAuthorization 认证请求头
	HeaderAuthorization = "Authorization"
	// HeaderContentType 内容类型请求头
	HeaderContentType = "Content-Type"
	// HeaderAccept 接受类型请求头
	HeaderAccept = "Accept"
)

// ==================== 用户角色 ====================

const (
	// RoleUser 普通用户
	RoleUser = "user"
	// RoleAdmin 管理员
	RoleAdmin = "admin"
)

// ==================== 环境类型 ====================

const (
	// EnvDevelopment 开发环境
	EnvDevelopment = "development"
	// EnvProduction 生产环境
	EnvProduction = "production"
	// EnvTesting 测试环境
	EnvTesting = "testing"
)

// ==================== 时间格式 ====================

const (
	// TimeFormatDefault 默认时间格式
	TimeFormatDefault = "2006-01-02 15:04:05"
	// TimeFormatDate 日期格式
	TimeFormatDate = "2006-01-02"
	// TimeFormatTime 时间格式
	TimeFormatTime = "15:04:05"
	// TimeFormatISO8601 ISO8601格式
	TimeFormatISO8601 = "2006-01-02T15:04:05Z07:00"
)

// ==================== 分页默认值 ====================

const (
	// DefaultPage 默认页码
	DefaultPage = 1
	// DefaultPageSize 默认每页数量
	DefaultPageSize = 20
	// MaxPageSize 最大每页数量
	MaxPageSize = 100
)

// ==================== Token 类型 ====================

const (
	// TokenTypeBearer Bearer Token类型
	TokenTypeBearer = "Bearer"
)

// ==================== 应用常量 ====================

const (
	// AppName 应用名称
	AppName = "SituationBak"
	// AppVersion 应用版本
	AppVersion = "1.0.0"
)

// ==================== 缓存键前缀 ====================

const (
	// CacheKeyPrefixUser 用户缓存前缀
	CacheKeyPrefixUser = "user:"
	// CacheKeyPrefixToken Token缓存前缀
	CacheKeyPrefixToken = "token:"
	// CacheKeyPrefixTLE TLE缓存前缀
	CacheKeyPrefixTLE = "tle:"
	// CacheKeyPrefixSatellite 卫星数据缓存前缀
	CacheKeyPrefixSatellite = "satellite:"
)

// ==================== 缓存过期时间(秒) ====================

const (
	// CacheTTLShort 短期缓存 5分钟
	CacheTTLShort = 300
	// CacheTTLMedium 中期缓存 1小时
	CacheTTLMedium = 3600
	// CacheTTLLong 长期缓存 24小时
	CacheTTLLong = 86400
)
