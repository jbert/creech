package pos

import (
	"fmt"
	"math"
)

func approxEqual(a, b float64) bool {
	eps := 1e-5
	return math.Abs(a-b) < eps
}

type Pos struct {
	X, Y float64
}

func (p Pos) Equals(q Pos) bool {
	return approxEqual(p.X, q.X) && approxEqual(p.Y, q.Y)
}

func (p Pos) Unit() Pos {
	return p.Scale(1 / p.Length())
}

func (p Pos) Polar() Polar {
	return Polar{
		R:     p.Length(),
		Theta: math.Atan2(p.Y, p.X),
	}
}

func (p Pos) Length() float64 {
	return math.Sqrt(p.X*p.X + p.Y*p.Y)
}

func (p Pos) PolarTo(q Pos) Polar {
	return q.Sub(p).Polar()
}

func (p Pos) DistanceTo(q Pos) float64 {
	return math.Sqrt(p.DistanceToSquared(q))
}

func (p Pos) DistanceToSquared(q Pos) float64 {
	dx := p.X - q.X
	dy := p.Y - q.Y
	return dx*dx + dy*dy
}

func (p Pos) Scale(r float64) Pos {
	return Pos{p.X * r, p.Y * r}
}

func (p Pos) Near(q Pos, r float64) bool {
	return p.DistanceToSquared(q) < r*r
}

func (p Pos) Add(q Pos) Pos {
	return Pos{p.X + q.X, p.Y + q.Y}
}

func (p Pos) Sub(q Pos) Pos {
	return Pos{p.X - q.X, p.Y - q.Y}
}

func (p Pos) Move(q Polar) Pos {
	return p.Add(q.Pos())
}

func (p Pos) String() string {
	return fmt.Sprintf("[%0.5f,%0.5f]", p.X, p.Y)
}

type Polar struct {
	R, Theta float64
}

var North = Polar{R: 1.0, Theta: math.Pi / 2}

func (p Polar) Turn(theta float64) Polar {
	newP := p
	newP.Theta += theta
	return newP.Normalise()
}

func (p Polar) TurnRight() Polar {
	return p.Turn(-math.Pi / 2)
}

func (p Polar) TurnLeft() Polar {
	return p.Turn(math.Pi / 2)
}

// To -math.Pi < theta <= math.Pi
func (p Polar) Normalise() Polar {
	theta := p.Theta
	if theta > math.Pi {
		theta -= 2 * math.Pi
	}
	theta = math.Mod(theta, 2*math.Pi)
	return Polar{p.R, theta}
}

func (p Polar) Scale(r float64) Polar {
	return Polar{R: p.R * r, Theta: p.Theta}
}

func (p Polar) Pos() Pos {
	return Pos{p.R * math.Cos(p.Theta), p.R * math.Sin(p.Theta)}
}

func (p Polar) String() string {
	return fmt.Sprintf("(%0.5f,%0.5f)", p.R, p.Theta)
}
