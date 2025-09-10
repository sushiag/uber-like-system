package ws

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

type Connection struct {
	ID     uint64
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

	// auto-increment for assigning IDs
	nextID uint64
}

// creates a new WebSocket manager
func NewWebSocketManager() *WebSocketManager {
	manager := &WebSocketManager{
		Connections: make(map[uint64]*Connection),
		register:    make(chan *Connection),
		unregister:  make(chan *Connection),
		broadcast:   make(chan []byte),
		nextID:      1,
	}

	go manager.run()
	return manager
}

// starts the event loop for handling connections
func (m *WebSocketManager) run() {
	for {
		select {
		case conn := <-m.register:
			m.Connections[conn.ID] = conn
			log.Printf("Client %d connected", conn.ID)

		case conn := <-m.unregister:
			if _, ok := m.Connections[conn.ID]; ok {
				close(conn.Send)
				delete(m.Connections, conn.ID)
				log.Printf("Client %d disconnected", conn.ID)
			}

		case msg := <-m.broadcast:
			for _, conn := range m.Connections {
				select {
				case conn.Send <- msg:
				default:
					close(conn.Send)
					delete(m.Connections, conn.ID)
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

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("WebSocket upgrade failed:", err)
		return
	}

	// client connection
	client := &Connection{
		ID:     m.nextID,
		Conn:   conn,
		Send:   make(chan []byte, 256),
		Closed: make(chan struct{}),
	}
	m.nextID++

	m.register <- client

	go m.readPump(client)
	go m.writePump(client)
}

func (m *WebSocketManager) readPump(c *Connection) {
	defer func() {
		m.unregister <- c
		c.Conn.Close()
	}()

	for {
		_, msg, err := c.Conn.ReadMessage()
		if err != nil {
			log.Printf("Read error from client %d: %v", c.ID, err)
			break
		}

		m.broadcast <- msg
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
