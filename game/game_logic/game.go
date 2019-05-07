package game_logic

import (
	"math"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

var Game = game{
	Rooms:        sync.Map{},
	WaitingRooms: nil,
}

var Settings = gameSettings{
	PlayerCircleRadius:    70,
	ShrinkPerSec:          100,
	RotatePerSec:          math.Pi / 3,
	TickDuration:          time.Millisecond * 40,
	IntensityToAngleRatio: 0.1,
	CursorRadius:          10,
	LineWidth:             5,
	MinHexagonSize:        40,
	MultiplierPerSec:      1e-2,
	MaxMultiplier:         2.5,
	SameDirectionDuration: 15, // seconds
}

var ShrinkPerTick float64
var RotatePerTick float64
var MultiplierPerTick float64
var SameDirectionNumTicks float64

type game struct {
	// the map stores pointers to Rooms, and despite being thread-safe
	// one room must not be modified concurrently, in the app a single room
	// is handled by a single goroutine
	Rooms        sync.Map
	WaitingRooms []*MultiPlayerRoom
	mu           sync.Mutex
}

type gameSettings struct {
	PlayerCircleRadius    float64
	ShrinkPerSec          float64
	RotatePerSec          float64
	TickDuration          time.Duration
	IntensityToAngleRatio float64
	CursorRadius          float64
	MinHexagonSize        float64
	LineWidth             float64
	// controls how fast does the shrinking/rotating speed increase
	MultiplierPerSec      float64
	MaxMultiplier         float64
	SameDirectionDuration float64
}

func LoadOrStoreRoom(room Room) (Room, bool) {
	actual, ok := Game.Rooms.LoadOrStore(string(room.Id()), room)
	r := actual.(Room)
	return r, ok
}

// TODO: send a test message to make sure the socket hasn't disconnected
func GetOrCreateMPRoom(
	r *http.Request, uid int64, conn *websocket.Conn) (*MultiPlayerRoom, bool) {
	Game.mu.Lock()
	defer Game.mu.Unlock()
	l := len(Game.WaitingRooms)
	if l != 0 {
		first := Game.WaitingRooms[0]
		Game.WaitingRooms = Game.WaitingRooms[1:]
		return first, true
	} else {
		Game.WaitingRooms = append(Game.WaitingRooms, NewMultiPlayerRoom(r, uid, conn))
		return Game.WaitingRooms[0], false
	}
}

func init() {
	ShrinkPerTick = Settings.ShrinkPerSec / float64(time.Second/Settings.TickDuration)
	RotatePerTick = Settings.RotatePerSec / float64(time.Second/Settings.TickDuration)
	MultiplierPerTick = Settings.MultiplierPerSec / float64(time.Second/Settings.TickDuration)
	SameDirectionNumTicks = Settings.SameDirectionDuration * float64(time.Second/Settings.TickDuration)
}
