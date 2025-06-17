package main

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"net"
)

var (
	Conn, _ = net.Dial("tcp", "localhost:4545")
)

func ui() {
	myApp := app.New()

	var UserName string
	f := myApp.NewWindow("Name")
	f.Resize(fyne.NewSize(300, 100))
	entry := widget.NewEntry()

	form := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "Enter something", Widget: entry},
		},
		OnSubmit: func() {
			UserName = entry.Text
			f.Close()
		},
	}
	c := container.NewVBox(form)
	f.SetContent(c)
	f.Show()

	myWindow := myApp.NewWindow("Chat")
	myWindow.Resize(fyne.NewSize(400, 400))

	historyLabel := widget.NewLabel("")
	historyLabel.Wrapping = fyne.TextWrapWord

	scrollContainer := container.NewScroll(historyLabel)
	scrollContainer.SetMinSize(fyne.NewSize(300, 360))
	scrollContainer.ScrollToBottom()

	inputEntry := widget.NewEntry()
	inputEntry.SetPlaceHolder("Введите сообщение")

	update := widget.NewButton("Появились сообщения", func() {
		mas, err := Writing()
		if err != nil {
			historyLabel.SetText(fmt.Sprintf("%v", err))
			return
		}
		if len(mas) == 0 {
			return
		} else {
			historyLabel.SetText("")
		}
		for i := 0; i < len(mas); i++ {
			historyLabel.SetText(fmt.Sprintf("%v%v: %v\n", historyLabel.Text, mas[i].Name, mas[i].Text))
		}
	})

	inputEntry.OnSubmitted = func(text string) {
		if UserName == "" {
			UserName = "Anonymous"
		}
		if text == "" {
			return
		}
		historyLabel.SetText(fmt.Sprintf("%v%v: %v\n", historyLabel.Text, UserName, text))
		Writer(UserName, text)
		inputEntry.SetText("")
		scrollContainer.ScrollToBottom()
	}

	inputContainer := container.NewBorder(nil, nil, nil, update, inputEntry)
	content := container.NewBorder(scrollContainer, inputContainer, nil, nil)
	myWindow.SetContent(content)
	myWindow.ShowAndRun()
}
