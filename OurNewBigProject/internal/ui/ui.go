package ui

import (
	"OurNewBigProject/internal"
	"OurNewBigProject/internal/MessageStruct"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

const (
	reprintFrequency = 50 * time.Millisecond
)

type App struct {
	Proto       *internal.Proto
	Chat        *Chat
	Sidebar     *Sidebar
	InfoField   *InformationField
	View        *tview.Flex
	UI          *tview.Application
	CurrentPeer *MessageStruct.Peer
}

func NewApp(proto *internal.Proto) *App {
	app := &App{
		Proto:       proto,
		Chat:        NewChat(),
		Sidebar:     NewSidebar(proto.Peers),
		InfoField:   NewInformationField(),
		View:        tview.NewFlex(),
		UI:          tview.NewApplication(),
		CurrentPeer: nil,
	}

	app.initView()
	app.initUI()
	app.initBindings()

	app.run()

	return app
}

func (app *App) Run() error {
	return app.UI.SetRoot(app.View, true).SetFocus(app.Sidebar.View).Run()
}

func (app *App) initView() {
	app.View.
		AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
			AddItem(app.InfoField.View, 3, 2, false).
			AddItem(app.Sidebar.View, 0, 1, false), 0, 1, false).
		AddItem(app.Chat.View, 0, 3, false)
}

func (app *App) initUI() {
	app.UI.SetRoot(app.View, true).SetFocus(app.Sidebar.View)
}

func (app *App) initBindings() {
	app.Sidebar.View.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Rune() == 'l' {
			app.UI.SetFocus(app.Chat.Messages)
		}

		if event.Key() == tcell.KeyEnter {
			if app.Sidebar.View.GetItemCount() > 0 {
				app.CurrentPeer = app.getCurrentPeer()
				app.UI.SetFocus(app.Chat.Messages)
			}
		}

		return event
	})

	app.Chat.Messages.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Rune() {
		case 'h':
			app.UI.SetFocus(app.Sidebar.View)
		case 'j':
			app.UI.SetFocus(app.Chat.InputField)
		}

		return event
	})

	app.Chat.InputField.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyUp {
			app.UI.SetFocus(app.Chat.Messages)
		}

		if event.Key() == tcell.KeyEnter {
			if app.CurrentPeer == nil || app.Chat.InputField.GetText() == "" {
				return event
			}

			message := app.Chat.InputField.GetText()
			app.CurrentPeer.AddMessage(message, app.Proto.Name)
			app.Chat.InputField.SetText("")
		}

		return event
	})
}

func (app *App) renderMessages() {
	if app.CurrentPeer != nil {
		app.Chat.RenderMessages(app.CurrentPeer.Messages, app.Proto.Name)
		app.Chat.View.SetTitle(app.CurrentPeer.Name)
	}
}

func (app *App) getCurrentPeer() *MessageStruct.Peer {
	_, pubKey := app.Sidebar.View.GetItemText(
		app.Sidebar.View.GetCurrentItem())

	peer, found := app.Proto.Peers.Get(pubKey)
	if !found {
		return nil
	}

	return peer
}

func (app *App) run() {
	ticker := time.NewTicker(reprintFrequency)

	go func() {
		for {
			<-ticker.C
			app.UI.QueueUpdateDraw(app.Sidebar.Reprint)
			app.UI.QueueUpdateDraw(app.renderMessages)
		}
	}()
}
