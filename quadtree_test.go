package quadtree

import (
	"math/rand"
	"reflect"
	"testing"
	"time"
)

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
		want := tp.contains
		got := b.ContainsPoint(tp.p)

		if want != got {
			t.Errorf("expected %v got %v for %+v", want, got, tp.p)
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

	sw := q.SouthEast.Boundary
	if sw.Start.X != b.Start.X+w/2 || sw.Start.Y != b.Start.Y+h/2 {
		t.Errorf("expected SE start to be at (%d, %d), got (%d, %d)", b.Start.X+w/2, b.Start.Y+h/2, sw.Start.X, sw.Start.Y)
	}

	se := q.SouthWest.Boundary
	if se.Start.X != b.Start.X || se.Start.Y != b.Start.Y+h/2 {
		t.Errorf("expected SW start to be at (%d, %d), got (%d, %d)", b.Start.X, b.Start.Y+h/2, se.Start.X, se.Start.Y)
	}
}

func TestInsert(t *testing.T) {
	t.Run("when point outside of boundary", func(t *testing.T) {
		b := &Boundary{Start: &Point{0, 0}, Width: 10, Height: 10}
		q := &Node{Boundary: b, Capacity: 1}
		p := &Point{X: 20, Y: 20}
		q.Insert(p)

		want := 0
		got := len(q.Points)

		if want != got {
			t.Errorf("expected number of points to equal %d, got %d", got, want)
		}
	})

	t.Run("when capacity hasn't exceeded", func(t *testing.T) {
		b := &Boundary{Start: &Point{0, 0}, Width: 10, Height: 10}
		q := &Node{Boundary: b, Capacity: 1}
		p := &Point{X: 0, Y: 0}
		q.Insert(p)

		want := 1
		got := len(q.Points)

		if q.NorthWest != nil {
			t.Fatalf("expected to not run subdivide")
		}

		if want != got {
			t.Errorf("expected number of points to equal %d, got %d", got, want)
		}
	})

	t.Run("when capacity has been exceeded", func(t *testing.T) {
		b := &Boundary{Start: &Point{0, 0}, Width: 10, Height: 10}
		q := &Node{Boundary: b, Capacity: 1}
		p := &Point{X: 1, Y: 1}
		q.Insert(p)
		p = &Point{X: 9, Y: 9}
		q.Insert(p)

		if q.NorthWest == nil {
			t.Fatalf("expected to run subdivide")
		}

		want := 1
		got := len(q.SouthEast.Points)

		if want != got {
			t.Errorf("expected number of SE points to equal %d, got %d", got, want)
		}
	})
}

func TestIntersectsWith(t *testing.T) {
	b := &Boundary{Start: &Point{30, 30}, Width: 20, Height: 20}
	testBoundaries := []struct {
		bound      *Boundary
		intersects bool
	}{
		{&Boundary{&Point{0, 0}, 10, 10}, false},
		{&Boundary{&Point{35, 35}, 10, 10}, true},
		{&Boundary{&Point{50, 50}, 5, 5}, false},
		{&Boundary{&Point{49, 49}, 20, 20}, true},
	}

	for _, tb := range testBoundaries {
		want := tb.intersects
		got := b.IntersectsWith(tb.bound)

		if want != got {
			t.Errorf("expected %v got %v for %+v", want, got, tb.bound)
		}
	}
}

func TestQuery(t *testing.T) {
	const cap = 1
	const w = 5
	const h = 6
	b := &Boundary{Start: &Point{0, 0}, Width: w, Height: h}
	q := &Node{Boundary: b, Capacity: cap}
	qb := &Boundary{&Point{1, 1}, 2, 2}

	for i := 0; i < w; i++ {
		for j := 0; j < h; j++ {
			q.Insert(&Point{i, j})
		}
	}

	points := make([]*Point, 0)
	q.Query(qb, &points)
	got := points

	expected := []*Point{&Point{1, 1}, &Point{2, 1}, &Point{1, 2}, &Point{2, 2}}

	if len(got) != len(expected) {
		t.Fatalf("expected to have %d points got %d", len(expected), len(got))
	}

	for _, e := range expected {
		found := false

		for _, g := range got {
			if reflect.DeepEqual(g, e) {
				found = true
				break
			}
		}

		if !found {
			t.Errorf("expected to got %v point", e)
		}
	}
}

func BenchmarkQuery(bench *testing.B) {
	rand.Seed(time.Now().UnixNano())

	width, height := 800, 800
	cap := 32
	region := 259
	b := &Boundary{Start: &Point{0, 0}, Width: width, Height: height}
	q := &Node{Boundary: b, Capacity: cap}
	qb := &Boundary{&Point{0, 0}, region, region}
	array := make([][]bool, height)

	for y := 0; y < height; y++ {
		array[y] = make([]bool, width)

		for x := 0; x < width; x++ {
			value := rand.Intn(2) == 1
			array[y][x] = value

			if value {
				q.Insert(&Point{x, y})
			}
		}
	}

	points := make([]*Point, 0)

	bench.ResetTimer()

	for i := 0; i < bench.N; i++ {
		q.Query(qb, &points)
	}
}
