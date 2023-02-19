package bot

import (
	"bot/internal/packets"
	"encoding/json"
	"log"
	"net"
	"os"
	"strconv"
)

var conn net.Conn

func sendPacket(p *packets.Packet) {
	if _, err := conn.Write(p.ToBytes()); err != nil {
		panic(err)
	}
}

func Start() error {
	// Get modulus
	mod, err := strconv.Atoi(os.Getenv("MOD"))
	if err != nil {
		return err
	}

	modR, err := strconv.Atoi(os.Getenv("MOD_R"))
	if err != nil {
		return err
	}

	// Connect to the server
	conn, err = net.Dial("tcp", os.Getenv("HOST"))
	if err != nil {
		panic(err)
	}

	// Send the name
	p := packets.NewIdentify(os.Getenv("AUTH"), os.Getenv("NAME"))
	sendPacket(p)

	// Listen for packets
	for {
		p, err := packets.ReadPacket(conn)
		if err != nil {
			panic(err)
		}

		switch p.Type {
		case packets.PacketTypeSuccess:
			data := &packets.Success{}
			err = json.Unmarshal(p.Raw, data)
			if err != nil {
				return err
			}

			if data.Success {
				log.Println("Successfully connected to the server")
			} else {
				log.Printf("Failed to connect to '%s' with password '%s'\n", os.Getenv("HOST"), os.Getenv("AUTH"))
				return nil
			}
		case packets.PacketTypeMessage:
			data := &packets.Message{}
			err = json.Unmarshal(p.Raw, data)
			if err != nil {
				return err
			}

			id, _ := strconv.Atoi(data.From)
			if id%mod == modR {
				log.Printf("Received message from %s: %s\n", data.From, data.Message)
				replyToMessage(data)
			}
		}
	}
}
