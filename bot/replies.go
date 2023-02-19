package bot

import (
	"bot/internal/packets"
	"strings"
)

func replyToMessage(p *packets.Message) {
	p.Message = strings.ToLower(p.Message)
	if p.Message == "hello" {
		r := packets.NewMessage("Hello!", p.From, "")
		sendPacket(r)
		return
	}

	if strings.Contains(p.Message, "meaning") && strings.Contains(p.Message, "life") {
		r := packets.NewMessage("42", p.From, "")
		sendPacket(r)
		return
	}

	if strings.Contains(p.Message, "what") && strings.Contains(p.Message, "time") {
		r := packets.NewMessage("It's time to get a watch!", p.From, "")
		sendPacket(r)
		return
	}

	if strings.Contains(p.Message, "what") && strings.Contains(p.Message, "weather") {
		r := packets.NewMessage("It's sunny outside!", p.From, "")
		sendPacket(r)
		return
	}

	if strings.Contains(p.Message, "what") && strings.Contains(p.Message, "your") && strings.Contains(p.Message, "name") {
		r := packets.NewMessage("My name is Bot!", p.From, "")
		sendPacket(r)
		return
	}

	// Actually help the client if their message contains "issue"
	if strings.Contains(p.Message, "issue") {
		if strings.Contains(p.Message, "help") {
			r := packets.NewMessage("I can help you with your issue!", p.From, "")
			sendPacket(r)
			return
		}

		if strings.Contains(p.Message, "fridge") {
			r := packets.NewMessage("If you have an issue with your fridge, have you considered checking the manual ?", p.From, "")
			sendPacket(r)
			return
		}

		if strings.Contains(p.Message, "car") {
			r := packets.NewMessage("If you have an issue with your car, have you already gone to a garage ?", p.From, "")
			sendPacket(r)
			return
		}
	}

	// Default message
	msg := "Sorry, but I don't understand what you mean. Please elevate your message to a higher level as I cannot help you."
	r := packets.NewMessage(msg, p.From, "")
	sendPacket(r)
}
