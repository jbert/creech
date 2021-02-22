package creech

import "fmt"

type Pos struct {
	X, Y int
}

func (p Pos) Equal(q Pos) bool {
	return p.X == q.X && p.Y == q.Y
}

func (p Pos) Add(q Pos) Pos {
	return Pos{p.X + q.X, p.Y + q.Y}
}

func (p Pos) Move(d Dir) Pos {
	return p.Add(Pos(d))
}

func (p Pos) String() string {
	return fmt.Sprintf("[%d,%d]", p.X, p.Y)
}

type Dir Pos

var North = Dir{0, 1}
var East = Dir{1, 0}
var South = Dir{0, -1}
var West = Dir{-1, 0}

func (d Dir) String() string {
	switch d {
	case North:
		return "N"
	case East:
		return "E"
	case South:
		return "S"
	case West:
		return "W"
	default:
		panic(fmt.Sprintf("wtf: %v", d))
	}
}

func (d Dir) TurnRight() Dir {
	switch d {
	case North:
		return East
	case East:
		return South
	case South:
		return West
	case West:
		return North
	default:
		panic(fmt.Sprintf("wtf: %v", d))
	}
}

func (d Dir) TurnLeft() Dir {
	switch d {
	case North:
		return West
	case East:
		return North
	case South:
		return East
	case West:
		return South
	default:
		panic(fmt.Sprintf("wtf: %v", d))
	}
}
