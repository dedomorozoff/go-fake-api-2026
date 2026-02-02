package models

import (
	"time"
)

// Board представляет интерактивную доску
type Board struct {
	ID        string                 `json:"id"`
	Hash      string                 `json:"hash"` // Публичный хеш для доступа без авторизации
	Name      string                 `json:"name"`
	OwnerID   int                    `json:"owner_id"`
	IsPublic  bool                   `json:"is_public"`
	Likes     int                    `json:"likes"`
	Objects   map[string]BoardObject   `json:"objects"` // map[object_id]Object
	CreatedAt time.Time              `json:"created_at"`
}

// BoardObject представляет объект на доске
type BoardObject struct {
	ID         string  `json:"id"`
	Type       string  `json:"type"` // text, image, rectangle, circle, line
	X          float64 `json:"x"`
	Y          float64 `json:"y"`
	Width      float64 `json:"width"`
	Height     float64 `json:"height"`
	Rotation   float64 `json:"rotation"`
	Content    string  `json:"content,omitempty"` // Текст или URL изображения
	Color      string  `json:"color,omitempty"`
	FocusedBy  *int    `json:"focused_by,omitempty"`   // ID пользователя, захватившего объект
	FocusedAt  *time.Time `json:"focused_at,omitempty"`
	OwnerName  string  `json:"owner_name,omitempty"`   // Имя пользователя, захватившего объект
}

// BoardAccess представляет права доступа к доске
type BoardAccess struct {
	BoardID string `json:"board_id"`
	UserID  int    `json:"user_id"`
	Email   string `json:"email"`
}

// Like представляет лайк доске
type Like struct {
	BoardID string `json:"board_id"`
	UserID  int    `json:"user_id"`
}

// BoardCreateRequest запрос на создание доски
type BoardCreateRequest struct {
	Name     string `json:"name"`
	IsPublic bool   `json:"is_public"`
}

// BoardShareRequest запрос на предоставление доступа
type BoardShareRequest struct {
	Email string `json:"email"`
}

// WSMessage структура сообщения WebSocket
type WSMessage struct {
	Type    string      `json:"type"`    // object_update, object_focus, object_blur, object_delete
	BoardID string      `json:"board_id"`
	Payload interface{} `json:"payload"`
}
