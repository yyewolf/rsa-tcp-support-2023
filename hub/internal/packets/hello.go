package packets

// Hello is the packet sent to the client when it connects.
type Hello struct {
	ID string
}

func NewHello(id string) *Packet {
	return NewPacket(PacketTypeHello, &Hello{
		ID: id,
	})
}
