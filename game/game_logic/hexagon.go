package game_logic

import (
	"encoding/json"
)

// we're thinking about ⬣ sort of hexagon, not the rotated one (width > height)

const (
	SkipTop = 1 << iota
	SkipTopRight
	SkipBottomRight
	SkipBottom
	SkipBottomLeft
	SkipTopLeft
)

const (
	SideTop = iota
	SideTopRight
	SideBottomRight
	SideBottom
	SideBottomLeft
	SideTopLeft
)

const sqrt3 = 1.7320508075688772

type Hexagon struct {
	Side      float64
	SidesMask int
	angle     float64
	Lines     [6]*Line
}

type hexagonData struct {
	Side      float64 `json:"side"`
	SidesMask int     `json:"sidesMask"`
	Angle     float64 `json:"angle"`
}

func (h Hexagon) MarshalJSON() ([]byte, error) {
	return json.Marshal(hexagonData{h.Side, h.SidesMask, h.angle})
}

func NewHexagon(sidesMask int, side float64) *Hexagon {
	h := &Hexagon{Side: side, SidesMask: sidesMask}

	if sidesMask&SkipTop == 0 {
		h.Lines[SideTop] = &Line{
			Start: Coords{-h.Side / 2, h.Height() / 2},
			End:   Coords{h.Side / 2, h.Height() / 2},
		}
	} else {
		h.Lines[SideTop] = nil
	}

	if sidesMask&SkipTopRight == 0 {
		h.Lines[SideTopRight] = &Line{
			Start: Coords{h.Side / 2, h.Height() / 2},
			End:   Coords{h.Width() / 2, 0},
		}
	} else {
		h.Lines[SideTopRight] = nil
	}

	if sidesMask&SkipBottomRight == 0 {
		h.Lines[SideBottomRight] = &Line{
			Start: Coords{h.Width() / 2, 0},
			End:   Coords{h.Side / 2, -h.Height() / 2},
		}
	} else {
		h.Lines[SideBottomRight] = nil
	}

	if sidesMask&SkipBottom == 0 {
		h.Lines[SideBottom] = &Line{
			Start: Coords{h.Side / 2, -h.Height() / 2},
			End:   Coords{-h.Side / 2, -h.Height() / 2},
		}
	} else {
		h.Lines[SideBottom] = nil
	}

	if sidesMask&SkipBottomLeft == 0 {
		h.Lines[SideBottomLeft] = &Line{
			Start: Coords{-h.Side / 2, -h.Height() / 2},
			End:   Coords{-h.Width() / 2, 0},
		}
	} else {
		h.Lines[SideBottomLeft] = nil
	}

	if sidesMask&SkipTopLeft == 0 {
		h.Lines[SideTopLeft] = &Line{
			Start: Coords{-h.Width() / 2, 0},
			End:   Coords{-h.Side / 2, h.Height() / 2},
		}
	} else {
		h.Lines[SideTopLeft] = nil
	}

	return h
}

func (h *Hexagon) Height() float64 {
	return h.Side * sqrt3
}

func (h *Hexagon) Width() float64 {
	return 2 * h.Side
}

func (h *Hexagon) Crosses(circle Circle) bool {
	for _, l := range h.Lines {
		// lines have width, and we compensate for the with cursorRadius. This works
		// fine, however on the edges lines then appear a little longer than they
		// should be. E.g., it's possible to lose even though cursor didn't touch
		// the line. Shortening the line before checking fixes that. Shortened
		// does not affect the original line.
		if l != nil && l.Shortened(Settings.LineWidth).Crosses(circle) {
			return true
		}
	}
	return false
}

func (h *Hexagon) Rotate(rad float64) {
	h.angle = rad
	for _, line := range h.Lines {
		if line != nil {
			line.Rotate(h.angle)
		}
	}
}

func (h *Hexagon) Shrink(diff float64) {
	h.Side -= diff
	h.Lines = NewHexagon(h.SidesMask, h.Side).Lines
	h.Rotate(0) // sync with the h.angle
}
