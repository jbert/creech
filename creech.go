package creech // import "github.com/jbert/creech"

import (
	"fmt"
	"time"

	"github.com/jbert/creech/render"
)

type Game struct {
	width, height int
	renderer      render.Renderer
}

func NewGame(r render.Renderer) *Game {
	return &Game{
		width:    10,
		height:   10,
		renderer: r,
	}
}

func (g *Game) Init() error {
	return g.renderer.Init(g.width, g.height)
}

func (g *Game) Run() error {
	tickDur := time.Second
	tickCh := time.Tick(tickDur)

	for range tickCh {
		err := g.Draw()
		if err != nil {
			return fmt.Errorf("Can't Draw: %w", err)
		}
	}
	return nil
}

func (g *Game) Draw() error {
	x := NewCreech("bob")

	r := g.renderer
	err := r.StartFrame()
	if err != nil {
		return fmt.Errorf("StartFrame: %w", err)
	}
	err = r.DrawAt(3, 4, x)
	if err != nil {
		return fmt.Errorf("DrawAt: %w", err)
	}
	err = r.FinishFrame()
	if err != nil {
		return fmt.Errorf("FinishFrame: %w", err)
	}

	return nil
}

type Creech struct {
	name string
}

func NewCreech(name string) *Creech {
	return &Creech{name: name}
}

func (c *Creech) Screen() byte {
	return 'X'
}
