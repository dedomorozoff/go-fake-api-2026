package api

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/alexl/go-fake-api/internal/models"
	"github.com/alexl/go-fake-api/internal/storage"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // Разрешаем всем (для разработки)
	},
}

// Client представляет подключенного пользователя
type Client struct {
	Hub     *Hub
	Conn    *websocket.Conn
	Send    chan []byte
	UserID  int
	UserName string
	BoardID string
}

// Hub управляет всеми подключениями
type Hub struct {
	clients    map[string]map[*Client]bool // boardID -> clients
	broadcast  chan models.WSMessage
	register   chan *Client
	unregister chan *Client
	storage    storage.Storage
	mu         sync.Mutex
}

func NewHub(s storage.Storage) *Hub {
	return &Hub{
		broadcast:  make(chan models.WSMessage),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[string]map[*Client]bool),
		storage:    s,
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			if h.clients[client.BoardID] == nil {
				h.clients[client.BoardID] = make(map[*Client]bool)
			}
			h.clients[client.BoardID][client] = true
			h.mu.Unlock()

		case client := <-h.unregister:
			h.mu.Lock()
			if clients, ok := h.clients[client.BoardID]; ok {
				if _, ok := clients[client]; ok {
					delete(clients, client)
					close(client.Send)
					if len(clients) == 0 {
						delete(h.clients, client.BoardID)
					}
				}
			}
			h.mu.Unlock()

		case message := <-h.broadcast:
			h.mu.Lock()
			msgBytes, _ := json.Marshal(message)
			for client := range h.clients[message.BoardID] {
				select {
				case client.Send <- msgBytes:
				default:
					close(client.Send)
					delete(h.clients[message.BoardID], client)
				}
			}
			h.mu.Unlock()
		}
	}
}

func (c *Client) ReadPump() {
	defer func() {
		c.Hub.unregister <- c
		c.Conn.Close()
	}()

	for {
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			break
		}

		var wsMsg models.WSMessage
		if err := json.Unmarshal(message, &wsMsg); err != nil {
			continue
		}

		wsMsg.BoardID = c.BoardID // Принудительно ставим BoardID клиента

		// Обработка разных типов сообщений
		switch wsMsg.Type {
		case "object_update":
			var obj models.BoardObject
			payloadBytes, _ := json.Marshal(wsMsg.Payload)
			json.Unmarshal(payloadBytes, &obj)
			
			// Проверяем фокус
			board, _ := c.Hub.storage.GetBoardByID(c.BoardID)
			if existing, ok := board.Objects[obj.ID]; ok {
				if existing.FocusedBy != nil && *existing.FocusedBy != c.UserID {
					continue // Нельзя менять чужой фокус
				}
			}
			
			c.Hub.storage.UpdateBoardObject(c.BoardID, obj)
			c.Hub.broadcast <- wsMsg

		case "object_focus":
			objectID := wsMsg.Payload.(string)
			now := time.Now()
			obj := models.BoardObject{
				ID:        objectID,
				FocusedBy: &c.UserID,
				FocusedAt: &now,
				OwnerName: c.UserName,
			}
			// Получаем текущий объект, чтобы не затереть данные
			board, _ := c.Hub.storage.GetBoardByID(c.BoardID)
			if existing, ok := board.Objects[objectID]; ok {
				obj.Type = existing.Type
				obj.X = existing.X
				obj.Y = existing.Y
				obj.Width = existing.Width
				obj.Height = existing.Height
				obj.Rotation = existing.Rotation
				obj.Content = existing.Content
				obj.Color = existing.Color
			}
			
			c.Hub.storage.UpdateBoardObject(c.BoardID, obj)
			c.Hub.broadcast <- wsMsg

		case "object_blur":
			objectID := wsMsg.Payload.(string)
			board, _ := c.Hub.storage.GetBoardByID(c.BoardID)
			if obj, ok := board.Objects[objectID]; ok {
				if obj.FocusedBy != nil && *obj.FocusedBy == c.UserID {
					obj.FocusedBy = nil
					obj.FocusedAt = nil
					obj.OwnerName = ""
					c.Hub.storage.UpdateBoardObject(c.BoardID, obj)
					
					// Рассылаем обновление о снятии фокуса
					wsMsg.Payload = obj
					c.Hub.broadcast <- wsMsg
				}
			}

		case "object_delete":
			objectID := wsMsg.Payload.(string)
			c.Hub.storage.DeleteBoardObject(c.BoardID, objectID)
			c.Hub.broadcast <- wsMsg
		}
	}
}

func (c *Client) WritePump() {
	for {
		select {
		case message, ok := <-c.Send:
			if !ok {
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			c.Conn.WriteMessage(websocket.TextMessage, message)
		}
	}
}

func ServeWs(hub *Hub, s storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		boardID := vars["board_id"]
		
		token := r.URL.Query().Get("token")
		user, err := s.GetUserByToken(token)
		if err != nil {
			// Для публичного просмотра тоже можно разрешить WS, 
			// но без права редактирования. 
			// Пока сделаем только для авторизованных.
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		hasAccess, _ := s.HasBoardAccess(boardID, user.ID)
		if !hasAccess {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}

		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println(err)
			return
		}

		client := &Client{
			Hub:      hub,
			Conn:     conn,
			Send:     make(chan []byte, 256),
			UserID:   user.ID,
			UserName: user.Name,
			BoardID:  boardID,
		}

		client.Hub.register <- client

		go client.WritePump()
		go client.ReadPump()
	}
}
