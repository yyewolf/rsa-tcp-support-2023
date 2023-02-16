package packets

// Hello is the packet sent to the client when it connects.
type Hello struct {
	ID int `json:"i"`
}

func NewHello(id int) *Packet {
	return NewPacket(PacketTypeHello, &Hello{
		ID: id,
	})
}
