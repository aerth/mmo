package main

import (
	"log"
	"time"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/mmogo/mmo/shared"
)

type inputProcessor struct {
	win          *pixelgl.Window
	centerMatrix pixel.Matrix
	req          chan *shared.Request
}

func newInputProcessor(win *pixelgl.Window, centerMatrix pixel.Matrix, requestsToSend chan *shared.Request) *inputProcessor {
	return &inputProcessor{
		win:          win,
		centerMatrix: centerMatrix,
		req:          requestsToSend,
	}
}

func (ip *inputProcessor) Process() error {
	if ip.win.Pressed(pixelgl.MouseButtonLeft) {
		mouse := ip.centerMatrix.Unproject(ip.win.MousePosition())
		ip.req <- &shared.Request{MoveRequest: &shared.MoveRequest{
			Direction: mouse.Unit(),
			Created:   time.Now(),
		}}
		log.Println("mouse %s unit %s", mouse, mouse.Unit())
	}

	return nil
}
