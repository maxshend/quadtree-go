package quadtree

import "testing"

func TestContainsPoint(t *testing.T) {
	b := Boundary{Start: &Point{0, 0}, Width: 10, Height: 10}
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
