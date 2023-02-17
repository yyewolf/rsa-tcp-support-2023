package hub

import (
	"errors"
	"hub/internal/packets"
	"log"
	"os"
	"strconv"
	"sync"
	"time"
)

func (s *Server) handleIdentify(data *packets.Identify, c *Connection) (err error) {
	if c.Client != nil || c.Agent != nil {
		err = errors.New("already identified")
		return
	}

	// If auth is a secret password then it's an agent.
	if data.Auth == os.Getenv("AUTH_SECRET") {
		if data.Name == "" {
			err = errors.New("invalid name")
			return
		}

		if data.Level < 0 || data.Level > 2 {
			err = errors.New("invalid level")
			return
		}

		// Create a new agent.
		agent := &Agent{
			ID:    c.ID,
			Name:  data.Name,
			Level: data.Level,
			Conn:  c.Conn,
			Mutex: &sync.Mutex{},
		}

		c.Agent = agent

		// Add the agent to the agents map.
		s.Agents[data.Level] = append(s.Agents[c.ID], agent)

		c.Send(packets.NewSuccess(true))
		if os.Getenv("DEBUG") == "true" {
			log.Println("[HUB]", c.ID, "is now identified as an agent.")
		}

		// Send the agent count to all the clients connected in this level.
		for _, client := range s.Clients {
			if client.Level != data.Level {
				continue
			}
			client.Send(packets.NewAgentCount(len(s.Agents)))
		}

		// Send the clients to the agent.
		var clients = make([]int, 0)
		for _, client := range s.Clients {
			if client.Level != data.Level {
				continue
			}
			clients = append(clients, client.ID)
		}
		agent.Send(packets.NewClientPresent(clients))

		return
	}

	if data.Auth != "" {
		err = errors.New("invalid auth")
		return
	}

	// Create a new client.
	client := &Client{
		ID:    c.ID,
		Level: 0,
		Conn:  c.Conn,
		Mutex: &sync.Mutex{},
	}

	c.Client = client

	// Add the client to the clients map.
	s.Clients[c.ID] = client

	c.Send(packets.NewSuccess(true))
	if os.Getenv("DEBUG") == "true" {
		log.Println("[HUB]", c.ID, "is now identified.")
	}

	// Send agent count every 5 seconds.
	go func() {
		for {
			err := c.Send(packets.NewAgentCount(len(s.Agents[c.Client.Level])))
			if err != nil {
				return
			}
			time.Sleep(5 * time.Second)
		}
	}()

	return
}

func (s *Server) handleMessage(data *packets.Message, c *Connection) (err error) {
	if c.Client == nil && c.Agent == nil {
		err = errors.New("not identified")
		return
	}

	if data.Message == "" {
		err = errors.New("invalid message")
		return
	}

	if c.Client != nil {
		// Send the message to all the agents connected in this level.
		for _, agent := range s.Agents[c.Client.Level] {
			agent.Send(packets.NewMessage(data.Message, "", strconv.Itoa(c.Client.ID)))
		}
		c.Client.History = append(c.Client.History, data)
		return
	}

	t, err := strconv.Atoi(data.Target)
	if err != nil {
		err = errors.New("invalid target")
		return
	}

	// Send the message to the target client.
	if target, ok := s.Clients[t]; ok {
		m := packets.NewMessage(data.Message, "", c.Agent.Name)
		target.Send(m)
		target.History = append(target.History, m.Data.(*packets.Message))
	} else {
		err = errors.New("invalid target")
	}

	return
}

func (s *Server) handleElevate(data *packets.Elevate, c *Connection) (err error) {
	if c.Client == nil {
		err = errors.New("not identified")
		return
	}

	if c.Client.Level == 2 {
		err = errors.New("already elevated")
		return
	}

	// Send an elevated message to all the agents connected in this level.
	for _, agent := range s.Agents[c.Client.Level] {
		agent.Send(packets.NewElevate(c.Client.ID))
	}

	c.Client.Level++
	// Send a message to the client.
	c.Client.Send(packets.NewMessage("You have been elevated to level "+strconv.Itoa(c.Client.Level), "", "SYSTEM"))
	return
}

func (s *Server) handleClientMessagesRequest(data *packets.ClientMessagesRequest, c *Connection) (err error) {
	if c.Agent == nil {
		err = errors.New("not identified")
		return
	}

	client, ok := s.Clients[data.ClientID]
	if !ok {
		err = errors.New("invalid client id")
		return
	}

	if client.Level != c.Agent.Level {
		err = errors.New("invalid client id")
		return
	}

	// Send the client's history to the agent.
	c.Agent.Send(packets.NewClientMessagesResponse(client.History))
	return
}