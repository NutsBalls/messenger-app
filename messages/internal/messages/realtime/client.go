package realtime

import (
	"encoding/json"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

type Client struct {
	UserID string
	Conn   *websocket.Conn
	send   chan []byte
	Hub    *Hub
}

type MessageFromClient struct {
	Type   string          `json:"type"`
	ChatID string          `json:"chat_id"`
	Data   json.RawMessage `json:"data,omitempty"`
}

type SendMessageToClient struct {
	Type   string      `json:"type"`
	ChatID string      `json:"chat_id"`
	Data   interface{} `json:"data"`
}

func (c *Client) ReadLoop() {
	defer func() {
		c.Hub.unregister <- c
		_ = c.Conn.Close()
	}()

	c.Conn.SetReadLimit(5120)
	c.Conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	c.Conn.SetPongHandler(func(string) error {
		c.Conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil
	})

	for {
		_, msg, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("ws: unexpected close: %v", err)
			}
			break
		}

		var m MessageFromClient
		if err := json.Unmarshal(msg, &m); err != nil {
			continue
		}

		switch m.Type {
		case "subscribe":
			if m.ChatID != "" {
				c.Hub.Subscribe(c, m.ChatID)
			}
		case "unsubscribe":
			if m.ChatID != "" {
				c.Hub.Unsubscribe(c, m.ChatID)
			}
		case "ping":
			_ = c.Conn.WriteMessage(websocket.PongMessage, nil)
		default:
			// другие типы сообщений — например, client->server chat actions
		}
	}
}

func (c *Client) WriteLoop() {
	ticker := time.NewTicker(30 * time.Second)
	defer func() {
		ticker.Stop()
		_ = c.Conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			_ = c.Conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if !ok {
				_ = c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.Conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			if _, err := w.Write(message); err != nil {
				_ = w.Close()
				return
			}
			if err := w.Close(); err != nil {
				return
			}

		case <-ticker.C:
			_ = c.Conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
