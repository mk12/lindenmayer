// Copyright 2015 Mitchell Kember. Subject to the MIT License.

package system

import "math"

// A vector is a two-dimensional Euclidean vector with x and y components.
type vector struct {
	X, Y float64
}

// A turtle is a friendly animal that controls drawing.
type turtle struct {
	pos   vector  // Current position.
	vel   vector  // Current velocity.
	dir   float64 // Current direction, in radians (standard position).
	step  float64 // Step size for moving forward.
	angle float64 // Turning angle, in radians.
}

// StepFactor is a rought estimate of how far the turtle should travel. It is
// used to set the turtle's step such that its position always stays near this
// order of magnitude.
const StepFactor = 600

// Returns the turtle to begin with when rendering the system.
func (s *System) initialTurtle(depth int) turtle {
	t := turtle{
		step:  StepFactor * math.Pow(s.base, float64(-depth)),
		angle: s.angle,
	}
	if s.turn {
		t.rotate(s.start * float64(depth))
	} else {
		t.rotate(s.start)
	}

	return t
}

// advance moves the turtle forward by one step.
func (t *turtle) advance() {
	t.pos.X += t.vel.X
	t.pos.Y -= t.vel.Y
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
	t.vel.X = t.step * math.Cos(t.dir)
	t.vel.Y = t.step * math.Sin(t.dir)
}
