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
	Web() []pos.Pos
}
