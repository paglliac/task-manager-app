package platform

import (
	"encoding/json"
	"log"
)

type Hub struct {
	clients    map[*Client]bool
	broadcast  chan WsEvent
	register   chan *Client
	unregister chan *Client
}

func (h *Hub) Handle(e WsEvent) {
	h.broadcast <- e
}

type WsEvent struct {
	Type  string      `json:"type"`
	Event interface{} `json:"event"`
}

func (we *WsEvent) toByte() []byte {
	bytes, err := json.Marshal(we)

	if err != nil {
		log.Printf("[web socket event to byte] error while encode ws event to bytes: %v", err)
	}

	return bytes
}

var hub *Hub

func InitHub() *Hub {
	hub = &Hub{
		clients:    make(map[*Client]bool),
		broadcast:  make(chan WsEvent),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}

	go hub.run()
	return hub
}

func (h *Hub) run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				close(client.send)
				delete(h.clients, client)
			}
		case message := <-h.broadcast:
			for client := range h.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}
		}
	}
}
