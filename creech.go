package creech // import "github.com/jbert/creech"

import (
	"fmt"
	"log"
	"math"
	"math/rand"
	"sort"
	"time"

	"github.com/jbert/creech/render"

	// pos.Pos shorthand
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
		value := rand.Float64() * 10
		density := 0.4
		f := NewFood(value, density)
		f.SetRandomPos(g, f.Size())
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
			if c.pos.Near(p, c.Size()+size) {
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
	ticker := time.NewTicker(g.tickDur)
	defer ticker.Stop()
	ticks := 0

	for range ticker.C {
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

func (g *Game) Observe(r Region, excludeID int64) []Entity {
	var es []Entity
	for _, c := range g.creeches {
		if c.ID() != excludeID && r.Contains(c.Pos()) {
			es = append(es, c)
		}
	}
	for _, f := range g.food {
		if f.ID() != excludeID && r.Contains(f.Pos()) {
			es = append(es, f)
		}
	}
	return es
}

type Entity interface {
	ID() int64
	Pos() Pos
	Size() float64
}

type BaseEntity struct {
	id  int64
	pos Pos
}

func NewBaseEntity(p Pos) BaseEntity {
	return BaseEntity{
		pos: p,
		id:  rand.Int63(),
	}
}

func (be *BaseEntity) SetRandomPos(g *Game, size float64) {
	be.pos = g.randomEmptyPos(size)
}

func (be *BaseEntity) ID() int64 {
	return be.id
}

func (be *BaseEntity) Pos() Pos {
	return be.pos
}

type Creech struct {
	BaseEntity
	size float64

	name   string
	facing Polar

	plan func()
}

func NewCreech(name string, pos Pos) *Creech {
	creechSize := 1.0
	return &Creech{
		name:       name,
		size:       creechSize,
		facing:     North,
		BaseEntity: NewBaseEntity(pos),
	}
}

func (c *Creech) Size() float64 {
	return c.size
}

func (c *Creech) String() string {
	return fmt.Sprintf("%s: %s %s", c.name, c.pos, c.facing)
}

func (c *Creech) MakePlan(g *Game) {
	region := c.ViewRegion()
	entities := g.Observe(region, c.ID())
	sort.Slice(entities, func(i, j int) bool {
		return c.Pos().DistanceToSquared(entities[i].Pos()) <
			c.Pos().DistanceToSquared(entities[j].Pos())
	})
	c.plan = c.makeRandomPlan()
	for i, ei := range entities {
		switch e := ei.(type) {
		case *Food:
			log.Printf("%d ============================= EAT =================", i)
			c.plan = func() {
				c.TurnToward(e)
			}
			break
		case *Creech:
			log.Printf("%d ============================= FLEE =================", i)
			c.plan = func() {
				c.TurnAway(e)
				dist := c.maxMove() * rand.Float64()
				c.pos = c.pos.Move(c.facing.Scale(dist))
			}
			break
		default:
			panic(fmt.Sprintf("wtf: %T", ei))
		}
	}
}

func (c *Creech) DoPlan() {
	c.plan()
}

func (c *Creech) TurnAway(e Entity) {
	dTheta := turnHelper(c.facing, c.Pos(), e.Pos(), c.maxTurn(), false)
	c.facing = c.facing.Turn(dTheta)
}

func (c *Creech) TurnToward(e Entity) {
	dTheta := turnHelper(c.facing, c.Pos(), e.Pos(), c.maxTurn(), true)
	c.facing = c.facing.Turn(dTheta)
}

func turnHelper(facing Polar, p Pos, target Pos, maxTurn float64, towards bool) float64 {
	joiningLine := target.Sub(p).Polar()
	angleToTarget := joiningLine.Theta - facing.Theta
	if angleToTarget == 0 {
		if towards {
			return 0
		} else {
			return maxTurn
		}
	}

	isNegative := math.Signbit(angleToTarget)
	angleToTarget = math.Abs(angleToTarget)

	var dTheta float64
	if towards {
		// We want to reduce angleToTarget as much as possible
		// but not beyond zero
		dTheta = angleToTarget
		if dTheta > maxTurn {
			dTheta = maxTurn
		}
	} else {
		// We want to increase angleToTarget as much as possible
		// but not beyond pi
		dTheta = maxTurn
		if angleToTarget+dTheta > math.Pi {
			dTheta = math.Pi - angleToTarget
		}
		dTheta = -dTheta
	}

	// Restore sign
	if isNegative {
		dTheta = -dTheta
	}

	return dTheta
}

func (c *Creech) maxMove() float64 {
	return 0.5
}

func (c *Creech) maxTurn() float64 {
	return math.Pi * 0.125
}

func (c *Creech) viewDistance() float64 {
	return 7.0
}

func (c *Creech) viewSideDistance() float64 {
	return 4.0
}

func (c *Creech) makeRandomPlan() func() {
	return func() {
		r := rand.Intn(10)
		if r < 4 {
			turn := (rand.Float64() - 0.5) * c.maxTurn()
			c.facing = c.facing.Turn(turn)
		}
		dist := c.maxMove() + rand.Float64()
		c.pos = c.pos.Move(c.facing.Scale(dist))
	}
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
	i := int(c.pos.X)
	j := int(c.pos.Y)
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

func (c *Creech) ViewRegion() Region {
	viewDist := c.viewDistance()
	sideDist := c.viewSideDistance()

	frontSideStep := c.facing.Turn(math.Pi / 2).Scale(sideDist).Pos()
	backSideStep := frontSideStep.Scale(0.2)

	base := c.pos
	frontLeft := base.Add(c.facing.Scale(viewDist).Pos()).Add(frontSideStep)
	frontRight := frontLeft.Sub(frontSideStep.Scale(2.0))

	backLeft := base.Add(backSideStep)
	backRight := backLeft.Sub(backSideStep.Scale(2.0))

	pts := []Pos{
		backLeft,
		frontLeft,
		frontRight,
		backRight,
		backLeft,
	}
	r := NewRegion(pts)
	log.Printf("Viewegion: %+v\n", r)
	return r
}

func (c *Creech) Web() []render.DrawCommand {
	dir := c.facing.Pos().Scale(c.size)
	pts := arrow(c.pos, c.pos.Add(dir), 0.3)
	region := c.ViewRegion()
	viewPoly := render.Poly(region.ClosedPoints())
	viewPoly.DoFill = true
	colour := render.RGBA{0.5, 0.1, 0.1, 0.2}
	viewPoly.FillColour = colour
	viewPoly.LineColour = colour
	return []render.DrawCommand{
		render.Poly(pts),
		viewPoly,
	}
}

type Food struct {
	BaseEntity

	value   float64
	density float64
}

func NewFood(value float64, density float64) *Food {
	f := &Food{
		BaseEntity: NewBaseEntity(Pos{0, 0}),
		value:      value,
		density:    density,
	}
	return f
}

func (f *Food) Size() float64 {
	return f.value * f.density
}

func (f *Food) Screen() (int, int, byte) {
	var b byte
	if f.value < 3 {
		b = '.'
	} else if f.value < 6 {
		b = 'o'
	} else if f.value < 9 {
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

func (f *Food) Web() []render.DrawCommand {
	pts := closedPolygon(6, f.pos, f.Size()/2)
	return []render.DrawCommand{render.Poly(pts)}
}
