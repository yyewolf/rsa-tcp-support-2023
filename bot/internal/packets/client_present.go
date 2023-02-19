package packets

type ClientPresent struct {
	IDs []int `json:"ids"`
}

func NewClientPresent(ids []int) *Packet {
	return NewPacket(PacketTypeClientPresent, &ClientPresent{
		IDs: ids,
	})
}
