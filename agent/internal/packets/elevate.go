package packets

type Elevate struct {
	ID int `json:"i"`
}

func NewElevate(id int) *Packet {
	return NewPacket(PacketTypeElevate, &Elevate{
		ID: id,
	})
}
