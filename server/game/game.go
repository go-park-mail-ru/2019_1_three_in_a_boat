package game

import (
	"math"
	"sync"
	"time"
)

var Game game

var Settings = gameSettings{
	PlayerCircleRadius:    10,
	ShrinkPerSec:          0.1,
	RotatePerSec:          math.Pi / 4,
	TickDuration:          time.Millisecond * 50,
	IntensityToAngleRatio: 0.1,
	CursorRadius:          3,
}

var ShrinkPerTick float64
var RotatePerTick float64

type game struct {
	// the map stores pointers to Rooms, and despite being thread-safe
	// one room must not be modified concurrently, in the app a single room
	// is handled by a single goroutine
	Rooms sync.Map
}

type gameSettings struct {
	PlayerCircleRadius    float64
	ShrinkPerSec          float64
	RotatePerSec          float64
	TickDuration          time.Duration
	IntensityToAngleRatio float64
	CursorRadius          float64
}

func LoadOrStoreSinglePlayRoom(playerId int64, room Room) (Room, bool) {
	actual, ok := Game.Rooms.LoadOrStore(playerId, room)
	r := actual.(Room)
	return r, ok

}

func init() {
	ShrinkPerTick = Settings.ShrinkPerSec / float64(time.Second/Settings.TickDuration)
	RotatePerTick = Settings.RotatePerSec / float64(time.Second/Settings.TickDuration)
}
