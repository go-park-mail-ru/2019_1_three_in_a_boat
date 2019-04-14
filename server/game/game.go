package game

import (
	"math"
	"sync"
	"time"
)

var Game game

var Settings = settings{
	PlayerCircleRadius:    5,
	ShrinkPerSec:          5,
	RotatePerSec:          math.Pi / 4,
	TickDuration:          time.Millisecond * 35,
	IntensityToAngleRatio: 0.1,
	CursorRadius:          0.5,
}

var ShrinkPerTick float64

type game struct {
	// the map stores pointers to Rooms, and despite being thread-safe
	// one room must not be modified concurrently, in the app a single room
	// is handled by a single goroutine
	Rooms sync.Map
}

type settings struct {
	PlayerCircleRadius    float64
	ShrinkPerSec          float64
	RotatePerSec          float64
	TickDuration          time.Duration
	IntensityToAngleRatio float64
	CursorRadius          float64
}

func GetSinglePlayRoom(playerId int64) (Room, bool) {
	if room, ok := Game.Rooms.Load(playerId); ok {
		r := room.(Room)
		return r, true
	} else {
		return nil, false
	}
}

func LoadOrStoreSinglePlayRoom(playerId int64, room Room) (Room, bool) {
	actual, ok := Game.Rooms.LoadOrStore(playerId, room)
	r := actual.(Room)
	return r, ok

}

func init() {
	ShrinkPerTick = Settings.ShrinkPerSec / float64(time.Second/Settings.TickDuration)
}
