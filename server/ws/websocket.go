package ws

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/websocket"
)

type Connection struct {
	UserID uint64
	Conn   *websocket.Conn
	Send   chan []byte
	Closed chan struct{}
}

// all client connections and messages
type WebSocketManager struct {
	Connections map[uint64]*Connection
	register    chan *Connection
	unregister  chan *Connection
	broadcast   chan []byte
}

// creates a new WebSocket manager
func NewWebSocketManager() *WebSocketManager {
	manager := &WebSocketManager{
		Connections: make(map[uint64]*Connection),
		register:    make(chan *Connection),
		unregister:  make(chan *Connection),
		broadcast:   make(chan []byte),
	}

	go manager.run()
	return manager
}

// starts the event loop for handling connections
func (m *WebSocketManager) run() {
	for {
		select {
		case conn := <-m.register:
			m.Connections[conn.UserID] = conn
			log.Printf("User %d connected", conn.UserID)

		case conn := <-m.unregister:
			if _, ok := m.Connections[conn.UserID]; ok {
				close(conn.Send)
				delete(m.Connections, conn.UserID)
				log.Printf("Client %d disconnected", conn.UserID)
			}

		case msg := <-m.broadcast:
			for _, conn := range m.Connections {
				select {
				case conn.Send <- msg:
				default:
					close(conn.Send)
					delete(m.Connections, conn.UserID)
				}
			}
		}
	}
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// new client connections
func (m *WebSocketManager) WebSocketHandler(w http.ResponseWriter, r *http.Request) {

	userIDStr := r.URL.Query().Get("user_id")
	if userIDStr == "" {
		http.Error(w, "userID is required to connect", http.StatusBadRequest)
		return
	}

	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		http.Error(w, "user is invalid", http.StatusUnauthorized)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("WebSocket upgrade failed:", err)
		return
	}

	// client connection
	client := &Connection{
		UserID: uint64(userID),
		Conn:   conn,
		Send:   make(chan []byte, 256),
		Closed: make(chan struct{}),
	}

	m.register <- client

	go m.readPump(client)
	go m.writePump(client)
}

type RideEvent struct {
	Event  string `json:"event"`
	RideID int64  `json:"ride_id"`
	FromID uint64 `json:"from_id"`
}

func (m *WebSocketManager) readPump(c *Connection) {
	defer func() {
		m.unregister <- c
		c.Conn.Close()
	}()

	for {
		_, msg, err := c.Conn.ReadMessage()
		if err != nil {
			break
		}

		var event RideEvent
		if err := json.Unmarshal(msg, &event); err == nil {
			// Example: notify rider that driver accepted
			if event.Event == "ride_accepted" {
				m.SendToUser(event.FromID, msg)
			}
		}
	}
}

// utgoing messages to the client
func (m *WebSocketManager) writePump(c *Connection) {
	defer c.Conn.Close()

	for {
		select {
		case msg, ok := <-c.Send:
			if !ok {
				// Connection closed
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			c.Conn.WriteMessage(websocket.TextMessage, msg)

		case <-c.Closed:
			return
		}
	}
}
