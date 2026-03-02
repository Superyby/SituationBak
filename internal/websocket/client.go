package websocket

import (
	"encoding/json"
	"time"

	"SituationBak/shared/logger"
	"github.com/fasthttp/websocket"
)

const (
	// 写超�?
	writeWait = 10 * time.Second

	// Pong超时
	pongWait = 60 * time.Second

	// Ping间隔
	pingPeriod = (pongWait * 9) / 10

	// 消息大小限制
	maxMessageSize = 512 * 1024 // 512KB
)

// Client WebSocket客户�?
type Client struct {
	hub           *Hub
	conn          *websocket.Conn
	send          chan *Message
	userID        uint
	subscriptions map[int]bool // 订阅的卫�?noradID
}

// NewClient 创建新客户端
func NewClient(hub *Hub, conn *websocket.Conn, userID uint) *Client {
	return &Client{
		hub:           hub,
		conn:          conn,
		send:          make(chan *Message, 256),
		userID:        userID,
		subscriptions: make(map[int]bool),
	}
}

// ReadPump 读取消息
func (c *Client) ReadPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()

	c.conn.SetReadLimit(maxMessageSize)
	_ = c.conn.SetReadDeadline(time.Now().Add(pongWait))

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				logger.Error("WebSocket read error", logger.Err(err))
			}
			break
		}

		// 重置读取超时
		_ = c.conn.SetReadDeadline(time.Now().Add(pongWait))

		// 处理消息
		c.handleMessage(message)
	}
}

// WritePump 写入消息
func (c *Client) WritePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			_ = c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// Hub关闭了通道
				_ = c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			data, err := message.ToJSON()
			if err != nil {
				logger.Error("Failed to marshal message", logger.Err(err))
				continue
			}

			if err := c.conn.WriteMessage(websocket.TextMessage, data); err != nil {
				logger.Error("WebSocket write error", logger.Err(err))
				return
			}

		case <-ticker.C:
			_ = c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			// 发送ping消息
			pingMsg, _ := NewMessage(MessageTypePing, nil)
			data, _ := pingMsg.ToJSON()
			if err := c.conn.WriteMessage(websocket.TextMessage, data); err != nil {
				return
			}
		}
	}
}

// handleMessage 处理接收到的消息
func (c *Client) handleMessage(data []byte) {
	var msg Message
	if err := json.Unmarshal(data, &msg); err != nil {
		logger.Error("Failed to unmarshal message", logger.Err(err))
		c.sendError(1001, "消息格式错误")
		return
	}

	switch msg.Type {
	case MessageTypePing:
		c.handlePing()
	case MessageTypeSubscribe:
		c.handleSubscribe(&msg)
	case MessageTypeUnsubscribe:
		c.handleUnsubscribe(&msg)
	default:
		c.sendError(1002, "未知消息类型")
	}
}

// handlePing 处理Ping消息
func (c *Client) handlePing() {
	pongMsg, _ := NewMessage(MessageTypePong, nil)
	c.send <- pongMsg
}

// handleSubscribe 处理订阅消息
func (c *Client) handleSubscribe(msg *Message) {
	var payload SubscribePayload
	if err := msg.ParsePayload(&payload); err != nil {
		c.sendError(1001, "订阅参数格式错误")
		return
	}

	if len(payload.NoradIDs) == 0 {
		c.sendError(1001, "订阅列表不能为空")
		return
	}

	c.hub.Subscribe(c, payload.NoradIDs)

	// 发送确认消�?
	confirmMsg, _ := NewMessage(MessageTypeNotification, &NotificationPayload{
		Title:   "订阅成功",
		Message: "已成功订阅卫星数�?,
		Level:   "info",
	})
	c.send <- confirmMsg
}

// handleUnsubscribe 处理取消订阅消息
func (c *Client) handleUnsubscribe(msg *Message) {
	var payload UnsubscribePayload
	if err := msg.ParsePayload(&payload); err != nil {
		c.sendError(1001, "取消订阅参数格式错误")
		return
	}

	c.hub.Unsubscribe(c, payload.NoradIDs)

	// 发送确认消�?
	confirmMsg, _ := NewMessage(MessageTypeNotification, &NotificationPayload{
		Title:   "取消订阅成功",
		Message: "已成功取消订�?,
		Level:   "info",
	})
	c.send <- confirmMsg
}

// sendError 发送错误消�?
func (c *Client) sendError(code int, message string) {
	errMsg, _ := NewMessage(MessageTypeError, &ErrorPayload{
		Code:    code,
		Message: message,
	})
	c.send <- errMsg
}

// SendSatelliteUpdate 发送卫星更�?
func (c *Client) SendSatelliteUpdate(satellites []SatellitePosition) {
	msg, _ := NewMessage(MessageTypeSatelliteUpdate, &SatelliteUpdatePayload{
		Satellites: satellites,
	})
	c.send <- msg
}
