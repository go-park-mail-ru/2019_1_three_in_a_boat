package game

import (
	"encoding/json"
	"go.uber.org/atomic"
	"math"
)

const (
	StateWaiting = iota
	StateRunning
	StateOver
	StateDisconnect
)

type Input struct {
	angle *atomic.Float64
}

func (i *Input) UnmarshalJSON(data []byte) error {
	var f float64
	err := json.Unmarshal(data, &f)
	if err != nil {
		return err
	}
	i.angle.Store(f)
	return nil
}

func (i *Input) Angle() float64 {
	return i.angle.Load()
}

func NewInput(angle float64) *Input {
	return &Input{atomic.NewFloat64(angle)}
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

func (ss *Snapshot) Update(in *Input) bool {

	// check previous snapshot doesn't end the game
	for _, h := range ss.Hexagons {
		if h.Crosses(ss.GetCursor()) {
			ss.State = StateOver
			return true
		}
	}

	// update snapshot
	for i := range ss.Hexagons {
		if ss.Hexagons[i].Side <= ShrinkPerTick {
			ss.Score += 10
			ss.Hexagons[i] = *NewHexagon(InitialHexagons[0].SidesMask, InitialHexagons[0].Side)
		} else {
			ss.Hexagons[i].Shrink(ShrinkPerTick)
			ss.Hexagons[i].Rotate(RotatePerTick)
		}
	}

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
