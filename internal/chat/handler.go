package chat

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/websocket"
)

// ServeWS обрабатывает WebSocket /api/chat/ws. Требует заголовок Authorization (JWT).
// После апгрейда клиент отправляет/получает JSON: {"type":"message","body":"..."}.
func (h *Hub) ServeWS(getUserID func(*http.Request) int64) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := getUserID(r)
		if userID == 0 {
			http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
			return
		}
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Printf("chat ws upgrade: %v", err)
			return
		}
		client := &Client{
			hub:  h,
			conn: conn,
			send: make(chan []byte, 256),
		}
		client.hub.register <- client
		go client.writePump()
		client.readPump()
	}
}

func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()
	for {
		_, msg, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("chat read: %v", err)
			}
			break
		}
		// Broadcast to all (в т.ч. отправителю)
		c.hub.Broadcast(msg)
	}
}

func (c *Client) writePump() {
	defer c.conn.Close()
	for msg := range c.send {
		if err := c.conn.WriteMessage(websocket.TextMessage, msg); err != nil {
			break
		}
	}
}

// GetUserIDFromHeader извлекает user_id из заголовка X-User-ID (после JWT middleware).
func GetUserIDFromHeader(r *http.Request) int64 {
	s := r.Header.Get("X-User-ID")
	if s == "" {
		return 0
	}
	id, _ := strconv.ParseInt(s, 10, 64)
	return id
}
