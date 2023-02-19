package packets

type Identify struct {
	Auth string `json:"a"`
	Name string `json:"n"`
}

func NewIdentify(auth, name string) *Packet {
	return NewPacket(PacketTypeIdentify, &Identify{
		Auth: auth,
		Name: name,
	})
}
