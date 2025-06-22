package internal

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"image/color"
	"os"
	"time"
)

var AllUsersNames map[string]struct{}

func Ui() {
	myApp := app.New()

	nameWindow := myApp.NewWindow("welcome")
	nameWindow.Resize(fyne.NewSize(400, 400))

	bgImage := canvas.NewImageFromFile("internal/images/fish.jpg")
	bgImage.FillMode = canvas.ImageFillStretch

	var UserName string
	entry := widget.NewEntry()
	entry.SetPlaceHolder("Enter your name:")
	entryIP := widget.NewEntry()
	entryIP.SetPlaceHolder("Enter IP of your server:")

	submitBtn := createImageButton("internal/images/greenbtn.png", func() {
		UserName = entry.Text
		if UserName == "" {
			UserName = "unnamed guy"
		}
		if UserName == "Oleg" {
			UserName = "silly guy"
		}
		IP = fmt.Sprintf("%v:8080", entryIP.Text)
		nameWindow.Close()
		createChatWindow(myApp, UserName)
	})

	form := container.NewVBox(
		entry,
		entryIP,
		container.NewCenter(submitBtn),
		container.NewCenter(
			canvas.NewText("                                                ", color.White),
		),
	)

	nameWindow.SetContent(container.NewStack(
		bgImage,
		container.NewCenter(form),
	))

	icon, err := fyne.LoadResourceFromPath("logo256.png")
	if err != nil {
		// Если файл не найден, используем встроенную иконку Fyne
		icon = theme.FyneLogo()
	}

	nameWindow.SetIcon(icon)
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
		color.NRGBA{R: 0, G: 100, B: 255, A: 230},
		color.NRGBA{R: 0, G: 100, B: 255, A: 128},
		"internal/images/blue.png",
	)
	PinkTheme = Theme(
		color.NRGBA{R: 255, G: 190, B: 240, A: 200},
		color.NRGBA{R: 255, G: 150, B: 200, A: 200},
		"internal/images/pink.png",
	)
	GreenTheme = Theme(
		color.NRGBA{R: 152, G: 251, B: 152, A: 230},
		color.NRGBA{R: 0, G: 180, B: 100, A: 255},
		"internal/images/green.png",
	)
	YellowTheme = Theme(
		color.NRGBA{R: 255, G: 220, B: 110, A: 235},
		color.NRGBA{R: 255, G: 180, B: 70, A: 235},
		"internal/images/yellow.png",
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

	currentTheme := BlueTheme

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
				roundedSquareButton("theme", color.NRGBA{R: 0, G: 100, B: 255, A: 230}, func() {
					currentTheme = BlueTheme
					applyTheme(currentTheme, panelBg, chatBg, avatar)
				}),
			),
			container.NewPadded(
				roundedSquareButton("theme", color.NRGBA{R: 255, G: 190, B: 240, A: 200}, func() {
					currentTheme = PinkTheme
					applyTheme(currentTheme, panelBg, chatBg, avatar)
				}),
			),
			container.NewPadded(
				roundedSquareButton("theme", color.NRGBA{R: 255, G: 220, B: 110, A: 235}, func() {
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

	//scrollContainer := container.NewScroll(historyLabel)
	//scrollContainer.SetMinSize(fyne.NewSize(300, 360))
	//scrollContainer.ScrollToBottom()

	inputEntry := widget.NewEntry()
	inputEntry.SetPlaceHolder("Сообщение")

	AllUsers := widget.NewLabel("")

	AllUsers.SetText(fmt.Sprintf("All Users: %v", name))

	var mas []Pair

	go func() {
		ticker := time.NewTicker(2 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				mas, _, AllUsersNames = Writing()
				if len(mas) == 0 {
					continue
				} else {
					historyLabel.SetText("")
				}
				for {
					if len(mas) >= 14 {
						mas = mas[13:]
					} else {
						break
					}
				}
				for i := 0; i < len(mas); i++ {
					historyLabel.SetText(fmt.Sprintf("%v%v: %v\n", historyLabel.Text, mas[i].Name, mas[i].Text))
				}
				AllUsers.SetText("All Users: ")
				for u := range AllUsersNames {
					if u == "" {
						continue
					}
					AllUsers.SetText(fmt.Sprintf("%v%v ", AllUsers.Text, u))
				}
			}
		}
	}()

	//update := widget.NewButton("Появились сообщения", func() {
	//	mas, _, AllUsersNames = Writing()
	//	if len(mas) == 0 {
	//		return
	//	} else {
	//		historyLabel.SetText("")
	//	}
	//	for i := 0; i < len(mas); i++ {
	//		historyLabel.SetText(fmt.Sprintf("%v%v: %v\n", historyLabel.Text, mas[i].Name, mas[i].Text))
	//	}
	//	AllUsers.SetText("All Users: ")
	//	for u := range AllUsersNames {
	//		if u == "" {
	//			continue
	//		}
	//		AllUsers.SetText(fmt.Sprintf("%v%v ", AllUsers.Text, u))
	//	}
	//})

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
				Writer(name, "EXIT")
				os.Exit(0)
			}
		}

		historyLabel.SetText(fmt.Sprintf("%s%s: %s\n", historyLabel.Text, name, text))
		Writer(name, text)
		inputEntry.SetText("")
		//scrollContainer.ScrollToBottom()
	}

	inputContainer := container.NewBorder(
		nil,
		nil,
		nil,
		nil,
		inputEntry)

	AllUsersContainer := container.NewHBox(
		AllUsers,
	)
	mainContent := container.NewBorder(
		AllUsersContainer,
		inputContainer,
		nil,
		nil,
		historyLabel)

	content := container.NewStack(
		chatBg,
		container.NewBorder(
			nil, nil,
			panel,
			nil,
			mainContent,
		),
	)

	chatWindow.SetOnClosed(func() {
		Writer(name, "EXIT")
		os.Exit(0)
	})

	icon, err := fyne.LoadResourceFromPath("logo256.png")
	if err != nil {
		// Если файл не найден, используем встроенную иконку Fyne
		icon = theme.FyneLogo()
	}
	chatWindow.SetIcon(icon)
	chatWindow.SetContent(content)
	chatWindow.Show()
}
