package game

// we're thinking about ⬣ sort of hexagon, not the rotated one (width > height)

type Coords struct {
	X float64
	Y float64
}

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
	Center     Coords
	Side       float64
	EmptySides map[int]struct{}
	Lines      [6]*Line
}

func NewHexagon(skipSides []int, side float64) *Hexagon {
	h := &Hexagon{Center: Coords{0, 0}, Side: side}

	for _, x := range skipSides {
		h.EmptySides[x] = struct{}{}
	}

	if _, skip := h.EmptySides[SideTop]; !skip {
		h.Lines[SideTop] = &Line{
			Start: Coords{-h.Side / 2, h.Height() / 2},
			End:   Coords{h.Side / 2, h.Height() / 2},
		}
	} else {
		h.Lines[SideTop] = nil
	}

	if _, skip := h.EmptySides[SideTopRight]; !skip {
		h.Lines[SideTopRight] = &Line{
			Start: Coords{h.Side / 2, h.Height() / 2},
			End:   Coords{h.Width() / 2, 0},
		}
	} else {
		h.Lines[SideTopRight] = nil
	}

	if _, skip := h.EmptySides[SideBottomRight]; !skip {
		h.Lines[SideBottomRight] = &Line{
			Start: Coords{h.Width() / 2, 0},
			End:   Coords{h.Side / 2, -h.Height() / 2},
		}
	} else {
		h.Lines[SideBottomRight] = nil
	}

	if _, skip := h.EmptySides[SideBottom]; !skip {
		h.Lines[SideBottom] = &Line{
			Start: Coords{h.Side / 2, -h.Height() / 2},
			End:   Coords{-h.Side / 2, -h.Height() / 2},
		}
	} else {
		h.Lines[SideBottom] = nil
	}

	if _, skip := h.EmptySides[SideBottomLeft]; !skip {
		h.Lines[SideBottomLeft] = &Line{
			Start: Coords{-h.Side / 2, -h.Height() / 2},
			End:   Coords{-h.Width() / 2, 0},
		}
	} else {
		h.Lines[SideBottomLeft] = nil
	}

	if _, skip := h.EmptySides[SideTopLeft]; !skip {
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
		if l != nil && l.Crosses(circle) {
			return true
		}
	}
	return false
}
