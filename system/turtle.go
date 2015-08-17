// Copyright 2015 Mitchell Kember. Subject to the MIT License.

package system

import "math"

// A vector is a two-dimensional Euclidean vector with x and y components.
type vector struct {
	x, y float64
}

// A turtle is a friendly animal that controls drawing.
type turtle struct {
	pos   vector  // Current position.
	vel   vector  // Current velocity.
	dir   float64 // Current direction, in radians (standard position).
	step  float64 // Step size for moving forward.
	angle float64 // Turning angle, in radians.
}

// stepFactor is a rought estimate of how far the turtle should travel. It is
// used to set the turtle's step such that its position always stays near this
// order of magnitude.
const stepFactor = 600

// Returns the turtle to begin with when rendering the system.
func (s *System) initialTurtle(depth int) turtle {
	t := turtle{
		step:  stepFactor * math.Pow(s.growth, float64(-depth)),
		angle: s.angle,
	}
	t.rotate(s.start)
	return t
}

// advance moves the turtle forward by one step.
func (t *turtle) advance() {
	t.pos.x += t.vel.x
	t.pos.y -= t.vel.y
}

// turnCCW turns the turtle counterclockwise by one step.
func (t *turtle) turnCCW() {
	t.rotate(t.angle)
}

// turnCCW turns the turtle clockwise by one step.
func (t *turtle) turnCW() {
	t.rotate(-t.angle)
}

// rotate changes the turtle's direction by delta (in radians).
func (t *turtle) rotate(delta float64) {
	t.dir += delta
	t.vel.x = t.step * math.Cos(t.dir)
	t.vel.y = t.step * math.Sin(t.dir)
}
