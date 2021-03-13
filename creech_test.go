package creech

import (
	"math"
	"testing"

	. "github.com/jbert/creech/pos"
)

func TestTurnHelper(t *testing.T) {
	pi2 := math.Pi / 2
	//	pi4 := math.Pi / 4
	o := Pos{0, 0}
	mt := math.Pi / 10
	east := 0.0
	north := pi2
	south := -pi2
	west := math.Pi

	a := Pos{1, 0}
	b := Pos{1, 1}

	bigX := 100.0
	smallAng := math.Atan2(1, bigX)
	t.Logf("smallAng %f", smallAng)

	testCases := []struct {
		facingTheta float64
		p           Pos
		target      Pos
		maxTurn     float64
		towards     bool
		expected    float64
	}{
		// Copy failing test first
		{west, o, Pos{bigX, 1}, mt, false, +smallAng},
		// -----

		{east, o, a, mt, true, 0},
		{north, o, a, mt, true, -mt},
		{south, o, a, mt, true, +mt},

		{east, o, b, mt, true, +mt},
		{north, o, b, mt, true, -mt},
		{south, o, b, mt, true, +mt},

		{east, o, Pos{bigX, 1}, mt, true, +smallAng},

		{east, o, a, mt, false, +mt},
		{north, o, a, mt, false, +mt},
		{south, o, a, mt, false, -mt},

		{east, o, b, mt, false, -mt},
		{north, o, b, mt, false, +mt},
		{south, o, b, mt, false, -mt},

		{west, o, Pos{bigX, 1}, mt, false, +smallAng},
	}

	for _, tc := range testCases {
		t.Logf("%+v", tc)
		got := turnHelper(Polar{R: 0, Theta: tc.facingTheta}, tc.p, tc.target, tc.maxTurn, tc.towards)
		if !approxEqual(got, tc.expected) {
			t.Fatalf("got %f expected %f", got, tc.expected)
		}
	}
}

func approxEqual(a, b float64) bool {
	eps := 1e-8
	return math.Abs(a-b) < eps
}
