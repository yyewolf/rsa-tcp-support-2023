package ui

import "github.com/rivo/tview"

var Pages *tview.Pages

func loadPages() *tview.Pages {
	Pages = tview.NewPages()
	Pages.SetTitle("Pages")
	Pages.SetTitleAlign(tview.AlignLeft)

	return Pages
}
