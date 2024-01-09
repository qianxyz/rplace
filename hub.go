package main

type Hub struct {
	// The counter.
	counter int

	// The hashset of clients.
	clients map[*Client]bool

	// Incoming message from clients.
	// For the counter app, the only message is to increment.
	broadcast chan bool

	// Register requests from clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client
}

func newHub() *Hub {
	return &Hub{
		counter:    0,
		clients:    make(map[*Client]bool),
		broadcast:  make(chan bool),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

func (h *Hub) run() {
	for {
		select {
		case c := <-h.register:
			h.clients[c] = true
			c.send <- h.counter // first time loading
		case c := <-h.unregister:
			if _, ok := h.clients[c]; ok {
				delete(h.clients, c)
				close(c.send)
			}
		case <-h.broadcast:
			h.counter++
			for c := range h.clients {
				select {
				case c.send <- h.counter:
				default: // client not listening, probably dead
					close(c.send)
					delete(h.clients, c)
				}
			}
		}
	}
}
