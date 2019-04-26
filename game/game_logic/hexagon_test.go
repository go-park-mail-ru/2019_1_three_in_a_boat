package game_logic

import (
	"fmt"
	"math"
	"testing"
)

func TestNewHexagon(t *testing.T) {
	h := NewHexagon(SkipBottomLeft, 10)
	eps := 1e-9

	for i, l := range h.Lines {
		if i != SideBottomLeft {
			if math.Abs(euclidianDistance(l.Start, l.End)-10) > eps {
				fmt.Println()
				t.Error("side length differs from the required one")
			}
		} else if l != nil {
			t.Error("side was required to be skipped but wasn't skipped")
		}
	}
}

func TestHexagon_Crosses(t *testing.T) {
	h := NewHexagon(0, 10)

	cases := []struct {
		Circle
		Crosses bool
	}{
		{
			Circle: Circle{
				Center: Coords{0, 0},
				Radius: 5,
			},
			Crosses: false,
		},
		{
			Circle: Circle{
				Center: Coords{50, 50},
				Radius: 15,
			},
			Crosses: false,
		},
		{
			Circle: Circle{
				Center: Coords{50, 50},
				Radius: 105,
			},
			Crosses: true,
		},
		{
			Circle: Circle{
				Center: Coords{5, -7},
				Radius: 3,
			},
			Crosses: true,
		},
	}

	for _, c := range cases {
		if crosses := h.Crosses(c.Circle); crosses != c.Crosses {
			t.Errorf("crosses error: expected %t, got %t", c.Crosses, crosses)
		}
	}
}

func euclidianDistance(one Coords, two Coords) float64 {
	return math.Sqrt((one.X-two.X)*(one.X-two.X) + (one.Y-two.Y)*(one.Y-two.Y))
}
