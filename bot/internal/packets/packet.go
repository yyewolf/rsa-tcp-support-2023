package packets

import (
	"encoding/json"
	"net"
)

// PacketType is the type of a packet.
type PacketType int

const (
	PacketTypeHello PacketType = iota
	PacketTypeIdentify
	PacketTypeSuccess
	PacketTypeError
	PacketTypeAgentCount
	PacketTypeMessage
	PacketTypeElevate
	PacketTypeClientPresent
	PacketTypeClientMessagesRequest
	PacketTypeClientMessagesResponse
)

// Packet is a packet sent to the client.
type Packet struct {
	Type PacketType  `json:"t"`
	Data interface{} `json:"d"`
	Raw  []byte      `json:"-"`
}

// NewPacket creates a new packet.
func NewPacket(t PacketType, d interface{}) *Packet {
	return &Packet{
		Type: t,
		Data: d,
	}
}

// ToBytes converts the packet to bytes.
func (p *Packet) ToBytes() []byte {
	d, _ := json.Marshal(p)
	return d
}

// ReadPacket reads a packet from a connection.
func ReadPacket(c net.Conn) (*Packet, error) {
	// Read the packet.
	buf := make([]byte, 1024)
	n, err := c.Read(buf)
	if err != nil {
		return nil, err
	}

	// Unmarshal the packet.
	pkt := &Packet{}
	err = json.Unmarshal(buf[:n], pkt)
	if err != nil {
		return nil, err
	}

	// Remarshal the packet data.
	pkt.Raw, err = json.Marshal(pkt.Data)
	if err != nil {
		return nil, err
	}

	return pkt, nil
}
