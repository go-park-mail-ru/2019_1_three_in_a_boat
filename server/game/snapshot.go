package game

import (
	"encoding/json"
	"math"
	"math/rand"

	"go.uber.org/atomic"
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
	Angle:             -math.Pi / 2,
	Hexagons:          nil,
	State:             StateWaiting,
	Score:             0,
	Ticks:             0,
	ClockWise:         true,
	CursorCircleAngle: 0,
}

var InitialHexagons = []struct {
	SidesMask int
	Side      float64
}{
	{SkipTop, 850},
	{SkipBottomLeft, 600},
	{SkipTopRight, 450},
	{SkipBottom, 300},
}

type Snapshot struct {
	Angle             float64
	Hexagons          []Hexagon
	State             int
	Score             int64
	Ticks             int64
	ClockWise         bool
	CursorCircleAngle float64
}

func NewSnapshot() Snapshot {
	snap := InitialSnapshot
	snap.Hexagons = make([]Hexagon, len(InitialHexagons))
	for i, h := range InitialHexagons {
		snap.Hexagons[i] = *NewHexagon(h.SidesMask, h.Side)
	}

	return snap
}

var masks = [6]int{SkipTop, SkipTopRight, SkipBottomRight, SkipBottom,
	SkipBottomLeft, SkipTopLeft}

func RandomSidesMask() int {
	return masks[1+rand.Intn(5)]
}

func RandomAngle() float64 {
	return rand.Float64() * math.Pi * 2
}

func (ss *Snapshot) Update(in *Input) bool {
	ss.Ticks += 1
	ss.Angle = in.Angle()
	// check previous snapshot doesn't end the game
	cur := ss.GetCursor()
	for _, h := range ss.Hexagons {
		if h.Crosses(cur) {
			ss.State = StateOver
			return true
		}
	}

	difficultyIncrement := 1 + MultiplierPerTick*float64(ss.Ticks)
	if difficultyIncrement > Settings.MaxMultiplier {
		difficultyIncrement = Settings.MaxMultiplier
	}

	ticksSinceRotation := float64(ss.Ticks % int64(math.Round(SameDirectionNumTicks)))
	if ticksSinceRotation == 0 {
		ss.ClockWise = !ss.ClockWise
	}

	var rotationAmplitude float64
	// to or from rotation really
	ticksToRotation := math.Min(ticksSinceRotation, math.Abs(ticksSinceRotation-SameDirectionNumTicks))
	rotationAmplitude = math.Min(1, math.Max(0.5, 3*ticksToRotation/SameDirectionNumTicks))

	// 1 = clockwise, -1 = ccw
	var rotationDirection float64 = 1
	if !ss.ClockWise {
		rotationDirection = -1
	}

	angleIncrement := rotationAmplitude * rotationDirection * RotatePerTick * difficultyIncrement
	// update snapshot
	for i := range ss.Hexagons {
		if ss.Hexagons[i].Side <= Settings.MinHexagonSize {
			ss.Score += 10
			ss.Hexagons[i] = *NewHexagon(RandomSidesMask(), InitialHexagons[0].Side)
			ss.Hexagons[i].Rotate(RandomAngle())
		} else {
			angle := ss.Hexagons[i].angle + angleIncrement
			ss.Hexagons[i].Shrink(ShrinkPerTick * difficultyIncrement)
			ss.Hexagons[i].Rotate(angle)
		}
	}
	ss.CursorCircleAngle += angleIncrement

	return false
}

func (ss *Snapshot) GetCursor() Circle {
	return Circle{
		Coords{
			X: Settings.PlayerCircleRadius * math.Cos(ss.Angle+ss.CursorCircleAngle),
			Y: Settings.PlayerCircleRadius * math.Sin(ss.Angle+ss.CursorCircleAngle),
		},
		Settings.CursorRadius,
	}
}
