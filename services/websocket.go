package services

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"

	"github.com/jake-schuler/fll-event-screen/models"
	"github.com/gorilla/websocket"
	"gorm.io/gorm"
)

// Global state variables to persist WebSocket state
// var current_match_state *models.WebSocketMatchPayload

var Upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		log.Printf("WebSocket upgrade check origin: %s", r.Header.Get("Origin"))
		return true // Allow all origins for simplicity; adjust as needed
	},
}

func HandleWebSocketConnection(conn *websocket.Conn, db *gorm.DB) {
	defer func() {
		Manager.RemoveConnection(conn)
		conn.Close()
	}()

	log.Println("WebSocket connection established")

	// Add connection to manager
	Manager.AddConnection(conn)

	// Handle WebSocket messages in a loop
	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Printf("WebSocket read error: %v", err)
			break
		}

		log.Printf("Received message: %s", message)

		// Parse JSON message
		var wsMessage models.WebSocketMessage
		err = json.Unmarshal(message, &wsMessage)
		if err != nil {
			log.Printf("Error parsing WebSocket message: %v", err)
			continue
		}
	}
	log.Println("WebSocket connection closed")
}

// WebSocket connection manager
type ConnectionManager struct {
	connections map[*websocket.Conn]bool
	mutex       sync.RWMutex
}

var Manager = &ConnectionManager{
	connections: make(map[*websocket.Conn]bool),
}

// Add connection to manager
func (cm *ConnectionManager) AddConnection(conn *websocket.Conn) {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()
	cm.connections[conn] = true
	log.Printf("WebSocket connection added. Total connections: %d", len(cm.connections))
}

// Remove connection from manager
func (cm *ConnectionManager) RemoveConnection(conn *websocket.Conn) {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()
	delete(cm.connections, conn)
	log.Printf("WebSocket connection removed. Total connections: %d", len(cm.connections))
}

// Broadcast message to all connected clients
func (cm *ConnectionManager) Broadcast(message models.WebSocketMessage) {
	cm.mutex.RLock()
	defer cm.mutex.RUnlock()

	for conn := range cm.connections {
		err := conn.WriteJSON(message)
		if err != nil {
			log.Printf("Error broadcasting to client: %v", err)
			// Remove dead connection
			go func(c *websocket.Conn) {
				cm.RemoveConnection(c)
				c.Close()
			}(conn)
		}
	}
}
