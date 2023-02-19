package settings

import (
	"agent/internal/packets"
	"net"
	"sync"
)

type Setting struct {
	// Agent settings
	Name string `json:"name"`
	Auth string `json:"auth"`

	SelectedClient int `json:"selected_client"`

	Conn net.Conn `json:"-"`
	*sync.Mutex
}

var Settings = &Setting{
	Mutex: &sync.Mutex{},
}

func (s *Setting) Send(p *packets.Packet) {
	s.Lock()
	defer s.Unlock()

	s.Conn.Write(p.ToBytes())
}
