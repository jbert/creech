package render

type Renderer interface {
	Init(w, h int) error
	StartFrame() error
	FinishFrame() error
	DrawAt(i, j int, d Drawable) error
}

type Drawable interface {
	Screen() byte
}
