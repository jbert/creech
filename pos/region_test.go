package pos

import "testing"

func TestSegmentIntersectLine(t *testing.T) {
	testCases := []struct {
		from                   Pos
		to                     Pos
		lineFrom               Pos
		lineTo                 Pos
		expectSegmentIntersect bool
		expectLineIntersect    bool
	}{
		{Pos{0, 0}, Pos{1, 1}, Pos{0, 1}, Pos{1, 0}, true, true},
	}

	for _, tc := range testCases {
		t.Logf("%+v", tc)
		ls := NewLineSegment(tc.from, tc.to)
		l := NewLineSegment(tc.lineFrom, tc.lineTo)

		p, got := ls.LineIntersects(l)
		if tc.expectLineIntersect != got {
			t.Fatalf("Failed LineIntersects : got %v expected %v", got, tc.expectLineIntersect)
		}
		t.Logf("%s and %s intersect at %s is %v", ls, l, p, got)

		got = ls.SegmentIntersectsLine(l)
		if tc.expectSegmentIntersect != got {
			t.Fatalf("Failed SegmentIntersectsLine : got %v expected %v", got, tc.expectSegmentIntersect)
		}
	}
}

func TestBoundingRectContains(t *testing.T) {
	testCases := []struct {
		from     Pos
		to       Pos
		p        Pos
		expected bool
	}{
		{Pos{0, 0}, Pos{1, 1}, Pos{0.5, 0.5}, true},
		{Pos{0, 0}, Pos{1, 1}, Pos{1, 1}, true},
		{Pos{0, 0}, Pos{1, 1}, Pos{2, 2}, false},
		{Pos{0, 0}, Pos{1, 1}, Pos{-0.1, -0.1}, false},

		{Pos{0, 1}, Pos{0, 0}, Pos{0, 0.4}, true},
	}

	for _, tc := range testCases {
		t.Logf("%+v", tc)
		ls := NewLineSegment(tc.from, tc.to)
		got := ls.BoundingRectContains(tc.p)
		if tc.expected != got {
			t.Fatalf("Failed : got %v expected %v", got, tc.expected)
		}
	}
}

func TestRegionContains(t *testing.T) {
	unitSquare := NewRegion(
		[]Pos{
			Pos{0, 0},
			Pos{1, 0},
			Pos{1, 1},
			Pos{0, 1},
		},
	)

	// Origin is special, this doesn't contain origin
	// or have corners aligned
	otherSquare := NewRegion(
		[]Pos{
			Pos{2, 2},
			Pos{2, 3},
			Pos{3, 2},
			Pos{3, 3},
		},
	)

	testCases := []struct {
		region   Region
		p        Pos
		expected bool
	}{
		// Put current failure first, so logging shows failure
		{unitSquare, Pos{-1, 0}, false},
		// ...

		{unitSquare, Pos{0.3, 0.4}, true},
		{unitSquare, Pos{1.3, 0.4}, false},
		{unitSquare, Pos{0.3, 1.4}, false},
		{unitSquare, Pos{-0.3, 0.4}, false},
		{unitSquare, Pos{0.3, -0.4}, false},

		// Corners included
		{unitSquare, Pos{0, 0}, true},
		{unitSquare, Pos{1, 0}, true},
		{unitSquare, Pos{0, 1}, true},
		{unitSquare, Pos{1, 1}, true},

		// Sides included
		{unitSquare, Pos{0, 0.5}, true},
		{unitSquare, Pos{1, 0.5}, true},
		{unitSquare, Pos{0.5, 0}, true},
		{unitSquare, Pos{0.5, 1}, true},

		{unitSquare, Pos{3, 2}, false},
		{unitSquare, Pos{-3, -2}, false},
		{unitSquare, Pos{3, -2}, false},
		{unitSquare, Pos{-3, 2}, false},

		{unitSquare, Pos{2, 0}, false},
		{unitSquare, Pos{0, 2}, false},
		{unitSquare, Pos{2, 2}, false},

		//		{unitSquare, Pos{-1, 0}, false},
		{unitSquare, Pos{0, -1}, false},
		{unitSquare, Pos{-1, -1}, false},

		{unitSquare, Pos{0.5, 0.5}, true},

		{unitSquare, Pos{0, 0.5}, true},
		{unitSquare, Pos{0, -0.5}, false},

		{unitSquare, Pos{0.5, 0}, true},
		{unitSquare, Pos{-0.5, 0}, false},

		{unitSquare, Pos{1.5, 1.5}, false},

		{otherSquare, Pos{0, 0}, false},
		{otherSquare, Pos{1, 0}, false},
		{otherSquare, Pos{1, 1}, false},
		{otherSquare, Pos{0, 1}, false},

		{otherSquare, Pos{2, 1}, false},

		{otherSquare, Pos{2.1, 2.1}, true},
		{otherSquare, Pos{2.5, 2.5}, true},
	}

	for _, tc := range testCases {
		t.Logf("TC: %v", tc)
		got := tc.region.Contains(tc.p)
		if got != tc.expected {
			t.Fatalf("Got %v expected %v", got, tc.expected)
		}
	}
}
