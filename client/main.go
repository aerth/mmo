package main

import (
	_ "image/png"

	"flag"
	"fmt"
	"image/color"
	"log"
	"math"
	"net"
	"time"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
	"github.com/ilackarms/pkg/errors"
	"github.com/mmogo/mmo/shared"
	"github.com/xtaci/smux"
	"golang.org/x/image/colornames"
	"golang.org/x/image/font/basicfont"
)

const (
	UP        = shared.UP
	DOWN      = shared.DOWN
	LEFT      = shared.LEFT
	RIGHT     = shared.RIGHT
	UPLEFT    = shared.UPLEFT
	UPRIGHT   = shared.UPRIGHT
	DOWNLEFT  = shared.DOWNLEFT
	DOWNRIGHT = shared.DOWNRIGHT
)

func init() {
	log.SetFlags(log.Lmicroseconds | log.Lshortfile)
}

func main() {
	addr := flag.String("addr", "localhost:8080", "address of server")
	id := flag.String("id", "", "playerid to use")
	protocol := flag.String("protocol", "udp", fmt.Sprintf("network protocol to use. available %s | %s", shared.ProtocolTCP, shared.ProtocolUDP))
	flag.Parse()
	if *id == "" {
		log.Fatal("id must be provided")
	}
	pixelgl.Run(func() {
		if err := run(*protocol, *addr, *id); err != nil {
			log.Fatal(err)
		}
	})
}

func run(protocol, addr, id string) error {
	conn, err := dialServer(protocol, addr, id)
	if err != nil {
		return errors.New("failed to dial server", err)
	}

	msg, err := shared.GetMessage(conn)
	if err != nil {
		return errors.New("failed reading message", err)
	}

	if msg.Update == nil || msg.Update.WorldState == nil || msg.Update.WorldState.World == nil {
		return errors.New("expected Sync message on server handshake, got "+msg.String(), nil)
	}

	cs := newClientState(id, msg.Update.WorldState.World)

	go func() { cs.readUpdates(conn) }()

	cfg := pixelgl.WindowConfig{
		Title:  "loading",
		Bounds: pixel.R(0, 0, 800, 600),
		VSync:  true,
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		return fmt.Errorf("creating window: %v", err)
	}

	// load assets
	lootImage, err := loadPicture("sprites/loot.png")
	if err != nil {
		return err
	}
	lootSprite := pixel.NewSprite(lootImage, lootImage.Bounds())

	playerSprite, err := LoadSpriteSheet("sprites/char1.png", nil)
	if err != nil {
		return shared.FatalErr(err)
	}

	fps := 0 // calculated frames per second
	second := time.Tick(time.Second)
	ping := time.Tick(time.Second * 2)
	last := time.Now()
	atlas := text.NewAtlas(basicfont.Face7x13, text.ASCII)
	tilebatch := debugTiles()
	wincenter := win.Bounds().Center()
	centerMatrix := pixel.IM.Moved(wincenter)
	go func() {
		for {
			err := <-cs.errc
			if shared.IsFatal(err) {
				log.Fatal(err)
			}
			log.Printf("Error: %v", err)
		}
	}()
	camPos := pixel.ZV
	playerText := text.New(pixel.ZV, atlas)
	requestsToSend := make(chan *shared.Request)
	input := newInputProcessor(win, centerMatrix, requestsToSend)
	for !win.Closed() {
		win.Clear(colornames.Yellow)
		dt := time.Since(last).Seconds()
		last = time.Now()

		err := input.Process()
		if err != nil {
			log.Println(err)
		}
		//g.applySimulations()

		//		playerSprite.Animate(dt, g.facing, g.action)

		tilebatch.Draw(win)

		lootSprite.Draw(win, pixel.IM.Scaled(pixel.ZV, 2.0))
		//g.lock.RLock()
		me, ok := cs.world.GetPlayer(id)
		if !ok {
			log.Println("didnt make it to new world")
			return shared.FatalErr(fmt.Errorf("cant find self"))
		}
		pos := me.Position
		camPos = pixel.Lerp(camPos, win.Bounds().Center().Sub(pos), 1-math.Pow(1.0/128, dt))
		cam := pixel.IM.Moved(camPos)
		win.SetMatrix(cam)
		for _, id := range cs.world.Players() {
			player, ok := cs.world.GetPlayer(id)
			if !ok {
				log.Println("player 404")
				continue
			}
			playerPos := pixel.IM.Moved(player.Position)
			playerSprite.Draw(win, playerPos, colornames.Blue)
		}

		// show mouse coordinates
		mousePos := cam.Unproject(win.MousePosition())
		playerText.Clear()
		playerText.Dot = playerText.Orig
		playerText.WriteString(fmt.Sprintf("%v", mousePos))
		playerText.DrawColorMask(win, pixel.IM.Moved(mousePos), colornames.White)

		win.Update()

		fps++
		select {
		default:
		case <-ping:
			shared.SendMessage(&shared.Message{}, conn)
		}
		select {
		default:
		case <-second:
			win.SetTitle(fmt.Sprintf("%v fps", fps))
			fps = 0
		}
	}
	return nil
}

func dialServer(protocol, addr, id string) (net.Conn, error) {
	log.Printf("dialing %s", addr)
	conn, err := shared.Dial(protocol, addr)
	if err != nil {
		return nil, err
	}
	session, err := smux.Client(conn, smux.DefaultConfig())
	if err != nil {
		return nil, err
	}
	stream, err := session.OpenStream()
	if err != nil {
		return nil, err
	}
	conn = stream

	if err := shared.SendMessage(&shared.Message{
		Request: &shared.Request{
			ConnectRequest: &shared.ConnectRequest{
				ID: id,
			},
		}}, conn); err != nil {
		return nil, err
	}
	return conn, nil
}

func requestMove(direction pixel.Vec, conn net.Conn) error {
	return shared.SendMessage(&shared.Message{
		Request: &shared.Request{
			MoveRequest: &shared.MoveRequest{
				Direction: direction,
				Created:   time.Now(),
			},
		},
	}, conn)
}

func requestSpeak(txt string, conn net.Conn) error {
	return shared.SendMessage(&shared.Message{
		Request: &shared.Request{
			SpeakRequest: &shared.SpeakRequest{
				Text: txt,
			},
		},
	}, conn)
}

func stringToColor(str string) color.Color {
	colornum := 0
	for _, s := range str {
		colornum += int(s)
	}
	all := len(colornames.Names)
	name := colornames.Names[colornum%all]
	return colornames.Map[name]
}
