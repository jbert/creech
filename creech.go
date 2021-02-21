package creech // import "github.com/jbert/creech"

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/jbert/creech/render"
)

type Game struct {
	width, height int

	creeches []*Creech

	renderer render.Renderer
}

func NewGame(r render.Renderer) *Game {
	return &Game{
		width:    10,
		height:   10,
		renderer: r,
	}
}

func (g *Game) Init() error {
	origin := Pos{0, 0}
	bob := NewCreech("bob", origin)

	g.creeches = append(g.creeches, bob)
	return g.renderer.Init(g.width, g.height)
}

func (g *Game) Run() error {
	tickDur := time.Second
	tickCh := time.Tick(tickDur)
	ticks := 0

	for range tickCh {
		err := g.Draw(ticks)
		if err != nil {
			return fmt.Errorf("Can't Draw: %w", err)
		}
		ticks++
		g.Update()
	}
	return nil
}

func (g *Game) Update() {
	for _, creech := range g.creeches {
		creech.MakePlan()
	}
	for _, creech := range g.creeches {
		creech.DoPlan()
	}
}

func (g *Game) Draw(ticks int) error {
	r := g.renderer
	err := r.StartFrame()
	if err != nil {
		return fmt.Errorf("StartFrame: %w", err)
	}
	for _, creech := range g.creeches {
		p := creech.Pos()
		err = r.DrawAt(p.X, p.Y, creech)
	}
	if err != nil {
		return fmt.Errorf("DrawAt: %w", err)
	}
	err = r.FinishFrame()
	if err != nil {
		return fmt.Errorf("FinishFrame: %w", err)
	}

	fmt.Printf("%d ticks\n", ticks)

	return nil
}

type Creech struct {
	name   string
	pos    Pos
	facing Dir
}

func NewCreech(name string, pos Pos) *Creech {
	return &Creech{
		name:   name,
		pos:    pos,
		facing: East,
	}
}

func (c *Creech) Pos() Pos {
	return c.pos
}

func (c *Creech) MakePlan() {
}

func (c *Creech) DoPlan() {
	r := rand.Intn(10)
	if r < 2 {
		c.facing = c.facing.TurnLeft()
	} else if r < 4 {
		c.facing = c.facing.TurnRight()
	}
	c.pos = c.pos.Move(c.facing)
}

func (c *Creech) Screen() byte {
	return 'X'
}
