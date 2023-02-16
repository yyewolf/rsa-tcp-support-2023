package packets

type Identify struct {
	Level int    `json:"l"`
	Auth  string `json:"a"`
	Name  string `json:"n"`
}

func NewIdentify(level int, auth, name string) *Packet {
	return NewPacket(PacketTypeIdentify, &Identify{
		Level: level,
		Auth:  auth,
		Name:  name,
	})
}
