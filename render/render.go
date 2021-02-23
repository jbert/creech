package render

type Renderer interface {
	Init(w, h float64) error
	StartFrame() error
	FinishFrame() error
	DrawAt(x, y float64, d Drawable) error
}

type Drawable interface {
	Screen() byte
}
