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

type inputData struct {
	Angle float64 `json:"angle"`
}

func (i *Input) UnmarshalJSON(data []byte) error {
	f := inputData{}
	err := json.Unmarshal(data, &f)
	if err != nil {
		return err
	}
	// invert the angle because of some stupid canvas shitто
	i.angle.Store(-f.Angle)
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
	{SkipTop | SkipBottom, 850},
	{SkipTopRight | SkipBottomLeft, 600},
	{SkipTopLeft | SkipBottomRight, 450},
	{SkipTop | SkipBottom, 300},
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
	ss.Angle = in.Angle()
	// check previous snapshot doesn't end the game
	cur := ss.GetCursor()
	for _, h := range ss.Hexagons {
		if h.Crosses(cur) {
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
			angle := ss.Hexagons[i].angle + RotatePerTick
			ss.Hexagons[i].Shrink(ShrinkPerTick)
			ss.Hexagons[i].Rotate(angle)
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
