package packets

type AgentCount struct {
	Count int `json:"c"`
}

func NewAgentCount(count int) *Packet {
	return NewPacket(PacketTypeAgentCount, &AgentCount{
		Count: count,
	})
}
