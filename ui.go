package main

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"image/color"
	"os"
)

type ChatTheme struct {
	PanelColor  color.Color // Цвет боковой панели
	ChatBgColor color.Color // Цвет фона чата
	AvatarImage string      // Путь к изображению аватара
}

var AllUsersNames map[string]struct{}

func ui() {
	myApp := app.New()

	nameWindow := myApp.NewWindow("welcome slay")
	nameWindow.Resize(fyne.NewSize(400, 400))

	bgImage := canvas.NewImageFromFile("images/fish.jpg")
	bgImage.FillMode = canvas.ImageFillStretch

	var UserName string
	entry := widget.NewEntry()
	entry.SetPlaceHolder("Enter your SLAY NAME")

	submitBtn := createImageButton("images/horse.jpg", func() {
		UserName = entry.Text
		if UserName == "" {
			UserName = "unnamed guy"
		}
		if UserName == "Oleg" {
			UserName = "silly guy"
		}
		nameWindow.Close()
		createChatWindow(myApp, UserName)
	})

	form := container.NewVBox(
		entry,
		container.NewCenter(submitBtn),
		container.NewCenter(
			canvas.NewText("(only SLAY NAME)!!!!!!!!or u will be kicked:", color.White),
		),
	)

	nameWindow.SetContent(container.NewStack(
		bgImage,
		container.NewCenter(form),
	))

	nameWindow.Show()
	myApp.Run()
}

func createImageButton(imagePath string, tapped func()) fyne.CanvasObject {
	img := canvas.NewImageFromFile(imagePath)
	img.SetMinSize(fyne.NewSize(150, 50))

	btn := widget.NewButton("", tapped)
	btn.Importance = widget.LowImportance
	btn.Resize(fyne.NewSize(150, 50))

	return container.NewStack(
		img,
		btn,
	)
}
func roundedSquareButton(text string, bgColor color.Color, tapped func()) *fyne.Container {

	bg := canvas.NewRectangle(bgColor)
	bg.CornerRadius = 8
	bg.SetMinSize(fyne.NewSize(60, 60))

	label := canvas.NewText(text, color.White)
	label.Alignment = fyne.TextAlignCenter
	label.TextSize = 12

	btnContent := container.NewStack(
		bg,
		container.NewCenter(label),
	)
	btn := widget.NewButton("", tapped)
	btn.Importance = widget.LowImportance

	return container.NewStack(
		btn,
		btnContent,
	)
}

func circleWidget(size float32, avatar *canvas.Image) fyne.CanvasObject {
	avatar.FillMode = canvas.ImageFillContain
	avatar.SetMinSize(fyne.NewSize(size, size))

	circle := canvas.NewCircle(color.Transparent)
	circle.Resize(fyne.NewSize(size, size))

	return container.NewStack(
		circle,
		avatar,
	)
}
func Theme(bgColor color.Color, bgc color.Color, image string) ChatTheme {
	return ChatTheme{
		PanelColor:  bgColor,
		ChatBgColor: bgc,
		AvatarImage: image,
	}
}

var (
	BlueTheme = Theme(
		color.NRGBA{R: 173, G: 216, B: 230, A: 255}, // Панель
		color.NRGBA{R: 240, G: 248, B: 255, A: 255}, // Фон чата
		"images/blue.png",                           // Аватар
	)
	PinkTheme = Theme(
		color.NRGBA{R: 255, G: 182, B: 193, A: 230},
		color.NRGBA{R: 255, G: 240, B: 245, A: 128},
		"images/pink.png",
	)
	GreenTheme = Theme(
		color.NRGBA{R: 152, G: 251, B: 152, A: 230},
		color.NRGBA{R: 152, G: 251, B: 152, A: 128},
		"images/green.png",
	)
	YellowTheme = Theme(
		color.NRGBA{R: 255, G: 255, B: 153, A: 230},
		color.NRGBA{R: 255, G: 255, B: 153, A: 100},
		"images/yellow.png",
	)
)

func applyTheme(theme ChatTheme, panel *canvas.Rectangle, chatBg *canvas.Rectangle, avatar *canvas.Image) {
	panel.FillColor = theme.PanelColor
	chatBg.FillColor = theme.ChatBgColor
	avatar.File = theme.AvatarImage
	panel.Refresh()
	chatBg.Refresh()
	avatar.Refresh()
}
func createChatWindow(a fyne.App, name string) {
	AllUsersNames = make(map[string]struct{}, 4096)
	chatWindow := a.NewWindow("boring chat " + name)
	chatWindow.Resize(fyne.NewSize(400, 400))

	currentTheme := GreenTheme

	panelBg := canvas.NewRectangle(currentTheme.PanelColor)
	chatBg := canvas.NewRectangle(currentTheme.ChatBgColor)
	avatar := canvas.NewImageFromFile(currentTheme.AvatarImage)

	panel := container.NewStack(
		panelBg,
		container.NewVBox(
			container.NewPadded(
				circleWidget(60, avatar),
			),
			container.NewCenter(
				func() *canvas.Text {
					t := canvas.NewText("boringchat", color.White)
					t.TextSize = 12
					t.TextStyle = fyne.TextStyle{Bold: true}
					return t
				}(),
			),
			container.NewPadded(
				roundedSquareButton("theme", color.NRGBA{R: 173, G: 216, B: 230, A: 230}, func() {
					currentTheme = BlueTheme
					applyTheme(currentTheme, panelBg, chatBg, avatar)
				}),
			),
			container.NewPadded(
				roundedSquareButton("theme", color.NRGBA{R: 255, G: 182, B: 193, A: 230}, func() {
					currentTheme = PinkTheme
					applyTheme(currentTheme, panelBg, chatBg, avatar)
				}),
			),
			container.NewPadded(
				roundedSquareButton("theme", color.NRGBA{R: 255, G: 255, B: 153, A: 230}, func() {
					currentTheme = YellowTheme
					applyTheme(currentTheme, panelBg, chatBg, avatar)
				}),
			),
			container.NewPadded(
				roundedSquareButton("theme", color.NRGBA{R: 152, G: 251, B: 152, A: 230}, func() {
					currentTheme = GreenTheme
					applyTheme(currentTheme, panelBg, chatBg, avatar)
				}),
			),
		),
	)

	historyLabel := widget.NewLabel("")
	historyLabel.Wrapping = fyne.TextWrapWord

	scrollContainer := container.NewScroll(historyLabel)
	scrollContainer.SetMinSize(fyne.NewSize(300, 360))
	scrollContainer.ScrollToBottom()

	inputEntry := widget.NewEntry()
	inputEntry.SetPlaceHolder("Сообщение")

	AllUsers := widget.NewLabel("")

	AllUsers.SetText(fmt.Sprintf("All Users: %v", name))

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
			AllUsersNames[mas[i].Name] = struct{}{}
			historyLabel.SetText(fmt.Sprintf("%v%v: %v\n", historyLabel.Text, mas[i].Name, mas[i].Text))
		}
		AllUsers.SetText("All Users: ")
		for u := range AllUsersNames {
			fmt.Println(u)
			AllUsers.SetText(fmt.Sprintf("%v%v; ", AllUsers.Text, u))
		}
	})

	inputEntry.OnSubmitted = func(text string) {
		var chars []string
		for _, r := range text {
			chars = append(chars, string(r))
		}
		if text == "" {
			return
		}
		if chars[0] == "/" {
			var command string
			for i := 1; i < len(chars); i++ {
				command += string(chars[i])
			}
			switch command {
			case "exit":
				os.Exit(0)
			}
		}

		historyLabel.SetText(fmt.Sprintf("%s%s: %s\n", historyLabel.Text, name, text))
		Writer(name, text)
		inputEntry.SetText("")
		scrollContainer.ScrollToBottom()
	}
	inputContainer := container.NewBorder(
		nil,
		nil,
		nil,
		update,
		inputEntry)

	AllUsersContainer := container.NewHBox(
		AllUsers,
	)
	mainContent := container.NewBorder(
		AllUsersContainer,
		inputContainer,
		nil,
		nil,
		scrollContainer)

	content := container.NewStack(
		chatBg,
		container.NewBorder(
			nil, nil,
			panel,
			nil,
			mainContent,
		),
	)

	chatWindow.SetContent(content)
	chatWindow.Show()
}
