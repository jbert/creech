package render

import (
	"github.com/jbert/creech/pos"
)

type Renderer interface {
	Init(w, h float64) error
	StartFrame() error
	FinishFrame() error
	Draw(d Drawable) error
}

type Drawable interface {
	Screen() (int, int, byte)
	Web() []DrawCommand
}

type DrawType int

const (
	StartFrame DrawType = iota
	DrawPoly
	FinishFrame
)

type RGBA struct {
	R, G, B, A float64
}

type DrawCommand struct {
	What       DrawType
	Points     []pos.Pos
	LineColour RGBA
	DoFill     bool
	FillColour RGBA
}

var Black = RGBA{0x00, 0x00, 0x00, 0x00}
var White = RGBA{0xff, 0xff, 0xff, 0x00}

func Poly(pts []pos.Pos) DrawCommand {
	return DrawCommand{
		What:       DrawPoly,
		Points:     pts,
		LineColour: Black,
		DoFill:     false,
		FillColour: White,
	}
}
