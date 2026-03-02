package router

import (
	"SituationBak/internal/handler"
	"SituationBak/internal/middleware"
	wsHandler "SituationBak/internal/websocket"
	"github.com/gofiber/fiber/v3"
)

// SetupRoutes 配置路由
func SetupRoutes(app *fiber.App) {
	// 全局中间件
	app.Use(middleware.RecoveryMiddleware())
	app.Use(middleware.LoggerMiddleware())
	app.Use(middleware.CORSMiddleware())
	app.Use(middleware.RateLimitMiddleware())

	// Swagger 文档路由
	setupSwagger(app)

	// 健康检查
	app.Get("/health", func(c fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "ok",
			"service": "SituationBak",
			"docs":    "/swagger/index.html",
		})
	})

	// API v1
	api := app.Group("/api/v1")

	// 初始化处理器
	authHandler := handler.NewAuthHandler()
	userHandler := handler.NewUserHandler()
	satelliteHandler := handler.NewSatelliteHandler()
	favoriteHandler := handler.NewFavoriteHandler()
	proxyHandler := handler.NewProxyHandler()

	// 认证路由（无需认证）
	auth := api.Group("/auth")
	auth.Post("/register", authHandler.Register)
	auth.Post("/login", authHandler.Login)
	auth.Post("/refresh", authHandler.RefreshToken)

	// 认证路由（需要认证）
	authProtected := api.Group("/auth", middleware.AuthMiddleware())
	authProtected.Post("/logout", authHandler.Logout)
	authProtected.Get("/me", authHandler.GetMe)

	// 用户路由（需要认证）
	user := api.Group("/user", middleware.AuthMiddleware())
	user.Get("/profile", userHandler.GetProfile)
	user.Put("/profile", userHandler.UpdateProfile)
	user.Put("/password", userHandler.ChangePassword)
	user.Get("/settings", userHandler.GetSettings)
	user.Put("/settings", userHandler.UpdateSettings)

	// 卫星路由（无需认证）
	satellites := api.Group("/satellites")
	satellites.Get("/", satelliteHandler.GetSatellites)
	satellites.Get("/search", satelliteHandler.SearchSatellites)
	satellites.Get("/categories", satelliteHandler.GetCategories)
	satellites.Get("/:id", satelliteHandler.GetSatelliteByID)
	satellites.Get("/:id/tle", satelliteHandler.GetSatelliteTLE)

	// 收藏路由（需要认证）
	favorites := api.Group("/favorites", middleware.AuthMiddleware())
	favorites.Get("/", favoriteHandler.GetFavorites)
	favorites.Post("/", favoriteHandler.AddFavorite)
	favorites.Delete("/:id", favoriteHandler.DeleteFavorite)

	// 代理路由
	proxy := api.Group("/proxy")
	proxy.Get("/keeptrack/sats", proxyHandler.GetKeepTrackSatellites) // 无需认证

	// Space-Track代理（需要认证）
	proxyProtected := api.Group("/proxy", middleware.AuthMiddleware())
	proxyProtected.Post("/spacetrack/login", proxyHandler.SpaceTrackLogin)
	proxyProtected.Get("/spacetrack/tle", proxyHandler.GetSpaceTrackTLE)

	// WebSocket路由
	app.Get("/ws", wsHandler.HandleWebSocket)
}
