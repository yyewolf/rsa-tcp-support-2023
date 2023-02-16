package packets

type Message struct {
	Message string `json:"m"`
	Target  string `json:"t"`
	From    string `json:"f"`
}

func NewMessage(message string, target string, From string) *Packet {
	return NewPacket(PacketTypeMessage, &Message{
		Message: message,
		Target:  target,
		From:    From,
	})
}
