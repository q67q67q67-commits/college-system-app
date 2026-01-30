package chat

import (
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

// Client — подключённый клиент чата.
type Client struct {
	Hub      *Hub
	Conn     *websocket.Conn
	Send     chan []byte
	UserID   int64
	FullName string
}

// SaveMessageFunc сохраняет сообщение в БД и возвращает JSON для рассылки. nil = не сохранять.
type SaveMessageFunc func(userID int64, fullName, body string) (broadcastJSON []byte, err error)

// Hub — центр рассылки сообщений (общий чат).
type Hub struct {
	clients      map[*Client]bool
	broadcast    chan []byte
	register     chan *Client
	unregister   chan *Client
	SaveMessage  SaveMessageFunc
	mu           sync.RWMutex
}

// NewHub создаёт новый Hub.
func NewHub() *Hub {
	return &Hub{
		clients:    make(map[*Client]bool),
		broadcast:  make(chan []byte, 256),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

// Run запускает цикл рассылки.
func (h *Hub) Run() {
	for {
		select {
		case c := <-h.register:
			h.mu.Lock()
			h.clients[c] = true
			h.mu.Unlock()
		case c := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.clients[c]; ok {
				delete(h.clients, c)
				close(c.Send)
			}
			h.mu.Unlock()
		case msg := <-h.broadcast:
			h.mu.RLock()
			for c := range h.clients {
				select {
				case c.Send <- msg:
				default:
					close(c.Send)
					delete(h.clients, c)
				}
			}
			h.mu.RUnlock()
		}
	}
}

// Broadcast отправляет сообщение всем клиентам.
func (h *Hub) Broadcast(msg []byte) {
	h.broadcast <- msg
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}
