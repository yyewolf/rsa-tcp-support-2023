package ui

import "github.com/rivo/tview"

func loadLoadingPage() {
	// Create loading page
	loadingPage := tview.NewFlex().SetDirection(tview.FlexRow)
	loadingPage.SetBorder(true).SetTitle("Chargement...").SetTitleAlign(tview.AlignLeft)

	// Create loading text
	loadingText := tview.NewTextView().SetTextAlign(tview.AlignCenter).SetText("Chargement...")

	// Add loading text to loading page
	loadingPage.AddItem(loadingText, 0, 1, false)

	// Add loading page to pages
	Pages.AddPage("loading", loadingPage, true, false)
}
