package chat

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/websocket"
)

// ServeWS обрабатывает WebSocket /api/chat/ws. Требует заголовок Authorization (JWT).
func (h *Hub) ServeWS(getUserID func(*http.Request) int64, getFullName func(*http.Request) string) http.HandlerFunc {
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
			Hub:      h,
			Conn:     conn,
			Send:     make(chan []byte, 256),
			UserID:   userID,
			FullName: getFullName(r),
		}
		client.Hub.register <- client
		go client.writePump()
		client.readPump()
	}
}

func (c *Client) readPump() {
	defer func() {
		c.Hub.unregister <- c
		c.Conn.Close()
	}()
	for {
		_, msg, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("chat read: %v", err)
			}
			break
		}
		var toSend []byte = msg
		var payload struct {
			Type string `json:"type"`
			Body string `json:"body"`
		}
		if json.Unmarshal(msg, &payload) == nil && payload.Type == "message" && payload.Body != "" {
			if c.Hub.SaveMessage != nil {
				if saved, err := c.Hub.SaveMessage(c.UserID, c.FullName, payload.Body); err == nil && len(saved) > 0 {
					toSend = saved
				}
			}
		}
		c.Hub.Broadcast(toSend)
	}
}

func (c *Client) writePump() {
	defer c.Conn.Close()
	for msg := range c.Send {
		if err := c.Conn.WriteMessage(websocket.TextMessage, msg); err != nil {
			break
		}
	}
}

// GetUserIDFromHeader извлекает user_id из заголовка X-User-ID.
func GetUserIDFromHeader(r *http.Request) int64 {
	s := r.Header.Get("X-User-ID")
	if s == "" {
		return 0
	}
	id, _ := strconv.ParseInt(s, 10, 64)
	return id
}

// GetFullNameFromHeader извлекает full_name из заголовка X-User-Name.
func GetFullNameFromHeader(r *http.Request) string {
	return r.Header.Get("X-User-Name")
}
