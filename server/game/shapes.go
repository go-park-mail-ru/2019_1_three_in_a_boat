package game

import "math"

type Coords struct {
	X float64
	Y float64
}

func (c Coords) Sub(c2 Coords) Coords {
	c.X -= c2.X
	c.Y -= c2.Y
	return c
}

func (c Coords) Add(c2 Coords) Coords {
	c.X += c2.X
	c.Y += c2.Y
	return c
}

func (c Coords) Dot(c2 Coords) float64 {
	return c.X*c2.X + c.Y*c2.Y
}

type Circle struct {
	Center Coords
	Radius float64
}

type Line struct {
	Start Coords
	End   Coords
}

func (l Line) Crosses(circle Circle) bool {
	// programmersforum for the rescue bljad
	// http://www.programmersforum.ru/showthread.php?t=117078
	eps := 1e-10
	d01 := l.Start.Sub(circle.Center)
	d12 := l.End.Sub(l.Start)

	// solving at + 2kt + c = 0 for t
	a := d12.Dot(d12)
	k := d01.Dot(d12)
	c := d01.X*d01.X + d01.Y*d01.Y - circle.Radius*circle.Radius
	disc := k*k - a*c

	if disc < 0 {
		return false
	} else if math.Abs(disc) < eps {
		t := -k / a
		return t > 0-eps && t < 1+eps
	} else {
		t1 := (-k + math.Sqrt(disc)) / a
		t2 := (-k - math.Sqrt(disc)) / a
		if t1 > t2 {
			t1, t2 = t2, t1
		}
		if t1 >= 0-eps && t2 <= 1+eps {
			return true
		} else if t2 <= 0+eps || t1 >= 1-eps {
			return false
		} else {
			return true
		}

		// return !(t1 >= 0 - eps && t2 <= 1 + eps) && (t2 > 0+eps && t1 < 1-eps)
	}
}

func (l Line) norm() float64 {
	return math.Sqrt((l.Start.X-l.End.X)*(l.Start.X-l.End.X) +
		(l.Start.Y-l.End.Y)*(l.Start.Y-l.End.Y))
}

func (l Line) normVector() Coords {
	norm := l.norm()
	return Coords{
		X: (l.End.X - l.Start.X) / norm,
		Y: (l.End.Y - l.Start.Y) / norm,
	}
}

func (l Line) distance(circle Circle) float64 {
	nv := l.normVector()
	ctl := Coords{ // centerToLineStart
		circle.Center.X - l.Start.X,
		circle.Center.Y - l.Start.Y,
	}
	return math.Abs(nv.X*ctl.Y - nv.Y*ctl.X)
}
