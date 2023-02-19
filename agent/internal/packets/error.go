package packets

type Error struct {
	Error string `json:"e"`
}

func NewError(err string) *Packet {
	return NewPacket(PacketTypeError, &Error{
		Error: err,
	})
}
