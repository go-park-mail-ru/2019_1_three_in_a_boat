package game

import (
	"math"
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

var InitialHexagons = []struct {
	SidesMask int
	Side      float64
}{
	{SkipTop | SkipBottom, 25},
	{SkipTopRight | SkipBottomLeft, 20},
	{SkipTopLeft | SkipBottomRight, 15},
	{SkipTop | SkipBottom, 10},
}

type Snapshot struct {
	Angle    float64
	Hexagons []Hexagon
	State    int
	Score    int64
}

func NewSnapshot() Snapshot {
	snap := InitialSnapshot
	snap.Hexagons = make([]Hexagon, len(InitialHexagons))
	for i, h := range InitialHexagons {
		snap.Hexagons[i] = *NewHexagon(h.SidesMask, h.Side)
	}

	return snap
}

func (ss *Snapshot) Update(in Input) bool {

	// check previous snapshot doesn't end the game
	for _, h := range ss.Hexagons {
		if h.Crosses(ss.GetCursor()) {
			ss.State = StateOver
			return true
		}
	}

	// update snapshot
	for i, h := range ss.Hexagons {
		if h.Side <= ShrinkPerTick {
			ss.Hexagons[i] = *NewHexagon(InitialHexagons[0].SidesMask, InitialHexagons[0].Side)
		} else {
			h.Shrink(ShrinkPerTick)
			h.Rotate(in.Angle - ss.Angle)
		}
	}

	ss.Angle = in.Angle
	return false
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
