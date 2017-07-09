package main

import (
	"fmt"
	"log"

	"net"

	"github.com/ilackarms/pkg/errors"
	"github.com/mmogo/mmo/shared"
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
			if err := mgr.handleRequest("player1", req); err != nil {
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
		return mgr.playerMoved(id, req.MoveRequest)
	case req.SpeakRequest != nil:
		return mgr.playerSpoke(&shared.PlayerSpoke{
			ID:   id,
			Text: req.SpeakRequest.Text,
		})
	default:
		return fmt.Errorf("unknown request type: %#v", req)
	}
}
func (mgr *requestManager) playerMoved(id string, req *shared.MoveRequest) error { return nil }
func (mgr *requestManager) playerSpoke(req *shared.PlayerSpoke) error            { return nil }
