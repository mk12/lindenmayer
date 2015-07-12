package main

import (
	"github.com/ajstarks/svgo"
	"math"
)

// There are four elements repeated in Hilbert's space-filling curve. Each is a
// square with one of its sides missing. I use that side to identify them.
const (
	U = iota // opens up
	D        // opens down
	L        // opens left
	R        // opens right
)

// The grid is recursively divided into four quadrants.
const (
	TL = iota // top left
	TR        // top right
	BL        // bottom left
	BR        // bottom right
)

// The grammar describes the construction of the curve. There are four elements,
// and each has a line instructing how to replace it with a further level of
// recursion. For example, the rule {TR, R} means that a right-opening element
// should be drawn in the top-right corner. The four rules of the line must be
// followed in order.
var grammar = [4][4][2]int{
	{{TR, R}, {BR, U}, {BL, U}, {TL, L}}, // U
	{{BL, L}, {TL, D}, {TR, D}, {BR, R}}, // D
	{{BL, D}, {BR, L}, {TR, L}, {TL, U}}, // L
	{{TR, U}, {TL, R}, {BL, R}, {BR, D}}, // R
}

// A Hilbert stores the points in a Hilbert space-filling curve.
type Hilbert struct {
	xs    []int // list of x-coordinates
	ys    []int // list of y-coordinates
	depth int   // maximum recursion depth
}

// Creates a new Hilbert with the specified depth, allocating enough space for
// all the points required for that depth.
func newHilbert(depth int) Hilbert {
	num := int(math.Pow(4, float64(depth)))
	return Curve{
		xs:    make([]int, 0, num),
		ys:    make([]int, 0, num),
		depth: depth,
	}
}

func renderHilbert(depth int, s *svg.SVG) {
}
