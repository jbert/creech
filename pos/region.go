package pos

import (
	"errors"
	"fmt"
	"math"
	"strings"
)

type Region struct {
	points []Pos
}

func NewRegion(pts []Pos) Region {
	return Region{
		points: pts,
	}
}

func (r Region) String() string {
	s := make([]string, len(r.points))
	for i, p := range r.points {
		s[i] = p.String()
	}
	return strings.Join(s, ",")
}

func (r Region) ClosedPoints() []Pos {
	pts := make([]Pos, len(r.points))
	copy(pts, r.points)
	pts = append(pts, r.points[0])
	return pts
}

type LineSegment struct {
	From, To Pos

	isVertical     bool
	gradient, yInt float64
}

func NewLineSegment(from, to Pos) LineSegment {
	l := LineSegment{
		From: from,
		To:   to,
	}

	dY := l.To.Y - l.From.Y
	dX := l.To.X - l.From.X
	if approxEqual(dX, 0) {
		l.isVertical = true
	}
	l.gradient = dY / dX
	l.yInt = l.From.Y - l.gradient*l.From.X

	return l
}

func (ls LineSegment) String() string {
	return fmt.Sprintf("%s -> %s", ls.From, ls.To)
}

var errVertical = errors.New("Line is vertical")

func (ls *LineSegment) GradAndIntercept() (float64, float64, error) {
	if ls.isVertical {
		return 0, 0, errVertical
	}
	return ls.gradient, ls.yInt, nil
}

// Evaluate segment extended to line at X
func (ls LineSegment) Eval(x float64) (float64, error) {
	m, c, err := ls.GradAndIntercept()
	if err != nil {
		return 0, err
	}
	return m*x + c, nil
}

func (ls LineSegment) LineIntersects(ls2 LineSegment) (Pos, bool) {
	m, c, err := ls.GradAndIntercept()
	m2, c2, err2 := ls2.GradAndIntercept()
	if err != nil && err == err2 {
		// Both vertical
		if ls.From.X == ls2.From.X {
			return Pos{ls.From.X, 0}, true
		}
		return Pos{}, false
	}

	if err2 != nil {
		// Switch ls and ls2
		return ls2.LineIntersects(ls)
	}

	if err != nil {
		// Handle ls vertical, ls2 not
		x := ls.From.X
		var y float64
		y, err = ls2.Eval(x)
		if err != nil {
			panic("Expected ls2 to be non-vertical")
		}
		return Pos{x, y}, true
	}

	// Neither vertical

	if m == m2 {
		// Parallel
		if c != c2 {
			// not co-incident
			return Pos{}, false
		}
		// Co-incident
		return Pos{0, c}, true
	}

	x := (c2 - c) / (m - m2)
	y, err := ls.Eval(x)
	if err != nil {
		panic("Expected ls2 to be non-vertial")
	}
	return Pos{x, y}, true
}

func (ls LineSegment) BoundingRectContains(p Pos) bool {
	minX := math.Min(ls.From.X, ls.To.X)
	maxX := math.Max(ls.From.X, ls.To.X)
	minY := math.Min(ls.From.Y, ls.To.Y)
	maxY := math.Max(ls.From.Y, ls.To.Y)
	return minX <= p.X &&
		maxX >= p.X &&
		minY <= p.Y &&
		maxY >= p.Y
}

func (ls LineSegment) ContainsSeg(ls2 LineSegment) bool {
	return ls.ContainsPos(ls2.From) && ls.ContainsPos(ls2.To)
}

func (ls LineSegment) ContainsPos(p Pos) bool {
	if !ls.BoundingRectContains(p) {
		return false
	}
	if ls.isVertical {
		return approxEqual(p.X, ls.From.X)
	}
	y, err := ls.Eval(p.X)
	if err != nil {
		panic("ls should not be vertical")
	}
	return approxEqual(y, p.Y)
}

func (ls LineSegment) SegmentIntersectsLine(ray LineSegment) bool {
	p, intersect := ls.LineIntersects(ray)
	if !intersect {
		return false
	}
	ok1 := ls.BoundingRectContains(p)
	ok2 := ray.BoundingRectContains(p)
	return ok1 && ok2
}

func (r Region) Translate(v Pos) Region {
	pts := make([]Pos, len(r.points))
	for i := range r.points {
		pts[i] = r.points[i].Add(v)
	}
	return NewRegion(pts)
}

func (r Region) Overlaps(s Region) bool {
	for _, p := range s.ClosedPoints() {
		if r.Contains(p) {
			return true
		}
	}
	return false
}

func (r Region) Contains(q Pos) bool {
	// Take a line segment "to infinity"
	big := 1e7
	qRay := NewLineSegment(q, Pos{0, big})

	var last Pos
	numIntersects := 0
	for i, p := range r.ClosedPoints() {
		if i > 0 {
			seg := NewLineSegment(last, p)
			if seg.ContainsPos(q) {
				return true
			}
			// If collinear with one side, but the pt is not within the side
			// Do not count this tangential ray
			if !qRay.ContainsSeg(seg) {
				if seg.SegmentIntersectsLine(qRay) {
					numIntersects++
				}
			}
		}
		last = p
	}
	return numIntersects%2 == 1
}
