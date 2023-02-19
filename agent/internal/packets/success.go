package packets

type Success struct {
	Success bool `json:"s"`
}

func NewSuccess(success bool) *Packet {
	return NewPacket(PacketTypeSuccess, &Success{
		Success: success,
	})
}
