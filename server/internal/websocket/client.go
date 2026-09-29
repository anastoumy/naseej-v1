package websocket

import (
	"time"

	"github.com/google/uuid"
	gorilla "github.com/gorilla/websocket"
)

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 4096
)

// Client represents a single WebSocket connection.
type Client struct {
	Hub        *Hub
	Conn       *gorilla.Conn
	Send       chan []byte
	MerchantID uuid.UUID
}

// NewClient creates a new WebSocket client.
func NewClient(hub *Hub, conn *gorilla.Conn, merchantID uuid.UUID) *Client {
	return &Client{
		Hub:        hub,
		Conn:       conn,
		Send:       make(chan []byte, 256),
		MerchantID: merchantID,
	}
}

// ReadPump reads messages from the WebSocket connection (handles pong/close).
func (c *Client) ReadPump() {
	defer func() {
		c.Hub.Unregister(c)
		c.Conn.Close()
	}()

	c.Conn.SetReadLimit(maxMessageSize)
	c.Conn.SetReadDeadline(time.Now().Add(pongWait))
	c.Conn.SetPongHandler(func(string) error {
		c.Conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		_, _, err := c.Conn.ReadMessage()
		if err != nil {
			break
		}
	}
}

// WritePump sends messages to the WebSocket connection.
func (c *Client) WritePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.Conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.Send:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				c.Conn.WriteMessage(gorilla.CloseMessage, []byte{})
				return
			}

			w, err := c.Conn.NextWriter(gorilla.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)
			if err := w.Close(); err != nil {
				return
			}

		case <-ticker.C:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.Conn.WriteMessage(gorilla.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
