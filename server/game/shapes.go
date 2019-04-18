package game

import "math"

type Coords struct {
	X float64
	Y float64
}

func (c *Coords) ApplyTransform(T [2][2]float64) {
	tmpX := c.X
	c.X = T[0][0]*c.X + T[0][1]*c.Y
	c.Y = T[1][0]*tmpX + T[1][1]*c.Y
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

func (l *Line) Crosses(circle Circle) bool {
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

		// if t1 >= 0-eps && t2 <= 1+eps {
		// 	return true
		// } else if t2 <= 0+eps || t1 >= 1-eps {
		// 	return false
		// } else {
		// 	return true
		// }

		return !(t1 < 0-eps || t2 > 1+eps) || (t2 > 0+eps && t1 < 1-eps)
	}
}

func (l *Line) Rotate(rad float64) {
	sin, cos := math.Sincos(rad)
	l.ApplyTransform([2][2]float64{
		{cos, -sin},
		{sin, cos},
	})
}

func (l *Line) ApplyTransform(T [2][2]float64) {
	l.Start.ApplyTransform(T)
	l.End.ApplyTransform(T)
}
