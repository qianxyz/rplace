package main

type Server struct {
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

func newServer() *Server {
	return &Server{
		counter:    0,
		clients:    make(map[*Client]bool),
		broadcast:  make(chan bool),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

func (s *Server) run() {
	for {
		select {
		case c := <-s.register:
			s.clients[c] = true
			c.send <- s.counter // first time loading
		case c := <-s.unregister:
			if _, ok := s.clients[c]; ok {
				delete(s.clients, c)
				close(c.send)
			}
		case <-s.broadcast:
			s.counter++
			for c := range s.clients {
				select {
				case c.send <- s.counter:
				default: // client not listening, probably dead
					close(c.send)
					delete(s.clients, c)
				}
			}
		}
	}
}
