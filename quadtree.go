package quadtree

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

// Node represents a Quadtree node
type Node struct {
	Boundary Boundary
	Capacity int
	Points   []Point

	NorthWest *Node
	NorthEast *Node
	SouthWest *Node
	SouthEast *Node
}
