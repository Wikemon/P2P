package ui

import (
	"OurNewBigProject/internal/MessageStruct"
	"fmt"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

const (
	timeFormat        = "15:04:05"
	maxMessagesInView = 100
)

type Chat struct {
	View       *tview.Flex
	InputField *tview.InputField
	Messages   *tview.TextView
}

func NewChat() *Chat {
	view := tview.NewFlex().SetDirection(tview.FlexRow)
	view.SetTitle("chat").SetBorder(true)

	messages := tview.NewTextView().SetText("").
		SetDynamicColors(true)

	inputField := tview.NewInputField().SetFieldBackgroundColor(tview.Styles.PrimitiveBackgroundColor).
		SetDoneFunc(func(key tcell.Key) {})
	inputField.SetBorder(true)

	view.
		AddItem(messages, 0, 14, false).
		AddItem(inputField, 0, 1, false)

	return &Chat{
		View:       view,
		InputField: inputField,
		Messages:   messages,
	}
}

func (c *Chat) RenderMessages(messages []*MessageStruct.Message, protoName string) {
	text := strings.Repeat("\n", maxMessagesInView)
	for _, message := range messages {
		isAuthor := false
		if message.Author == protoName {
			isAuthor = true
		}

		text += fmt.Sprintf("%s %s: %s\n",
			formatTime(message),
			formatAuthor(message, isAuthor),
			formatText(message))
	}

	c.Messages.SetText(text[:len(text)-1]).ScrollToEnd()
}

func formatTime(message *MessageStruct.Message) string {
	now := message.Time.UTC()
	return fmt.Sprintf("%s%s", "[blue]", now.Format(timeFormat))
}

func formatAuthor(message *MessageStruct.Message, isAuthor bool) string {
	if isAuthor {
		return fmt.Sprintf("%s%s", "[green]", message.Author)
	}
	return fmt.Sprintf("%s%s", "[red]", message.Author)
}

func formatText(message *MessageStruct.Message) string {
	return fmt.Sprintf("%s%s", "[white]", message.Text)
}
