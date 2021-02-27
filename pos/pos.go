package pos

import (
	"fmt"
	"math"
)

type Pos struct {
	X, Y float64
}

func (p Pos) Unit() Pos {
	return p.Scale(1 / p.Length())
}

func (p Pos) Dir() Dir {
	return Dir{
		R:     p.Length(),
		Theta: math.Atan2(p.Y, p.X),
	}
}

func (p Pos) Length() float64 {
	return math.Sqrt(p.X*p.X + p.Y*p.Y)
}

func (p Pos) Scale(r float64) Pos {
	return Pos{p.X * r, p.Y * r}
}

func (p Pos) Near(q Pos, r float64) bool {
	return p.DistSquard(q) < r*r
}

func (p Pos) DistSquard(q Pos) float64 {
	dx := p.X - q.X
	dy := p.Y - q.Y
	return dx*dx + dy*dy
}

func (p Pos) Add(q Pos) Pos {
	return Pos{p.X + q.X, p.Y + q.Y}
}

func (p Pos) Sub(q Pos) Pos {
	return Pos{p.X - q.X, p.Y - q.Y}
}

func (p Pos) Move(d Dir) Pos {
	return p.Add(d.Pos())
}

func (p Pos) String() string {
	return fmt.Sprintf("[%0.5f,%0.5f]", p.X, p.Y)
}

type Dir struct {
	R, Theta float64
}

var North = Dir{R: 1.0, Theta: math.Pi / 2}

func (d Dir) Turn(theta float64) Dir {
	newDir := d
	newDir.Theta += theta
	return newDir.Normalise()
}

func (d Dir) TurnRight() Dir {
	return d.Turn(-math.Pi / 2)
}

func (d Dir) TurnLeft() Dir {
	return d.Turn(math.Pi / 2)
}

// To -math.Pi < theta <= math.Pi
func (d Dir) Normalise() Dir {
	theta := d.Theta
	if theta > math.Pi {
		theta -= 2 * math.Pi
	}
	theta = math.Mod(theta, 2*math.Pi)
	return Dir{d.R, theta}
}

func (d Dir) Scale(r float64) Dir {
	return Dir{R: d.R * r, Theta: d.Theta}
}

func (d Dir) Pos() Pos {
	return Pos{d.R * math.Cos(d.Theta), d.R * math.Sin(d.Theta)}
}

func (d Dir) String() string {
	return fmt.Sprintf("(%0.5f,%0.5f)", d.R, d.Theta)
}
