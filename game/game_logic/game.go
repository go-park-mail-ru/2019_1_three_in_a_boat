package game_logic

import (
	"math"
	"sync"
	"time"
)

var Game game

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
	Rooms sync.Map
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

func LoadOrStoreSinglePlayRoom(room Room) (Room, bool) {
	actual, ok := Game.Rooms.LoadOrStore(string(room.Id()), room)
	r := actual.(Room)
	return r, ok

}

func init() {
	ShrinkPerTick = Settings.ShrinkPerSec / float64(time.Second/Settings.TickDuration)
	RotatePerTick = Settings.RotatePerSec / float64(time.Second/Settings.TickDuration)
	MultiplierPerTick = Settings.MultiplierPerSec / float64(time.Second/Settings.TickDuration)
	SameDirectionNumTicks = Settings.SameDirectionDuration * float64(time.Second/Settings.TickDuration)
}
