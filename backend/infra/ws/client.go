package ws

import (
	"github.com/gorilla/websocket"
)

type Client struct {
	UserID int
	Conn   *websocket.Conn
	Send   chan []byte // per-client write channel
}

// WritePump writes messages from the Send channel to the WebSocket
func (c *Client) WritePump() {
	for msg := range c.Send {
		err := c.Conn.WriteMessage(websocket.TextMessage, msg)
		if err != nil {
			break
		}
	}
	c.Conn.Close()
}