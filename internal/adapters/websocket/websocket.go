package websocket

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"
	applog "web-porto-backend/common/logger"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

// Manager manages WebSocket connections
type Manager struct {
	clients      map[*Client]bool
	broadcast    chan []byte
	register     chan *Client
	unregister   chan *Client
	mutex        sync.Mutex
	latestCounts map[string]interface{}
}

// Client is a middleman between the WebSocket connection and the Manager
type Client struct {
	manager *Manager
	conn    *websocket.Conn
	send    chan []byte
}

// Message represents the structure of messages sent over WebSocket
type Message struct {
	Type    string      `json:"type"`
	Data    interface{} `json:"data"`
	Channel string      `json:"channel,omitempty"`
}

// ViewCountsUpdate represents analytics data to be sent to clients
type ViewCountsUpdate struct {
	Total  int64  `json:"total"`
	Today  int64  `json:"today"`
	Week   int64  `json:"week"`
	Month  int64  `json:"month"`
	Unique int64  `json:"unique"`
	Page   string `json:"page,omitempty"`
}

// WebSocket connection upgrader with CORS support
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	// Allow all origins for development
	CheckOrigin: func(r *http.Request) bool {
		// Allow all origins in development mode, in production you might want to be more restrictive
		return true
	},
}

// NewManager creates a new WebSocket manager
func NewManager() *Manager {
	return &Manager{
		clients:      make(map[*Client]bool),
		broadcast:    make(chan []byte),
		register:     make(chan *Client),
		unregister:   make(chan *Client),
		latestCounts: make(map[string]interface{}),
	}
}

// Start begins listening for WebSocket events
func (m *Manager) Start() {
	for {
		select {
		case client := <-m.register:
			m.mutex.Lock()
			m.clients[client] = true
			m.mutex.Unlock()

			// Send latest counts immediately upon connection
			if len(m.latestCounts) > 0 {
				for channel, counts := range m.latestCounts {
					msg := Message{
						Type:    "view_counts",
						Data:    counts,
						Channel: channel,
					}
					data, err := json.Marshal(msg)
					if err == nil {
						client.send <- data
					}
				}
			}

		case client := <-m.unregister:
			m.mutex.Lock()
			if _, ok := m.clients[client]; ok {
				close(client.send)
				delete(m.clients, client)
			}
			m.mutex.Unlock()

		case message := <-m.broadcast:
			m.mutex.Lock()
			for client := range m.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(m.clients, client)
				}
			}
			m.mutex.Unlock()
		}
	}
}

// UpdateViewCounts broadcasts view count updates to all connected clients
func (m *Manager) UpdateViewCounts(counts ViewCountsUpdate, page string) {
	channel := "global"
	if page != "" {
		channel = "page:" + page
	}

	// Store latest counts for new connections
	m.mutex.Lock()
	m.latestCounts[channel] = counts
	m.mutex.Unlock()

	// Prepare message
	message := Message{
		Type:    "view_counts",
		Data:    counts,
		Channel: channel,
	}

	data, err := json.Marshal(message)
	if err != nil {
		applog.GetLogger().Error("Failed to marshal view counts update", applog.Fields{"error": err.Error()})
		return
	}

	m.broadcast <- data
}

// ServeWs handles WebSocket requests from clients
func (m *Manager) ServeWs(c *gin.Context) {
	log := applog.GetLogger().WithFields(applog.Fields{"handler": "websocket.ServeWs"})

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Error("Failed to upgrade connection", applog.Fields{"error": err.Error()})
		return
	}

	client := &Client{
		manager: m,
		conn:    conn,
		send:    make(chan []byte, 256),
	}

	m.register <- client

	// Start goroutines for reading and writing
	go client.writePump()
	go client.readPump()
}

// writePump pumps messages from the manager to the WebSocket connection
func (c *Client) writePump() {
	ticker := time.NewTicker(60 * time.Second)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			if !ok {
				// The manager closed the channel
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			if err := w.Close(); err != nil {
				return
			}

		case <-ticker.C:
			// Send ping to prevent connection from timing out
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// readPump pumps messages from the WebSocket connection to the manager
func (c *Client) readPump() {
	defer func() {
		c.manager.unregister <- c
		c.conn.Close()
	}()

	c.conn.SetReadLimit(512) // Limit message size
	c.conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil
	})

	for {
		_, _, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}
		// For now we don't process incoming messages
		// but we could handle client requests here
	}
}
