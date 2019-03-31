package game

import "math"

type Circle struct {
	Center Coords
	Radius float64
}

type Line struct {
	Start Coords
	End   Coords
}

func (l Line) Crosses(circle Circle) bool {
	// https://math.stackexchange.com/questions/275529/check-if-line-intersects-with-circles-perimeter
	x0 := circle.Center.X
	y0 := circle.Center.Y
	x1 := l.Start.X
	y1 := l.Start.Y
	x2 := l.End.X
	y2 := l.End.Y
	return math.Abs((x2-x1)*x0+(y1-y2)*y0+(x1-x2)*y1+x1*(y2-y1))/
		math.Sqrt((x2-x1)*(x2-x1)+(y1-y2)*(y1-y2)) <= circle.Radius
}
