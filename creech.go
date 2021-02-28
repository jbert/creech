package creech // import "github.com/jbert/creech"

import (
	"fmt"
	"math"
	"math/rand"
	"time"

	"github.com/jbert/creech/render"

	. "github.com/jbert/creech/pos"
)

type Game struct {
	size     Pos
	tickDur  time.Duration
	renderer render.Renderer

	creeches []*Creech
	food     []*Food
}

func NewGame(r render.Renderer, tickDur time.Duration) *Game {
	return &Game{
		size:     Pos{40, 40},
		tickDur:  tickDur,
		renderer: r,
	}
}

func (g *Game) Init() error {

	g.AddCreeches()
	g.AddFood()
	return g.renderer.Init(g.size.X, g.size.Y)
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
	tickCh := time.Tick(g.tickDur)
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
		creech.MakePlan(g)
	}
	for _, creech := range g.creeches {
		creech.DoPlan()
	}
}

func (g *Game) Draw(ticks int) error {
	r := g.renderer
	dlog := func(s string, args ...interface{}) {
		//		log.Printf(s, args...)
	}
	dlog("StartFrame")
	err := r.StartFrame()
	if err != nil {
		return fmt.Errorf("StartFrame: %w", err)
	}

	for i, f := range g.food {
		dlog("Draw Food %d", i)
		err = r.Draw(f)
		if err != nil {
			return fmt.Errorf("Draw: %w", err)
		}
	}

	for i, creech := range g.creeches {
		dlog("Draw Creech %d", i)
		err = r.Draw(creech)
		if err != nil {
			return fmt.Errorf("Draw: %w", err)
		}
	}

	dlog("FinishFrame")
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

func (e *Entity) Pos() Pos {
	return e.pos
}

func (e *Entity) Size() float64 {
	return e.size
}

type Creech struct {
	Entity

	name   string
	facing Polar
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

func (c *Creech) MakePlan(g *Game) {
	//	region := c.ViewRegion()
	//	itemsOfInterest := g.Observe(region)
}

func (c *Creech) DoPlan() {
	r := rand.Intn(10)
	if r < 4 {
		turn := (rand.Float64() - 0.5) * 0.5 * math.Pi
		c.facing = c.facing.Turn(turn)
	}
	dist := 0.5 + rand.Float64()
	c.pos = c.pos.Move(c.facing.Scale(dist))
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

func (c *Creech) Screen() (int, int, byte) {
	t := c.facing.Theta
	i := int(c.Pos().X)
	j := int(c.Pos().Y)
	var b byte
	if math.Abs(t) < math.Pi/4 {
		b = '>'
	} else if math.Pi/4 < t && t < 3*math.Pi/4 {
		b = '^'
	} else if -math.Pi/4 > t && t > -3*math.Pi/4 {
		b = 'v'
	} else if math.Abs(t) > 3*math.Pi/4 {
		b = '<'
	} else {
		panic("wtf")
	}
	return i, j, b
}

func arrow(from, to Pos, headSize float64) []Pos {
	p := to.Sub(from).Polar()
	p.R = headSize
	back1 := p.Turn(math.Pi * 3 / 4)
	back2 := p.Turn(-math.Pi * 3 / 4)
	return []Pos{
		from,
		to,
		to.Add(back1.Pos()),
		to.Add(back2.Pos()),
		to,
	}
}
func (c *Creech) Web() []Pos {
	dir := c.facing.Pos().Scale(c.size)
	return arrow(c.pos, c.pos.Add(dir), 0.3)
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

func (f *Food) Screen() (int, int, byte) {
	var b byte
	if f.value < 3 {
		b = '.'
	} else if f.value < 3 {
		b = 'o'
	} else if f.value < 6 {
		b = 'O'
	} else {
		b = '*'
	}
	i := int(f.Pos().X)
	j := int(f.Pos().Y)
	return i, j, b
}

func closedPolygon(sides int, p Pos, r float64) []Pos {
	pts := make([]Pos, sides+1)
	theta := 0.0
	dTheta := 2 * math.Pi / float64(sides)
	for i := 0; i <= sides; i++ {
		pts[i] = Pos{X: p.X + r*math.Cos(theta), Y: p.Y + r*math.Sin(theta)}
		theta += dTheta
	}
	return pts
}

func (f *Food) Web() []Pos {
	return closedPolygon(6, f.pos, f.size/2)
}
