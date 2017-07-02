package shared

import (
	"image/color"
	"net"
	"sync"
	"time"

	"github.com/faiface/pixel"
)

type ServerPlayer struct {
	*Player
	Conn         net.Conn
	RequestQueue []*Message
	QueueLock    sync.RWMutex
	Ping         time.Time
}

type ClientPlayer struct {
	*Player
	Color color.Color
}

type Player struct {
	ID       string
	Position pixel.Vec
}
