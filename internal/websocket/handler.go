package websocket

import (
	"SituationBak/shared/logger"
	"SituationBak/internal/service"
	"github.com/fasthttp/websocket"
	"github.com/gofiber/fiber/v3"
	"github.com/valyala/fasthttp"
)

var upgrader = websocket.FastHTTPUpgrader{
	CheckOrigin: func(ctx *fasthttp.RequestCtx) bool {
		return true // 允许所有来源
	},
}

// HandleWebSocket WebSocket连接处理
func HandleWebSocket(c fiber.Ctx) error {
	// 获取token参数（可选认证）
	token := c.Query("token")
	var userID uint = 0

	if token != "" {
		// 验证Token
		authService := service.NewAuthService()
		claims, err := authService.ValidateToken(token)
		if err == nil {
			userID = claims.UserID
		}
	}

	// 获取底层 fasthttp context
	ctx := c.Context()

	err := upgrader.Upgrade(ctx, func(conn *websocket.Conn) {
		hub := GetHub()
		client := NewClient(hub, conn, userID)

		// 注册客户端
		hub.register <- client

		logger.Info("WebSocket connection established",
			logger.Uint("user_id", userID),
			logger.String("remote_addr", conn.RemoteAddr().String()),
		)

		// 启动读写协程
		go client.WritePump()
		client.ReadPump()
	})

	if err != nil {
		logger.Error("WebSocket upgrade failed", logger.Err(err))
		return fiber.ErrUpgradeRequired
	}

	return nil
}

// BroadcastSatelliteUpdate 广播卫星位置更新（供外部调用）
func BroadcastSatelliteUpdate(satellites []SatellitePosition) {
	hub := GetHub()
	msg, _ := NewMessage(MessageTypeSatelliteUpdate, &SatelliteUpdatePayload{
		Satellites: satellites,
	})
	hub.Broadcast(msg)
}

// BroadcastNotification 广播系统通知
func BroadcastNotification(title, message, level string) {
	hub := GetHub()
	msg, _ := NewMessage(MessageTypeNotification, &NotificationPayload{
		Title:   title,
		Message: message,
		Level:   level,
	})
	hub.Broadcast(msg)
}

// GetConnectionCount 获取当前WebSocket连接数
func GetConnectionCount() int {
	return GetHub().GetClientCount()
}
