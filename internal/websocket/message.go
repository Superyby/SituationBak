package websocket

import (
	"encoding/json"
	"time"
)

// MessageType 消息类型
type MessageType string

const (
	// 客户�?-> 服务�?
	MessageTypePing        MessageType = "ping"
	MessageTypeSubscribe   MessageType = "subscribe"
	MessageTypeUnsubscribe MessageType = "unsubscribe"

	// 服务�?-> 客户�?
	MessageTypePong            MessageType = "pong"
	MessageTypeSatelliteUpdate MessageType = "satellite_update"
	MessageTypeNotification    MessageType = "notification"
	MessageTypeError           MessageType = "error"
)

// Message WebSocket消息结构
type Message struct {
	Type      MessageType     `json:"type"`
	Payload   json.RawMessage `json:"payload,omitempty"`
	Timestamp time.Time       `json:"timestamp"`
}

// SubscribePayload 订阅消息负载
type SubscribePayload struct {
	NoradIDs []int `json:"norad_ids"`
}

// UnsubscribePayload 取消订阅消息负载
type UnsubscribePayload struct {
	NoradIDs []int `json:"norad_ids"`
}

// SatellitePosition 卫星位置信息
type SatellitePosition struct {
	NoradID  int     `json:"norad_id"`
	Name     string  `json:"name"`
	Position Vector3 `json:"position"`
	Velocity Vector3 `json:"velocity"`
}

// Vector3 三维向量
type Vector3 struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
	Z float64 `json:"z"`
}

// SatelliteUpdatePayload 卫星更新消息负载
type SatelliteUpdatePayload struct {
	Satellites []SatellitePosition `json:"satellites"`
}

// NotificationPayload 通知消息负载
type NotificationPayload struct {
	Title   string `json:"title"`
	Message string `json:"message"`
	Level   string `json:"level"` // info, warning, error
}

// ErrorPayload 错误消息负载
type ErrorPayload struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// NewMessage 创建新消�?
func NewMessage(msgType MessageType, payload interface{}) (*Message, error) {
	var payloadBytes json.RawMessage
	if payload != nil {
		bytes, err := json.Marshal(payload)
		if err != nil {
			return nil, err
		}
		payloadBytes = bytes
	}

	return &Message{
		Type:      msgType,
		Payload:   payloadBytes,
		Timestamp: time.Now().UTC(),
	}, nil
}

// ParsePayload 解析消息负载
func (m *Message) ParsePayload(v interface{}) error {
	if m.Payload == nil {
		return nil
	}
	return json.Unmarshal(m.Payload, v)
}

// ToJSON 转换为JSON
func (m *Message) ToJSON() ([]byte, error) {
	return json.Marshal(m)
}
