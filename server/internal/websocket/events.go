package websocket

import "encoding/json"

// EventType defines the type of WebSocket event.
type EventType string

const (
	EventStockUpdate EventType = "stock_update"
	EventOrderPlaced EventType = "order_placed"
)

// Event represents a WebSocket event payload.
type Event struct {
	Type    EventType   `json:"type"`
	Payload interface{} `json:"payload"`
}

// Marshal serializes the event to JSON bytes.
func (e Event) Marshal() ([]byte, error) {
	return json.Marshal(e)
}
