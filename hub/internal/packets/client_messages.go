package packets

type ClientMessagesRequest struct {
	ClientID int `json:"i"`
}

func NewClientMessagesRequest(clientID int) *Packet {
	return NewPacket(PacketTypeClientMessagesRequest, &ClientMessagesRequest{
		ClientID: clientID,
	})
}

type ClientMessagesResponse struct {
	Messages []*Message `json:"m"`
}

func NewClientMessagesResponse(messages []*Message) *Packet {
	return NewPacket(PacketTypeClientMessagesResponse, &ClientMessagesResponse{
		Messages: messages,
	})
}
