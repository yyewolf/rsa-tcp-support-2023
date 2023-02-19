package ui

import (
	"agent/internal/packets"
	"agent/internal/settings"
	"net"
	"time"

	"github.com/rivo/tview"
)

var MainPage *tview.Form

func loadMainPage() {
	MainPage = tview.NewForm()
	MainPage.SetBorder(true)
	MainPage.SetTitle("Connexion agent")
	MainPage.SetTitleAlign(tview.AlignLeft)
	MainPage.SetBorderPadding(0, 0, 1, 0)
	MainPage.AddInputField("Authentification", "", 20, nil, func(text string) {
		settings.Settings.Auth = text
	})
	MainPage.AddInputField("Nom", "", 20, nil, func(text string) {
		settings.Settings.Name = text
	})
	MainPage.AddButton("Connect", func() {
		Pages.SwitchToPage("loading")
		App.ForceDraw()

		// Try to connect to TCP server
		c, err := net.Dial("tcp", "127.0.0.1:8000")
		if err != nil {
			LoadErrorPage("main", "Impossible de se connecter au serveur")
			App.ForceDraw()
			return
		}

		handleMessages(c)
		settings.Settings.Conn = c

		time.Sleep(1 * time.Second)

		// Send authentification
		settings.Settings.Send(packets.NewIdentify(settings.Settings.Auth, settings.Settings.Name))
	})
	MainPage.AddButton("Quit", func() {
		App.Stop()
	})

	Pages.AddPage("main", MainPage, true, true)

}
