package creech // import "github.com/jbert/creech"

import (
	"fmt"
	"math"
	"math/rand"
	"time"

	"github.com/jbert/creech/render"
)

type Game struct {
	size Pos

	creeches []*Creech
	food     []*Food

	renderer render.Renderer
}

func NewGame(r render.Renderer) *Game {
	return &Game{
		size:     Pos{20, 20},
		renderer: r,
	}
}

func (g *Game) Init() error {

	g.AddCreeches()
	g.AddFood()
	return g.renderer.Init(int(g.size.X), int(g.size.Y))
}

func (g *Game) AddCreeches() {
	bob := NewCreech("bob", Pos{0, 0})
	g.creeches = append(g.creeches, bob)

	alice := NewCreech("alice", Pos{2, 2})
	g.creeches = append(g.creeches, alice)
}

func (g *Game) AddFood() {
	numFood := 5
	for i := 0; i < numFood; i++ {
		value := rand.Intn(10)
		foodSize := 1.0
		f := NewFood(value, g.randomEmptyPos(foodSize), foodSize)
		g.food = append(g.food, f)
	}
}

// TODO: at a certain point, we'll want to avoid looping over everything to do this
func (g *Game) randomEmptyPos(size float64) Pos {
RANDOM_POSITION:
	for {
		p := Pos{rand.Float64() * g.size.X, rand.Float64() * g.size.Y}
		p = moduloPos(p, g.size)
		for _, c := range g.creeches {
			if c.Pos().Near(p, c.Size()+size) {
				continue RANDOM_POSITION
			}
		}
		for _, f := range g.food {
			if f.Pos().Near(p, f.Size()+size) {
				continue RANDOM_POSITION
			}
		}
		return p
	}
}

func (g *Game) Run() error {
	tickDur := 100 * time.Millisecond
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
		creech.ModuloPos(g.size)
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

	for _, f := range g.food {
		p := f.Pos()
		err = r.DrawAt(p.X, p.Y, f)
		if err != nil {
			return fmt.Errorf("DrawAt: %w", err)
		}
	}

	for _, creech := range g.creeches {
		p := creech.Pos()
		err = r.DrawAt(p.X, p.Y, creech)
		if err != nil {
			return fmt.Errorf("DrawAt: %w", err)
		}
	}

	err = r.FinishFrame()
	if err != nil {
		return fmt.Errorf("FinishFrame: %w", err)
	}

	fmt.Printf("%d ticks\n", ticks)
	for _, creech := range g.creeches {
		fmt.Printf("%s\n", creech)
	}

	return nil
}

type Entity struct {
	pos  Pos
	size float64
}

type Creech struct {
	Entity

	name   string
	facing Dir
}

func NewCreech(name string, pos Pos) *Creech {
	creechSize := 1.0
	return &Creech{
		name:   name,
		facing: North,
		Entity: Entity{
			pos:  pos,
			size: creechSize,
		},
	}
}

func (c *Creech) String() string {
	return fmt.Sprintf("%s: %s %s", c.name, c.pos, c.facing)
}

func (c *Creech) Pos() Pos {
	return c.pos
}

func (c *Creech) Size() float64 {
	return c.size
}

func (c *Creech) MakePlan() {
}

func (c *Creech) ModuloPos(worldSize Pos) {
	c.pos = moduloPos(c.pos, worldSize)
}

func moduloPos(p Pos, worldSize Pos) Pos {
	q := p // Struct copy
	if q.X > worldSize.X/2 {
		q.X -= worldSize.X
	}
	if q.X <= -worldSize.X/2 {
		q.X += worldSize.X
	}
	if q.Y > worldSize.Y/2 {
		q.Y -= worldSize.Y
	}
	if q.Y <= -worldSize.Y/2 {
		q.Y += worldSize.Y
	}
	return q
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
	t := c.facing.Theta
	if math.Abs(t) < math.Pi/4 {
		return '>'
	} else if math.Pi/4 < t && t < 3*math.Pi/4 {
		return '^'
	} else if -math.Pi/4 > t && t > -3*math.Pi/4 {
		return 'v'
	} else if math.Abs(t) > 3*math.Pi/4 {
		return '<'
	}
	panic("wtf")
}

type Food struct {
	Entity

	value int
}

func NewFood(v int, p Pos, size float64) *Food {
	foodSize := 1.0
	return &Food{
		Entity: Entity{
			pos:  p,
			size: foodSize,
		},
		value: v,
	}
}

func (f *Food) Size() float64 {
	return f.size
}

func (f *Food) Pos() Pos {
	return f.pos
}

func (f *Food) Screen() byte {
	if f.value < 3 {
		return '.'
	} else if f.value < 3 {
		return 'o'
	} else if f.value < 6 {
		return 'O'
	} else {
		return '*'
	}
}
