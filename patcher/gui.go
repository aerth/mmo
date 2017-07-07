package main

import (
	"log"

	"github.com/google/gxui"
	"github.com/google/gxui/drivers/gl"
	"github.com/google/gxui/gxfont"
	"github.com/google/gxui/themes/basic"
	"github.com/google/gxui/themes/dark"
)

var errtext string
var idtext string
var srvtext string

func gui() {
	idtext = *playerID
	srvtext = *addr
	gl.StartDriver(patcherApp)
}

func patcherApp(driver gxui.Driver) {
	theme := dark.CreateTheme(driver)
	splitter := theme.CreateSplitterLayout()
	// input id
	header := theme.CreateLabel()
	font, err := driver.CreateFont(gxfont.Default, 75)
	if err != nil {
		panic(err)
	}
	header.SetFont(font)
	header.SetText("MMO")
	splitter.AddChild(header)
	idlabel := theme.CreateLabel()
	idlabel.SetText("ID")
	idbox := basic.CreateTextBox(theme.(*basic.Theme))
	idbox.SetText(idtext)
	log.Println(idbox.Text())
	srvlabel := theme.CreateLabel()
	srvlabel.SetText("Server")
	addrbox := basic.CreateTextBox(theme.(*basic.Theme))
	addrbox.SetText(srvtext)
	log.Println(addrbox.Text())
	splitter.AddChild(idlabel)
	splitter.AddChild(idbox)
	splitter.AddChild(srvlabel)
	splitter.AddChild(addrbox)
	// errors
	errlabel := theme.CreateLabel()
	errlabel.SetColor(gxui.Red)
	window := theme.CreateWindow(800, 600, "MMO")
	clientName := getClientName()
	downloadBtn := theme.CreateButton()
	downloadBtn.SetText("Check for update")
	downloadBtn.OnClick(func(gxui.MouseEvent) {
		*addr = addrbox.Text()
		if err := downloadClient(clientName); err != nil {
			errtext += "\n" + err.Error()
			errtext = wrap(errtext)
			errlabel.SetText(errtext)
			log.Println(err)
		}

		// ungrey playBtn

	})

	playBtn := theme.CreateButton()
	playBtn.SetText("Play")

	playBtn.OnClick(func(gxui.MouseEvent) {
		*addr = addrbox.Text()
		*playerID = idbox.Text()
		window.Close()
		runClient(clientName)
	})

	splitter.AddChild(downloadBtn)
	splitter.AddChild(playBtn)
	splitter.AddChild(errlabel)

	window.AddChild(splitter)
	window.OnClose(driver.Terminate)

}

func wrap(s string) string {
	var out string
	for i, c := range s {
		if i%4 == 0 {
			out += "\n"
		}
		out += string(c)
	}

	return out
}
