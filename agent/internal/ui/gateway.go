package ui

import (
	"agent/internal/packets"
	"agent/internal/settings"
	"encoding/json"
	"fmt"
	"net"
	"strconv"
)

func handleMessages(c net.Conn) {
	go func() {
		for {
			// Read packet
			p, err := packets.ReadPacket(c)
			if err != nil {
				goto err
			}

			// Handle packet
			switch p.Type {
			case packets.PacketTypeSuccess:
				// We are connected and can switch to the chat page
				Pages.SwitchToPage("chat")
				App.ForceDraw()
			case packets.PacketTypeClientPresent:
				data := &packets.ClientPresent{}
				err = json.Unmarshal(p.Raw, data)
				if err != nil {
					goto err
				}

				// Add client to list
				for _, client := range data.IDs {
					UserList.AddItem(strconv.Itoa(client), "", 0, nil)
				}
				App.ForceDraw()
				// Request messages from the first client
				if len(data.IDs) > 0 {
					settings.Settings.Send(packets.NewClientMessagesRequest(data.IDs[0]))
					settings.Settings.SelectedClient = data.IDs[0]
				}

			case packets.PacketTypeClientMessagesResponse:
				data := &packets.ClientMessagesResponse{}
				err = json.Unmarshal(p.Raw, data)
				if err != nil {
					goto err
				}

				// Empty the chat
				MessageList.Clear()

				// Add messages to chat
				for _, message := range data.Messages {
					MessageList.SetText(fmt.Sprintf("%s%s: %s", MessageList.GetText(false), message.From, message.Message))
				}
				MessageInput.SetText("")
				App.ForceDraw()
			case packets.PacketTypeError:
				data := &packets.Error{}
				err = json.Unmarshal(p.Raw, data)
				if err != nil {
					goto err
				}

				// Error
				LoadErrorPage("main", data.Error)
			case packets.PacketTypeMessage:
				data := &packets.Message{}
				err = json.Unmarshal(p.Raw, data)
				if err != nil {
					goto err
				}

				if data.From == strconv.Itoa(settings.Settings.SelectedClient) {
					// Message from the selected client
					MessageList.SetText(fmt.Sprintf("%s%s: %s", MessageList.GetText(false), data.From, data.Message))
					App.ForceDraw()
				} else {
					// If client doesn't exist, add it
					items := UserList.FindItems(data.From, "", false, true)
					if len(items) == 0 {
						UserList.AddItem(data.From, "", 0, nil)
						App.ForceDraw()
					}

					// If there's no selected client, select the new one
					if settings.Settings.SelectedClient == 0 {
						settings.Settings.SelectedClient, _ = strconv.Atoi(data.From)
						settings.Settings.Send(packets.NewClientMessagesRequest(settings.Settings.SelectedClient))
					}
				}
			case packets.PacketTypeElevate:
				data := &packets.Elevate{}
				err = json.Unmarshal(p.Raw, data)
				if err != nil {
					goto err
				}

				// Remove the client from the list
				items := UserList.FindItems(strconv.Itoa(data.ID), "", false, true)
				if len(items) > 0 {
					UserList.RemoveItem(items[0])
					App.ForceDraw()
				}

				// If the client was selected, select the first client
				if settings.Settings.SelectedClient == data.ID {
					if UserList.GetItemCount() > 0 {
						main, _ := UserList.GetItemText(0)
						settings.Settings.SelectedClient, _ = strconv.Atoi(main)
						settings.Settings.Send(packets.NewClientMessagesRequest(settings.Settings.SelectedClient))
					} else {
						settings.Settings.SelectedClient = 0
						// Empty the chat
						MessageList.Clear()
						App.ForceDraw()
					}
				}

			}
		}
	err:
		// Error
		LoadErrorPage("main", "Erreur lors de la lecture d'un paquet")
	}()
}
