package main

import (
	"fmt"
	"log"

	"net"

	"github.com/mmogo/mmo/shared"
	"github.com/ilackarms/pkg/errors"
)

type requestManager struct {
	pendingRequests   <-chan *shared.Request
	updatePredictions chan *shared.Update
	conn              net.Conn
}

func (mgr *requestManager) processPending() error {
requestLoop:
	for {
		select {
		default:
			break requestLoop
		case req := <-mgr.pendingRequests:
			if err := mgr.handleRequest(req); err != nil {
				log.Printf("Error handling player request %#v: %v", req, err)
			}
		}
	}
	return nil
}

func (mgr *requestManager) handleRequest(id string, req *shared.Request) error {
	if err := shared.SendMessage(&shared.Message{Request: req}, mgr.conn); err != nil {
		return errors.New("failed to send request", err)
	}
	switch {
	case req.MoveRequest != nil:
		return mgr.playerMoved(player, req.MoveRequest)
	case req.SpeakRequest != nil:
		return mgr.playerSpoke(&shared.PlayerSpoke{
			ID:   player.ID,
			Text: req.SpeakRequest.Text,
		})
	default:
		return fmt.Errorf("unknown request type: %#v", req)
	}
}
