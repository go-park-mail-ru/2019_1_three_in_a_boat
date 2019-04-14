package game

import (
	"math"
	"time"
)

const (
	StateWaiting = iota
	StateRunning
	StateOver
	StateDisconnect
)

type Input struct {
	Angle float64 `json:"angle"`
}

var InitialSnapshot = Snapshot{
	Angle:    math.Pi / 2,
	Hexagons: nil,
	State:    StateWaiting,
	Score:    0,
}

type Snapshot struct {
	Angle    float64
	Hexagons []Hexagon
	State    int
	Score    int64
}

func NewSnapshot() Snapshot {
	return InitialSnapshot
}

func (ss *Snapshot) Update(in Input) {

	// check previous snapshot doesn't end the game
	for _, h := range ss.Hexagons {
		if h.Crosses(ss.GetCursor()) {
			ss.State = StateOver
		}
	}

	// update snapshot
	for _, h := range ss.Hexagons {
		h.Shrink(ShrinkPerTick)
		h.Rotate(in.Angle - ss.Angle)
	}

	ss.Angle = in.Angle
}

func (ss *Snapshot) GetCursor() Circle {
	return Circle{
		Coords{
			X: Settings.PlayerCircleRadius * math.Cos(ss.Angle),
			Y: Settings.PlayerCircleRadius * math.Sin(ss.Angle),
		},
		Settings.CursorRadius,
	}
}

func init() {
	ShrinkPerTick = Settings.ShrinkPerSec / float64(time.Second/Settings.TickLength)
}

