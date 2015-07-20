// Copyright 2015 Mitchell Kember. Subject to the MIT License.

package curve

import "math"

// A hilbertShape is one of four shapes used in Hilbert's space-filling curve.
// Each is a square with one of its four sides missing.
type hilbertShape int

// A quadrant designates one of the four parts of the unit square.
type quadrant int

const (
	hsUp hilbertShape = iota
	hsDown
	hsLeft
	hsRight
)

const (
	qTL quadrant = iota
	qTR
	qBL
	qBR
)

// A hilbertRule is an instruction to a shape in a quadrant.
type hilbertRule struct {
	quad  quadrant
	shape hilbertShape
}

// hilbertGrammar contains the rules for drawing Hilbert's space-filling curve.
// It describes how to replace each of the four shapes with a further level of
// recursion. For example, hsUp should be replaced by following four rules,
// starting by drawing hsRight in its top-right quadrant. The rules must be
// followed in order for the points to connect properly.
var hilbertGrammar = map[hilbertShape][4]hilbertRule{
	hsUp:    {{qTR, hsRight}, {qBR, hsUp}, {qBL, hsUp}, {qTL, hsLeft}},
	hsDown:  {{qBL, hsLeft}, {qTL, hsDown}, {qTR, hsDown}, {qBR, hsRight}},
	hsLeft:  {{qBL, hsDown}, {qBR, hsLeft}, {qTR, hsLeft}, {qTL, hsUp}},
	hsRight: {{qTR, hsUp}, {qTL, hsRight}, {qBL, hsRight}, {qBR, hsDown}},
}

// drawHilbert approximates Hilbert's space-filling curve to the given depth of
// recursion, returning a sequence of points to be connected.
func drawHilbert(depth int) []point {
	if depth < 0 {
		return []point{}
	}

	numDots := int(math.Pow(4, float64(depth)))
	dots := make([]point, 0, numDots)
	centre := point{0.5, 0.5}
	return hilbert(dots, centre, hsDown, 0, depth)
}

// toQuad returns a new point in the centre of one of the four quadrants. The
// size of the current square is determined by the level of recursion.
func (pt point) toQuad(quad quadrant, level int) point {
	s := math.Pow(2, float64(-level))
	switch quad {
	case qTL:
		pt.x -= s
		pt.y -= s
	case qTR:
		pt.x += s
		pt.y -= s
	case qBL:
		pt.x -= s
		pt.y += s
	case qBR:
		pt.x += s
		pt.y += s
	}
	return pt
}

// hilbert recursively draws Hilbert's space-filling curve by adding points to
// the dots slice (which is assumed to have the correct capacity).
func hilbert(dots []point, pos point, shape hilbertShape, level, depth int) []point {
	if level >= depth {
		return append(dots, pos)
	}

	for _, rule := range hilbertGrammar[shape] {
		newPos := pos.toQuad(rule.quad, level+2)
		dots = hilbert(dots, newPos, rule.shape, level+1, depth)
	}
	return dots
}
