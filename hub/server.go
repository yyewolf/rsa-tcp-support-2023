package hub

import "net"

// Server is a hub server, it knows all the entities in the system.
type Server struct {
	Clients map[string]*Client
	Agents  map[string][]*Agent
}

// Client is a client connected to the hub via TCP.
type Client struct {
	ID string

	// Level should be 0, 1 or 2. (robots, employees or managers)
	Level int
	Conn  net.Conn
}

// Agent is an agent connected to the hub via TCP.
type Agent struct {
	ID    string
	Level int
	Conn  net.Conn
}

// NewServer creates a new server.
func NewServer() *Server {
	return &Server{
		Clients: make(map[string]*Client),
		Agents:  make(map[string][]*Agent),
	}
}

// Listen starts the server.
func (s *Server) Listen() error {
	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		return err
	}

	for {
		conn, err := ln.Accept()
		if err != nil {
			return err
		}

		go s.handleConn(conn)
	}
}

// handleConn handles a connection to the hub.
func (s *Server) handleConn(conn net.Conn) {
	// When a client connect we send an identify packet containing the client ID.
	// The client will then send a packet with authentication if it's an agent and a level.

}
