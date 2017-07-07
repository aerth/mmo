package main

import (
	"log"

	"github.com/google/gxui"
	"github.com/google/gxui/drivers/gl"
	"github.com/google/gxui/themes/basic"
	"github.com/google/gxui/themes/dark"
)

func gui() {
	gl.StartDriver(patcherApp)
}

func patcherApp(driver gxui.Driver) {
	theme := dark.CreateTheme(driver)
	label1 := theme.CreateLabel()
	label1.SetColor(gxui.White)
	label1.SetText("")
	splitter := theme.CreateSplitterLayout()
	splitter.AddChild(label1)
	// input id
	textbox := basic.CreateTextBox(theme.(*basic.Theme))

	// errors
	errlabel := theme.CreateLabel()
	errlabel.SetColor(gxui.Red)
	splitter.AddChild(textbox)
	window := theme.CreateWindow(800, 600, "MMO")
	clientName := getClientName()
	downloadBtn := theme.CreateButton()
	downloadBtn.SetText("Check for update")
	downloadBtn.OnClick(func(gxui.MouseEvent) {
		if err := downloadClient(clientName); err != nil {
			errlabel.SetText(err.Error())
			log.Println(err)
		}

		// ungrey playBtn

	})

	playBtn := theme.CreateButton()
	playBtn.SetText("Play")

	playBtn.OnClick(func(gxui.MouseEvent) {
		window.Close()
		runClient(clientName)
	})

	splitter.AddChild(downloadBtn)
	splitter.AddChild(playBtn)
	splitter.AddChild(errlabel)

	window.AddChild(splitter)
	window.OnClose(driver.Terminate)

}
