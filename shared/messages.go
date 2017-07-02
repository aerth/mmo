package shared

import (
	"fmt"
	"time"

	"github.com/faiface/pixel"
)

type Message struct {
	Sent    time.Time
	Ping    *time.Duration
	Request *Request
	Update  *Update
}

type Update struct {
	PlayerMoved        *PlayerMoved
	PlayerSpoke        *PlayerSpoke
	WorldState         *WorldState
	PlayerDisconnected *PlayerDisconnected
}

type Request struct {
	ConnectRequest *ConnectRequest
	MoveRequest    *MoveRequest
	SpeakRequest   *SpeakRequest
}

type ConnectRequest struct {
	ID string
}

type MoveRequest struct {
	Direction Direction
	Created   time.Time
}

type SpeakRequest struct {
	Text string
}

type PlayerMoved struct {
	ID          string
	NewPosition pixel.Vec
	RequestTime time.Time
}

type PlayerSpoke struct {
	ID   string
	Text string
}

type WorldState struct {
	Players []*Player
}

type PlayerDisconnected struct {
	ID string
}

func (m Message) String() string {
	if m.Request != nil {
		return m.Request.String()
	}

	if m.Update != nil {
		return m.Update.String()
	}

	if m.Sent.IsZero() {
		return "invalid packet"
	}

	if m.Ping != nil {
		return fmt.Sprintf("Ping: %s", m.Ping)
	}

	return "empty packet"
}

func (u Update) String() string {
	if u.PlayerMoved != nil {
		return fmt.Sprintf("PlayerMoved: %s: %s", u.PlayerMoved.ID, u.PlayerMoved.NewPosition)
	}

	if u.PlayerSpoke != nil {
		return fmt.Sprintf("PlayerSpoke: %s: %s", u.PlayerSpoke.ID, u.PlayerSpoke.Text)
	}

	if u.WorldState != nil {

		return fmt.Sprintf("WorldState: %v players", len(u.WorldState.Players))
	}
	if u.PlayerDisconnected != nil {
		return fmt.Sprintf("PlayerDisconnected: %s", u.PlayerDisconnected)
	}

	return "empty update"

}

func (r Request) String() string {
	if r.ConnectRequest != nil {
		return fmt.Sprintf("ConnectRequest: %v", r.ConnectRequest.ID)
	}
	if r.MoveRequest != nil {
		return fmt.Sprintf("MoveRequest: %s", r.MoveRequest.Direction)
	}
	if r.SpeakRequest != nil {
		return fmt.Sprintf("SpeakRequest: %s", r.SpeakRequest.Text)
	}

	return "empty request"
}
