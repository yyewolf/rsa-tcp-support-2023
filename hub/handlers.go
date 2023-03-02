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
	if data.Auth == os.Getenv("AUTH_SECRET_0") || data.Auth == os.Getenv("AUTH_SECRET_1") || data.Auth == os.Getenv("AUTH_SECRET_2") {
		if data.Name == "" {
			err = errors.New("invalid name")
			return
		}

		// Create a new agent.
		agent := &Agent{
			ID:    c.ID,
			Name:  data.Name,
			Conn:  c.Conn,
			Mutex: &sync.Mutex{},
		}

		if data.Auth == os.Getenv("AUTH_SECRET_0") {
			agent.Level = 0
		} else if data.Auth == os.Getenv("AUTH_SECRET_1") {
			agent.Level = 1
		} else if data.Auth == os.Getenv("AUTH_SECRET_2") {
			agent.Level = 2
		}

		c.Agent = agent

		// Add the agent to the agents map.
		s.Agents[agent.Level] = append(s.Agents[agent.Level], agent)

		c.Send(packets.NewSuccess(true))
		if os.Getenv("DEBUG") == "true" {
			log.Println("[HUB]", c.ID, "is now identified as an agent.")
		}

		// Send the agent count to all the clients connected in this level.
		for _, client := range s.Clients {
			if client.Level != agent.Level {
				continue
			}
			client.Send(packets.NewAgentCount(len(s.Agents[agent.Level])))
		}

		// Send the clients to the agent.
		var clients = make([]int, 0)
		for _, client := range s.Clients {
			if client.Level != agent.Level {
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
		p := packets.NewMessage(data.Message, "", strconv.Itoa(c.Client.ID))
		c.Client.History = append(c.Client.History, p.Data.(*packets.Message))
		for _, agent := range s.Agents[c.Client.Level] {
			agent.Send(p)
		}
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
