package main

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/mmogo/mmo/shared"
)

type inputProcessor struct {
	win            *pixelgl.Window
	centerMatrix   pixel.Matrix
	requestsToSend chan *shared.Request
}

func newInputProcessor(win *pixelgl.Window, centerMatrix pixel.Matrix, requestsToSend chan *shared.Request) *inputProcessor {
	return &inputProcessor{
		win:            win,
		centerMatrix:   centerMatrix,
		requestsToSend: requestsToSend,
	}
}

func (ip *inputProcessor) handleInputs() error {
	if win.Pressed(pixelgl.MouseButtonLeft) {
		mouse := g.centerMatrix.Unproject(win.MousePosition())
		loc := g.players[g.playerID].Position
		g.queueSimulation(func() {
			g.setPlayerPosition(g.playerID, loc.Add(mouse.Unit().Scaled(2)))
		})

		// send to server
		if err := requestMove(mouse.Unit().Scaled(2), conn); err != nil {
			return err
		}
	}

}
