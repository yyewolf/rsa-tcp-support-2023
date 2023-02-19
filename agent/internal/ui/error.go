package ui

import "github.com/rivo/tview"

var ErrorPage *tview.Form
var ErrorPageReturn string

func LoadErrorPage(r, msg string) {
	ErrorPage = tview.NewForm().
		AddTextView("", msg, 0, 0, true, true).
		AddButton("Retour", func() {
			Pages.SwitchToPage(r)
		})
	ErrorPage.SetBorder(true).SetTitle("Erreur").SetTitleAlign(tview.AlignLeft)
	Pages.AddPage("error", ErrorPage, true, true)
	App.ForceDraw()
}
