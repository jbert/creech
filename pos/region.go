package pos

type Region struct {
	points []Pos
}

func NewRegion(pts []Pos) Region {
	return Region{
		points: pts,
	}
}

func (r Region) ClosedPoints() []Pos {
	pts := make([]Pos, len(r.points))
	copy(pts, r.points)
	pts = append(pts, r.points[0])
	return pts
}
