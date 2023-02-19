package hub

import (
	"encoding/json"
	"hub/internal/packets"
	"hub/internal/uid"
	"log"
	"net"
	"os"
	"sync"
)

// Server is a hub server, it knows all the entities in the system.
type Server struct {
	Clients map[int]*Client
	Agents  map[int][]*Agent
}

// Client is a client connected to the hub via TCP.
type Client struct {
	ID int

	// Level should be 0, 1 or 2. (robots, employees or managers)
	Level int
	Conn  net.Conn

	History []*packets.Message

	*sync.Mutex
}

func (c *Client) Send(p *packets.Packet) error {
	c.Lock()
	defer c.Unlock()

	_, err := c.Conn.Write(p.ToBytes())

	return err
}

// Agent is an agent connected to the hub via TCP.
type Agent struct {
	ID    int
	Name  string
	Level int
	Conn  net.Conn

	*sync.Mutex
}

func (a *Agent) Send(p *packets.Packet) error {
	a.Lock()
	defer a.Unlock()

	_, err := a.Conn.Write(p.ToBytes())

	if os.Getenv("DEBUG") == "true" {
		log.Println("[HUB] Sent packet to agent", a.ID, ":", p.Type, " => ", string(p.ToBytes()))
	}

	return err
}

type Connection struct {
	ID     int
	Agent  *Agent
	Client *Client
	Conn   net.Conn

	*sync.Mutex
}

func (c *Connection) Send(p *packets.Packet) error {
	if c.Agent != nil {
		return c.Agent.Send(p)
	}
	if c.Client != nil {
		return c.Client.Send(p)
	}
	c.Lock()
	defer c.Unlock()

	_, err := c.Conn.Write(p.ToBytes())
	return err
}

// NewServer creates a new server.
func NewServer() *Server {
	return &Server{
		Clients: make(map[int]*Client),
		Agents:  make(map[int][]*Agent),
	}
}

// Listen starts the server.
func (s *Server) Listen(port string) error {
	ln, err := net.Listen("tcp", port)
	if err != nil {
		return err
	}

	if os.Getenv("DEBUG") == "true" {
		log.Println("[HUB] Listening on port " + port)
	}

	for {
		conn, err := ln.Accept()
		if err != nil {
			return err
		}

		var c = &Connection{
			ID:    uid.New(),
			Conn:  conn,
			Mutex: &sync.Mutex{},
		}

		go func() {
			s.handleConn(c)
			conn.Close()

			if c.Agent != nil {
				for i, a := range s.Agents[c.Agent.Level] {
					if a == c.Agent {
						s.Agents[c.Agent.Level] = append(s.Agents[c.Agent.Level][:i], s.Agents[c.Agent.Level][i+1:]...)
					}
				}
			}

			if c.Client != nil {
				delete(s.Clients, c.Client.ID)
			}
		}()
	}
}

// handleConn handles a connection to the hub.
func (s *Server) handleConn(c *Connection) {
	if os.Getenv("DEBUG") == "true" {
		log.Println("[HUB] Received connection from " + c.Conn.RemoteAddr().String())
	}

	err := c.Send(packets.NewHello(c.ID))
	if err != nil {
		return
	}

	for {
		// Read the packet.
		buf := make([]byte, 1024)
		n, err := c.Conn.Read(buf)
		if err != nil {
			return
		}

		if os.Getenv("DEBUG") == "true" {
			log.Println("[HUB] Received packet from " + c.Conn.RemoteAddr().String())
		}

		// Unmarshal the packet.
		pkt := &packets.Packet{}
		err = json.Unmarshal(buf[:n], pkt)
		if err != nil {
			c.Send(packets.NewError(err.Error()))
			return
		}

		// Remarshal the packet data.
		pkt.Raw, err = json.Marshal(pkt.Data)
		if err != nil {
			c.Send(packets.NewError(err.Error()))
			return
		}

		if os.Getenv("DEBUG") == "true" {
			log.Println("[HUB] Packet type:", pkt.Type)
		}

		// Handle the packet.
		switch pkt.Type {
		case packets.PacketTypeIdentify:
			data := &packets.Identify{}
			err = json.Unmarshal(pkt.Raw, data)
			if err != nil {
				c.Send(packets.NewError(err.Error()))
				return
			}

			err = s.handleIdentify(data, c)
			if err != nil {
				c.Send(packets.NewError(err.Error()))
				return
			}
		case packets.PacketTypeMessage:
			data := &packets.Message{}
			err = json.Unmarshal(pkt.Raw, data)
			if err != nil {
				c.Send(packets.NewError(err.Error()))
				return
			}

			err = s.handleMessage(data, c)
			if err != nil {
				c.Send(packets.NewError(err.Error()))
				return
			}
		case packets.PacketTypeElevate:
			data := &packets.Elevate{}
			err = json.Unmarshal(pkt.Raw, data)
			if err != nil {
				c.Send(packets.NewError(err.Error()))
				return
			}

			err = s.handleElevate(data, c)
			if err != nil {
				c.Send(packets.NewError(err.Error()))
				return
			}
		case packets.PacketTypeClientMessagesRequest:
			data := &packets.ClientMessagesRequest{}
			err = json.Unmarshal(pkt.Raw, data)
			if err != nil {
				c.Send(packets.NewError(err.Error()))
				return
			}

			err = s.handleClientMessagesRequest(data, c)
			if err != nil {
				c.Send(packets.NewError(err.Error()))
				return
			}
		}
	}
}
