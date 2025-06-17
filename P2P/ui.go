package main

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
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
	historyLabel.Wrapping = fyne.TextWrapWord // Перенос слов

	// Создаем контейнер с фиксированным размером и ползунком
	scrollContainer := container.NewScroll(historyLabel)
	scrollContainer.SetMinSize(fyne.NewSize(300, 360))
	scrollContainer.ScrollToBottom()

	inputEntry := widget.NewEntry()
	inputEntry.SetPlaceHolder("Введите сообщение")

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

	go func() {
		name, text := Listener()
		if name == "" || text == "" || name == UserName {
			return
		}
		historyLabel.SetText(fmt.Sprintf("%v%v: %v\n", historyLabel.Text, name, text))
	}()

	// Размещаем элементы в контейнерах
	inputContainer := container.NewBorder(nil, nil, nil, nil, inputEntry)
	content := container.NewBorder(scrollContainer, inputContainer, nil, nil)
	// Устанавливаем содержимое окна и показываем его
	myWindow.SetContent(content)
	myWindow.ShowAndRun()
	Listener()
}
