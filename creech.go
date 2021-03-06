package creech // import "github.com/jbert/creech"

import (
	"fmt"
	"math"
	"math/rand"
	"sort"
	"strings"
	"time"

	"github.com/jbert/creech/render"

	// pos.Pos shorthand
	. "github.com/jbert/creech/pos"
)

type Game struct {
	worldSize Pos
	tickDur   time.Duration
	renderer  render.Renderer

	state State
}

type State struct {
	creeches []*Creech
	food     []*Food
}

func (s *State) String() string {
	var lines []string
	for _, c := range s.creeches {
		lines = append(lines, c.String())
	}
	for _, f := range s.food {
		lines = append(lines, f.String())
	}
	return strings.Join(lines, "\n")
}

func (s *State) AddCreeches() {
	bob := NewCreech("bob", Pos{0, 0})
	s.creeches = append(s.creeches, bob)

	alice := NewCreech("alice", Pos{2, 2})
	s.creeches = append(s.creeches, alice)
}

func (s *State) AddFood(worldSize Pos) {
	numFood := 5
	for i := 0; i < numFood; i++ {
		value := rand.Float64() * 10
		f := NewFood(value)
		f.SetRandomPos(s, worldSize, f.Size())
		s.food = append(s.food, f)
	}
}

// TODO: at a certain point, we'll want to avoid looping over everything to do this
func (s *State) randomEmptyPos(worldSize Pos, size float64) Pos {
RANDOM_POSITION:
	for {
		p := Pos{rand.Float64() * worldSize.X, rand.Float64() * worldSize.Y}
		p = moduloPos(p, worldSize)
		for _, c := range s.creeches {
			if c.pos.Near(p, c.Size()+size) {
				continue RANDOM_POSITION
			}
		}
		for _, f := range s.food {
			if f.Pos().Near(p, f.Size()+size) {
				continue RANDOM_POSITION
			}
		}
		return p
	}
}

func (s *State) Draw(r render.Renderer, ticks int) error {
	dlog := func(s string, args ...interface{}) {
		//		log.Printf(s, args...)
	}
	dlog("StartFrame")
	err := r.StartFrame()
	if err != nil {
		return fmt.Errorf("StartFrame: %w", err)
	}

	for i, f := range s.food {
		dlog("Draw Food %d", i)
		err = r.Draw(f)
		if err != nil {
			return fmt.Errorf("Draw: %w", err)
		}
	}

	for i, creech := range s.creeches {
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

	fmt.Printf("%d ticks\n%s\n\n", ticks, s)
	//	for _, creech := range s.creeches {
	//		fmt.Printf("%s\n", creech)
	//	}

	return nil
}

func NewGame(r render.Renderer, tickDur time.Duration) *Game {
	return &Game{
		worldSize: Pos{40, 40},
		tickDur:   tickDur,
		renderer:  r,
	}
}

func (g *Game) Init() error {
	g.state.AddCreeches()
	g.state.AddFood(g.worldSize)
	return g.renderer.Init(g.worldSize.X, g.worldSize.Y)
}

func (g *Game) Run() error {
	ticker := time.NewTicker(g.tickDur)
	defer ticker.Stop()
	ticks := 0

	for range ticker.C {
		err := g.state.Draw(g.renderer, ticks)
		if err != nil {
			return fmt.Errorf("Can't Draw: %w", err)
		}
		ticks++
		g.Update()
	}
	return nil
}

func (g *Game) Update() {
	for _, creech := range g.state.creeches {
		if !creech.Dead() {
			creech.ModuloPos(g.worldSize)
			creech.MakePlan(g)
		}
	}
	for _, creech := range g.state.creeches {
		creech.DoPlan()
	}
}

func (g *Game) Observe(r Region, excludeID int64) []Entity {
	var es []Entity
	for _, c := range g.state.creeches {
		if c.ID() != excludeID && r.Contains(c.Pos()) {
			es = append(es, c)
		}
	}
	for _, f := range g.state.food {
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

func (be *BaseEntity) SetRandomPos(s *State, worldSize Pos, size float64) {
	be.pos = s.randomEmptyPos(worldSize, size)
}

func (be *BaseEntity) ID() int64 {
	return be.id
}

func (be *BaseEntity) Pos() Pos {
	return be.pos
}

type Plan struct {
	name   string // for display
	action func()
	cost   float64
}

func NewPlan(name string, action func()) *Plan {
	return &Plan{name: name, action: action, cost: 0.1}
}

func (p *Plan) Cost() float64 {
	return p.cost
}

func (p *Plan) Execute() {
	p.action()
}

type Creech struct {
	BaseEntity
	size float64

	name   string
	facing Polar

	food float64
	plan *Plan
}

func NewCreech(name string, pos Pos) *Creech {
	creechSize := 1.0
	c := &Creech{
		name:       name,
		size:       creechSize,
		facing:     North,
		BaseEntity: NewBaseEntity(pos),
	}
	c.food = c.maxFood() / 2
	return c
}

func (c *Creech) Dead() bool {
	return c.food <= 0
}

func (c *Creech) Size() float64 {
	return c.size
}

func (c *Creech) String() string {
	return fmt.Sprintf("%s: %5.2f %s %s", c.name, c.food, c.pos, c.facing)
}

func (c *Creech) Full() bool {
	return c.food >= c.maxFood()
}

func (c *Creech) Eat(f *Food) {
	biteSize := c.biteSize()
	if biteSize > c.maxFood()-c.food {
		biteSize = c.maxFood() - c.food
	}
	f.Consume(biteSize)
	c.food += biteSize
}

func (c *Creech) MakePlan(g *Game) {
	region := c.ViewRegion()
	entities := g.Observe(region, c.ID())
	sort.Slice(entities, func(i, j int) bool {
		return c.Pos().DistanceToSquared(entities[i].Pos()) <
			c.Pos().DistanceToSquared(entities[j].Pos())
	})
	c.plan = c.makeRandomPlan()
	for _, ei := range entities {
		switch e := ei.(type) {
		case *Food:
			if c.Full() {
				continue
			}
			c.plan = NewPlan("FOOD", func() {
				c.TurnToward(e)

				eatDistance := c.eatDistance(e)
				if c.Pos().DistanceTo(e.Pos()) < eatDistance {
					c.Eat(e)
				} else {
					c.ApproachTo(e, eatDistance)
				}
			})
			break
		case *Creech:
			c.plan = NewPlan("FLEE", func() {
				c.TurnAway(e)
				dist := c.maxMove() * (0.5 + 0.5*rand.Float64())
				c.pos = c.pos.Move(c.facing.Scale(dist))
			})
			break
		default:
			panic(fmt.Sprintf("wtf: %T", ei))
		}
	}
}

func (c *Creech) DoPlan() {
	if c.plan == nil {
		return
	}
	c.plan.Execute()
	c.food -= c.plan.cost
	c.plan = nil
}

func (c *Creech) ApproachTo(e Entity, d float64) {
	p := c.Pos().PolarTo(e.Pos())
	if math.Abs(p.Theta) > math.Pi/2 {
		return
	}

	moveDist := c.maxMove()
	if moveDist > (p.R + d) {
		moveDist = p.R + d
	}
	c.MoveForward(moveDist)
}

func (c *Creech) TurnAway(e Entity) {
	dTheta := turnHelper(c.facing, c.Pos(), e.Pos(), c.maxTurn(), false)
	c.facing = c.facing.Turn(dTheta)
}

func (c *Creech) TurnToward(e Entity) {
	dTheta := turnHelper(c.facing, c.Pos(), e.Pos(), c.maxTurn(), true)
	c.facing = c.facing.Turn(dTheta)
}

func (c *Creech) MoveForward(d float64) {
	c.pos = c.pos.Move(c.facing.Scale(d))
}

func turnHelper(facing Polar, p Pos, target Pos, maxTurn float64, towards bool) float64 {
	joiningLine := p.PolarTo(target)
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

func (c *Creech) biteSize() float64 {
	return 1.5
}

func (c *Creech) maxFood() float64 {
	return 10
}

func (c *Creech) maxMove() float64 {
	return 0.5
}

func (c *Creech) maxTurn() float64 {
	return math.Pi * 0.125
}

func (c *Creech) eatDistance(f *Food) float64 {
	return f.Size() + 1
}

func (c *Creech) viewDistance() float64 {
	return 10.0
}

func (c *Creech) viewSideDistance() float64 {
	return 4.0
}

func (c *Creech) makeRandomPlan() *Plan {
	return NewPlan("RANDOM", func() {
		r := rand.Intn(10)
		if r < 4 {
			turn := (rand.Float64() - 0.5) * c.maxTurn()
			c.facing = c.facing.Turn(turn)
		}
		dist := c.maxMove() * rand.Float64()
		c.MoveForward(dist)
	})
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
	if c.Dead() {
		b = 'X'
	} else if math.Abs(t) < math.Pi/4 {
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
	return r
}

func (c *Creech) Web() []render.DrawCommand {
	if c.Dead() {
		step := c.facing.Scale(c.size).Pos()
		sideStep := c.facing.Turn(math.Pi / 2).Scale(c.size).Pos()
		p := c.Pos()
		return []render.DrawCommand{
			render.Poly([]Pos{p.Sub(step), p.Add(step)}),
			render.Poly([]Pos{p.Sub(sideStep), p.Add(sideStep)}),
		}
	}
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

	value float64
}

func NewFood(value float64) *Food {
	f := &Food{
		BaseEntity: NewBaseEntity(Pos{0, 0}),
		value:      value,
	}
	return f
}

func (f *Food) String() string {
	return fmt.Sprintf("%5.2f: %s", f.value, f.Pos())
}

func (f *Food) Consume(bite float64) {
	f.value -= bite
	//	fmt.Printf("Ate %f - %f left\n", bite, f.value)
}

func (f *Food) Size() float64 {
	return f.value
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
