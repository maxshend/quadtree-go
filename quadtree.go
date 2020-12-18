package quadtree

import "math"

// Point represents a point in a 2-dimensional space
type Point struct {
	X int
	Y int
}

// Boundary represents a Quadtree boundary
type Boundary struct {
	Start  *Point
	Width  int
	Height int
}

// ContainsPoint checks whether a point inside the boundary
func (b *Boundary) ContainsPoint(p *Point) bool {
	return (p.X >= b.Start.X && p.X < b.Start.X+b.Width && p.Y >= b.Start.Y && p.Y < b.Start.Y+b.Height)
}

// IntersectsWith checks whether two boundaries are intersect with each other
func (b *Boundary) IntersectsWith(bound *Boundary) bool {
	return !(b.Start.X >= bound.Start.X+bound.Width ||
		bound.Start.X >= b.Start.X+b.Width ||
		b.Start.Y >= bound.Start.Y+bound.Height ||
		bound.Start.Y >= b.Start.Y+b.Height)
}

// Node represents a Quadtree node
type Node struct {
	Boundary *Boundary
	Capacity int
	Points   []*Point

	NorthWest *Node
	NorthEast *Node
	SouthWest *Node
	SouthEast *Node
}

// Subdivide subdivides a Quadtree node
func (n *Node) Subdivide() {
	var b *Boundary
	var p *Point
	w := n.Boundary.Width / 2
	h := n.Boundary.Height / 2
	ww := int(math.Ceil(float64(n.Boundary.Width) / 2))
	hh := int(math.Ceil(float64(n.Boundary.Height) / 2))
	origin := n.Boundary.Start

	p = &Point{X: origin.X + w, Y: origin.Y}
	b = &Boundary{Width: w, Height: h, Start: p}
	n.NorthWest = &Node{Boundary: b, Capacity: n.Capacity}

	p = &Point{X: origin.X, Y: origin.Y}
	b = &Boundary{Width: ww, Height: h, Start: p}
	n.NorthEast = &Node{Boundary: b, Capacity: n.Capacity}

	p = &Point{X: origin.X + w, Y: origin.Y + h}
	b = &Boundary{Width: w, Height: hh, Start: p}
	n.SouthEast = &Node{Boundary: b, Capacity: n.Capacity}

	p = &Point{X: origin.X, Y: origin.Y + h}
	b = &Boundary{Width: ww, Height: hh, Start: p}
	n.SouthWest = &Node{Boundary: b, Capacity: n.Capacity}
}

// Insert inserts a point to a Quadtree Node
func (n *Node) Insert(p *Point) (success bool) {
	if !n.Boundary.ContainsPoint(p) {
		return false
	}

	if len(n.Points) < n.Capacity && n.NorthWest == nil {
		n.Points = append(n.Points, p)
		return true
	}

	if n.NorthWest == nil {
		n.Subdivide()
	}

	if n.NorthWest.Insert(p) {
		return true
	}
	if n.NorthEast.Insert(p) {
		return true
	}
	if n.SouthWest.Insert(p) {
		return true
	}
	if n.SouthEast.Insert(p) {
		return true
	}

	return false
}

// Query returns all points within a boundary
func (n *Node) Query(r *Boundary) (points []*Point) {
	points = make([]*Point, 0)

	if !n.Boundary.IntersectsWith(r) {
		return
	}

	for _, p := range n.Points {
		if r.ContainsPoint(p) {
			points = append(points, p)
		}
	}

	if n.NorthWest == nil {
		return
	}

	points = append(points, n.NorthWest.Query(r)...)
	points = append(points, n.NorthEast.Query(r)...)
	points = append(points, n.SouthWest.Query(r)...)
	points = append(points, n.SouthEast.Query(r)...)

	return
}
