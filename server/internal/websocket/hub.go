package websocket

import (
	"sync"

	"github.com/google/uuid"
)

// Hub manages WebSocket connections grouped by merchant (per-merchant rooms).
type Hub struct {
	mu         sync.RWMutex
	rooms      map[uuid.UUID]map[*Client]bool
	register   chan *Client
	unregister chan *Client
}

// NewHub creates a new WebSocket Hub.
func NewHub() *Hub {
	h := &Hub{
		rooms:      make(map[uuid.UUID]map[*Client]bool),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
	go h.run()
	return h
}

func (h *Hub) run() {
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			if h.rooms[client.MerchantID] == nil {
				h.rooms[client.MerchantID] = make(map[*Client]bool)
			}
			h.rooms[client.MerchantID][client] = true
			h.mu.Unlock()

		case client := <-h.unregister:
			h.mu.Lock()
			if room, ok := h.rooms[client.MerchantID]; ok {
				if _, exists := room[client]; exists {
					delete(room, client)
					close(client.Send)
					if len(room) == 0 {
						delete(h.rooms, client.MerchantID)
					}
				}
			}
			h.mu.Unlock()
		}
	}
}

// Register adds a client to its merchant room.
func (h *Hub) Register(client *Client) {
	h.register <- client
}

// Unregister removes a client from its merchant room.
func (h *Hub) Unregister(client *Client) {
	h.unregister <- client
}

// Broadcast sends an event to all connected clients in a merchant room.
func (h *Hub) Broadcast(merchantID uuid.UUID, event Event) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	room, ok := h.rooms[merchantID]
	if !ok {
		return
	}

	data, err := event.Marshal()
	if err != nil {
		return
	}

	for client := range room {
		select {
		case client.Send <- data:
		default:
			// Client buffer full, drop message
		}
	}
}
