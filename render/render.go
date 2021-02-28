package render

import (
	"fmt"

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

// CSS format for colour
func (rgba RGBA) MarshalJSON() ([]byte, error) {
	f2i := func(f float64) int { return int(255 * f) }
	s := fmt.Sprintf(`"rgba(%d,%d,%d,%f)"`, f2i(rgba.R), f2i(rgba.G), f2i(rgba.B), rgba.A)
	return []byte(s), nil
}

type DrawCommand struct {
	What       DrawType
	Points     []pos.Pos
	LineColour RGBA
	DoFill     bool
	FillColour RGBA
}

var Black = RGBA{0, 0, 0, 1}
var White = RGBA{1, 1, 1, 1}

func Poly(pts []pos.Pos) DrawCommand {
	return DrawCommand{
		What:       DrawPoly,
		Points:     pts,
		LineColour: Black,
		DoFill:     false,
		FillColour: White,
	}
}
