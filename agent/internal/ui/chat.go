package ui

import (
	"agent/internal/packets"
	"agent/internal/settings"
	"fmt"
	"strconv"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

var ChatPage *tview.Flex
var ChatBox *tview.Flex
var UserList *tview.List
var MessageList *tview.TextView
var MessageInput *tview.InputField

func loadChatPage() {
	// Chat has a list of users on the left and a chatbox on the right
	// The chatbox has a list of messages and an input field at the bottom

	ChatPage = tview.NewFlex()
	// ChatPage.SetBorder(true)
	// ChatPage.SetTitle("Chat")
	ChatPage.SetTitleAlign(tview.AlignLeft)
	// ChatPage.SetBorderPadding(0, 0, 1, 0)

	// Add the userlist
	UserList = tview.NewList()
	UserList.SetBorder(true)
	UserList.SetTitle("Utilisateurs")
	UserList.SetTitleAlign(tview.AlignLeft)
	UserList.SetBorderPadding(0, 0, 1, 0)

	UserList.SetSelectedFunc(func(index int, mainText string, secondaryText string, shortcut rune) {
		// Get client ID
		client, _ := strconv.Atoi(mainText)
		// Select client
		settings.Settings.SelectedClient = client
		// Change title
		ChatBox.SetTitle(fmt.Sprintf("Chat avec %s", mainText))
		// Request messages
		settings.Settings.Send(packets.NewClientMessagesRequest(client))
	})

	// Add the chatbox
	ChatBox = tview.NewFlex()
	// ChatBox.SetBorder(true)
	ChatBox.SetTitle("Chat")
	ChatBox.SetTitleAlign(tview.AlignLeft)
	// ChatBox.SetBorderPadding(0, 0, 1, 0)
	ChatBox.SetDirection(tview.FlexRow)

	// Add the message list
	MessageList = tview.NewTextView()
	MessageList.SetBorder(true)
	MessageList.SetTitle("Messages")
	MessageList.SetTitleAlign(tview.AlignLeft)
	MessageList.SetBorderPadding(0, 0, 1, 0)
	// Enable text wrapping
	MessageList.SetWrap(true)
	MessageList.SetWordWrap(true)

	// Add the message input
	MessageInput = tview.NewInputField()
	// MessageInput.SetBorder(true)
	MessageInput.SetTitle("Message")
	MessageInput.SetTitleAlign(tview.AlignLeft)
	MessageInput.SetLabel("Message: ")
	// MessageInput.SetBorderPadding(0, 0, 1, 0)
	MessageInput.SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEnter {
			// Check that the message is not empty and that there's a selected client
			if MessageInput.GetText() == "" || settings.Settings.SelectedClient == 0 {
				return
			}
			// Add the message to the message list
			MessageList.SetText(fmt.Sprintf("%s%s: %s", MessageList.GetText(false), settings.Settings.Name, MessageInput.GetText()))
			// Send the message to the server
			settings.Settings.Send(packets.NewMessage(MessageInput.GetText(), strconv.Itoa(settings.Settings.SelectedClient), ""))
			// Clear the input field
			MessageInput.SetText("")
		}
	})
	// Add the message list and input to the chatbox
	ChatBox.AddItem(MessageList, 0, 7, true)
	ChatBox.AddItem(MessageInput, 0, 1, false)

	// Add the userlist and chatbox to the chat page
	ChatPage.AddItem(UserList, 0, 1, false)
	ChatPage.AddItem(ChatBox, 0, 4, true)

	// Add the chat page to the pages
	Pages.AddPage("chat", ChatPage, true, false)
}
