package services

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/jake-schuler/fll-event-screen/models"
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

func HandleWebSocketConnection(conn *websocket.Conn) {
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
		switch wsMessage.Type {
		case "get_info":
			// Parse the schedule round from the payload
			var round scheduleRound
			if payload, ok := wsMessage.Payload.(string); ok {
				switch payload {
				case "practice":
					round = practice
				case "round1":
					round = round1
				case "round2":
					round = round2
				case "round3":
					round = round3
				default:
					log.Printf("Unknown schedule round: %s, defaulting to practice", payload)
					round = practice // default to practice if unknown
				}
			} else {
				log.Printf("Payload is not a string, defaulting to practice. Payload: %v", wsMessage.Payload)
				round = practice // default if payload is not a string
			}

			matches := ReadSchedule(round)

			Manager.Broadcast(models.WebSocketMessage{
				Type:    "matches",
				Payload: matches,
			})
		case "set_active_match":
			Manager.Broadcast(wsMessage)
		case "show_timer":
			Manager.Broadcast(wsMessage)
		case "start_timer":
			Manager.Broadcast(models.WebSocketMessage{
				Type:    "start_timer",
				Payload: time.Now().Add(time.Second * 149),
			})
		case "play_test_sound":
			Manager.Broadcast(wsMessage)
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
