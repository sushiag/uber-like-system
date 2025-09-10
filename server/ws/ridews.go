package ws

func (m *WebSocketManager) SendToUser(userID uint64, message []byte) {
	if conn, ok := m.Connections[userID]; ok {
		select {
		case conn.Send <- message:
		default:
			close(conn.Send)
			delete(m.Connections, userID)
		}
	}
}
