package ui

import (
	"github.com/rivo/tview"
)

type InformationField struct {
	View *tview.TextView
}

func NewInformationField() *InformationField {
	view := tview.NewTextView().
		SetText("P2P").
		SetTextAlign(tview.AlignCenter)

	view.SetTitle("P2P-chat").SetBorder(true)

	return &InformationField{
		View: view,
	}
}
