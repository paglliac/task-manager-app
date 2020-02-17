package platform

type Hub struct {
	clients    map[*Client]bool
	Broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
}

var hub *Hub

func InitHub() *Hub {
	hub = &Hub{
		clients:    make(map[*Client]bool),
		Broadcast:  make(chan []byte),
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
		case message := <-h.Broadcast:
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
