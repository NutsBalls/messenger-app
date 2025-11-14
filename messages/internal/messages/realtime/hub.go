package realtime

import (
	"sync"
)

type Broadcast struct {
	ChatID string
	Data   []byte
}

type Hub struct {
	clients    map[*Client]bool
	chatSubs   map[string]map[*Client]bool
	register   chan *Client
	unregister chan *Client
	broadcast  chan Broadcast
	mu         sync.RWMutex
}

func NewHub() *Hub {
	return &Hub{
		clients:    make(map[*Client]bool),
		chatSubs:   make(map[string]map[*Client]bool),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		broadcast:  make(chan Broadcast, 256),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case c := <-h.register:
			h.mu.Lock()
			h.clients[c] = true
			h.mu.Unlock()

		case c := <-h.unregister:
			h.mu.Lock()
			delete(h.clients, c)
			for chatID := range h.chatSubs {
				delete(h.chatSubs[chatID], c)
			}
			h.mu.Unlock()
			close(c.send)

		case msg := <-h.broadcast:
			h.mu.RLock()
			subs := h.chatSubs[msg.ChatID]
			for client := range subs {
				select {
				case client.send <- msg.Data:
				default:
				}
			}
			h.mu.RUnlock()
		}
	}
}

func (h *Hub) Subscribe(client *Client, chatID string) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if h.chatSubs[chatID] == nil {
		h.chatSubs[chatID] = make(map[*Client]bool)
	}
	h.chatSubs[chatID][client] = true
}

func (h *Hub) Unsubscribe(client *Client, chatID string) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if subs, ok := h.chatSubs[chatID]; ok {
		delete(subs, client)
		if len(subs) == 0 {
			delete(h.chatSubs, chatID)
		}
	}
}
