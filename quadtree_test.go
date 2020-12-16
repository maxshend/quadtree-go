package quadtree

import "testing"

func TestContainsPoint(t *testing.T) {
	b := &Boundary{Start: &Point{0, 0}, Width: 10, Height: 10}
	testPoints := []struct {
		p        *Point
		contains bool
	}{
		{&Point{-1, -1}, false},
		{&Point{0, 0}, true},
		{&Point{5, 20}, false},
		{&Point{3, 4}, true},
		{&Point{10, 10}, false},
	}

	for _, tp := range testPoints {
		got := b.ContainsPoint(tp.p)
		if tp.contains != got {
			t.Errorf("expected %v got %v", tp.contains, got)
		}
	}
}

func TestSubdivide(t *testing.T) {
	const cap = 5
	const w = 20
	const h = 10
	b := &Boundary{Start: &Point{0, 0}, Width: w, Height: h}
	q := &Node{Boundary: b, Capacity: cap}
	q.Subdivide()

	if q.NorthWest == nil || q.NorthEast == nil || q.SouthWest == nil || q.SouthEast == nil {
		t.Errorf("exptected to create 4 subdivisions")
	}

	if q.NorthWest.Capacity != cap || q.NorthEast.Capacity != cap || q.SouthWest.Capacity != cap || q.SouthEast.Capacity != cap {
		t.Errorf("exptected subdivisions have capacity %d", cap)
	}

	if q.NorthWest.Boundary.Width != w/2 || q.NorthEast.Boundary.Width != w/2 || q.SouthWest.Boundary.Width != w/2 || q.SouthEast.Boundary.Width != w/2 {
		t.Errorf("exptected subdivisions have width %d", w/2)
	}

	if q.NorthWest.Boundary.Height != h/2 || q.NorthEast.Boundary.Height != h/2 || q.SouthWest.Boundary.Height != h/2 || q.SouthEast.Boundary.Height != h/2 {
		t.Errorf("exptected subdivisions have height %d", h/2)
	}

	nw := q.NorthWest.Boundary
	if nw.Start.X != b.Start.X+w/2 || nw.Start.Y != b.Start.Y {
		t.Errorf("expected NW start to be at (%d, %d), got (%d, %d)", b.Start.X+w/2, b.Start.Y, nw.Start.X, nw.Start.Y)
	}

	ne := q.NorthEast.Boundary
	if ne.Start.X != b.Start.X || ne.Start.Y != b.Start.Y {
		t.Errorf("expected NE start to be at (%d, %d), got (%d, %d)", b.Start.X, b.Start.Y, ne.Start.X, ne.Start.Y)
	}

	sw := q.SouthWest.Boundary
	if sw.Start.X != b.Start.X+w/2 || sw.Start.Y != b.Start.Y+h/2 {
		t.Errorf("expected SW start to be at (%d, %d), got (%d, %d)", b.Start.X+w/2, b.Start.Y+h/2, sw.Start.X, sw.Start.Y)
	}

	se := q.SouthEast.Boundary
	if se.Start.X != b.Start.X || se.Start.Y != b.Start.Y+h/2 {
		t.Errorf("expected SE start to be at (%d, %d), got (%d, %d)", b.Start.X, b.Start.Y+h/2, se.Start.X, se.Start.Y)
	}
}

func TestInsert(t *testing.T) {
	b := &Boundary{Start: &Point{0, 0}, Width: 10, Height: 10}
	q := &Node{Boundary: b, Capacity: 1}
	p := &Point{X: 0, Y: 0}
	q.Insert(p)

	if len(q.Points) != 1 {
		t.Errorf("expected number of points to equal 1, got %d", len(q.Points))
	}

	// TODO: Add more test cases
}
