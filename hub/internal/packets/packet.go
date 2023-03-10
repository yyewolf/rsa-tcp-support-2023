package packets

import "encoding/json"

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
