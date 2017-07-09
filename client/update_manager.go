package main

import (
	"fmt"
	"log"
	"net"

	"github.com/mmogo/mmo/shared"
)

const (
	maxBufferedUpdates = 30
)

type clientState struct {
	playerID string
	world    *shared.World
	updates  chan *shared.Update
	errc     chan error
}

func newClientState(id string, world *shared.World) *clientState {
	return &clientState{
		playerID: id,
		world:    world,
		updates:  make(chan *shared.Update, maxBufferedUpdates),
		errc:     make(chan error),
	}
}

func (cs *clientState) readUpdates(conn net.Conn) {
	loop := func() error {
		msg, err := shared.GetMessage(conn)
		if err != nil {
			return shared.FatalErr(err)
		}
		log.Println("RECV", msg)
		if msg.Error != nil {
			return fmt.Errorf("server returned an error: %v", msg.Error.Message)
		}
		if msg.Update != nil {
			cs.updates <- msg.Update
		}
		return nil
	}
	for {
		if err := loop(); err != nil {
			cs.errc <- err
			continue
		}
	}
}
