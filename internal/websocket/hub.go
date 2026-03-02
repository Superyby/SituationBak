package websocket

import (
	"sync"

	"SituationBak/shared/logger"
)

// Hub WebSocket连接管理中心
type Hub struct {
	// 注册的客户端
	clients map[*Client]bool

	// 客户端订阅的卫星
	subscriptions map[int]map[*Client]bool // noradID -> clients

	// 广播通道
	broadcast chan *Message

	// 客户端注册通道
	register chan *Client

	// 客户端注销通道
	unregister chan *Client

	// 互斥�?
	mu sync.RWMutex
}

// 全局Hub实例
var hub *Hub
var once sync.Once

// GetHub 获取Hub单例
func GetHub() *Hub {
	once.Do(func() {
		hub = &Hub{
			clients:       make(map[*Client]bool),
			subscriptions: make(map[int]map[*Client]bool),
			broadcast:     make(chan *Message, 256),
			register:      make(chan *Client),
			unregister:    make(chan *Client),
		}
		go hub.run()
	})
	return hub
}

// run 运行Hub
func (h *Hub) run() {
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			h.clients[client] = true
			h.mu.Unlock()
			logger.Info("WebSocket client connected",
				logger.Uint("user_id", client.userID),
				logger.Int("total_clients", len(h.clients)),
			)

		case client := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)

				// 清理订阅
				for noradID, clients := range h.subscriptions {
					delete(clients, client)
					if len(clients) == 0 {
						delete(h.subscriptions, noradID)
					}
				}
			}
			h.mu.Unlock()
			logger.Info("WebSocket client disconnected",
				logger.Uint("user_id", client.userID),
				logger.Int("total_clients", len(h.clients)),
			)

		case message := <-h.broadcast:
			h.mu.RLock()
			for client := range h.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}
			h.mu.RUnlock()
		}
	}
}

// Subscribe 订阅卫星
func (h *Hub) Subscribe(client *Client, noradIDs []int) {
	h.mu.Lock()
	defer h.mu.Unlock()

	for _, noradID := range noradIDs {
		if h.subscriptions[noradID] == nil {
			h.subscriptions[noradID] = make(map[*Client]bool)
		}
		h.subscriptions[noradID][client] = true
		client.subscriptions[noradID] = true
	}

	logger.Debug("Client subscribed to satellites",
		logger.Uint("user_id", client.userID),
		logger.Any("norad_ids", noradIDs),
	)
}

// Unsubscribe 取消订阅卫星
func (h *Hub) Unsubscribe(client *Client, noradIDs []int) {
	h.mu.Lock()
	defer h.mu.Unlock()

	for _, noradID := range noradIDs {
		if clients, ok := h.subscriptions[noradID]; ok {
			delete(clients, client)
			if len(clients) == 0 {
				delete(h.subscriptions, noradID)
			}
		}
		delete(client.subscriptions, noradID)
	}

	logger.Debug("Client unsubscribed from satellites",
		logger.Uint("user_id", client.userID),
		logger.Any("norad_ids", noradIDs),
	)
}

// BroadcastToSubscribers 向订阅了指定卫星的客户端广播消息
func (h *Hub) BroadcastToSubscribers(noradID int, message *Message) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	if clients, ok := h.subscriptions[noradID]; ok {
		for client := range clients {
			select {
			case client.send <- message:
			default:
				// 发送失败，跳过
			}
		}
	}
}

// Broadcast 向所有客户端广播消息
func (h *Hub) Broadcast(message *Message) {
	h.broadcast <- message
}

// GetClientCount 获取当前连接�?
func (h *Hub) GetClientCount() int {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return len(h.clients)
}

// GetSubscriberCount 获取指定卫星的订阅者数�?
func (h *Hub) GetSubscriberCount(noradID int) int {
	h.mu.RLock()
	defer h.mu.RUnlock()
	if clients, ok := h.subscriptions[noradID]; ok {
		return len(clients)
	}
	return 0
}
