package ui

import (
	"github.com/rivo/tview"
)

var App *tview.Application

func Start() {
	App = tview.NewApplication()
	p := loadPages()

	loadMainPage()
	loadLoadingPage()
	loadChatPage()

	if err := App.SetRoot(p, true).EnableMouse(true).Run(); err != nil {
		panic(err)
	}
}
