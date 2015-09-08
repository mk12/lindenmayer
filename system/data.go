// Copyright 2015 Mitchell Kember. Subject to the MIT License.

package system

import "math"

// Named returns the system with the given name, or nil if it doesn't exist.
func Named(name string) *System {
	if sys, ok := namedSystems[name]; ok {
		return sys
	}
	return nil
}

// name returns the name of the system, if it is an instance in namedSystems.
// Otherwise, returns the empty string.
func (s *System) name() string {
	for key, sys := range namedSystems {
		if s == sys {
			return key
		}
	}
	return ""
}

// MaxDepth returns the maximum depth the system should be rendered at.
func (s *System) MaxDepth() int {
	return s.max - s.min
}

// namedSystems contains some common Lindenmayer systems.
var namedSystems = map[string]*System{
	"koch": {
		axiom: []byte("F++F++F"),
		rules: rewriteSet{
			'F': []byte("F-F++F-F"),
		},
		angle: math.Pi / 3,
		start: 0,
		turn:  false,
		base:  3,
		min:   0,
		max:   7,
	},
	"hilbert": {
		axiom: []byte("a"),
		rules: rewriteSet{
			'a': []byte("+bF-aFa-Fb+"),
			'b': []byte("-aF+bFb+Fa-"),
		},
		angle: math.Pi / 2,
		start: 0,
		turn:  false,
		base:  2,
		min:   1,
		max:   8,
	},
	"peano": {
		axiom: []byte("a"),
		rules: rewriteSet{
			'a': []byte("aFbFa-F-bFaFb+F+aFbFa"),
			'b': []byte("bFaFb+F+aFbFa-F-bFaFb"),
		},
		angle: math.Pi / 2,
		start: math.Pi / 2,
		turn:  false,
		base:  3,
		min:   1,
		max:   5,
	},
	"gosper": {
		axiom: []byte("A"),
		rules: rewriteSet{
			'A': []byte("A-B--B+A++AA+B-"),
			'B': []byte("+A-BB--B-A++A+B"),
		},
		angle: math.Pi / 3,
		start: math.Pi / 9,
		turn:  true,
		base:  2.6,
		min:   0,
		max:   5,
	},
	"sierpinski": {
		axiom: []byte("A"),
		rules: rewriteSet{
			'A': []byte("+B-A-B+"),
			'B': []byte("-A+B+A-"),
		},
		angle: math.Pi / 3,
		start: 0,
		turn:  false,
		base:  2,
		min:   1,
		max:   9,
	},
	"rings": {
		axiom: []byte("F+F+F+F"),
		rules: rewriteSet{
			'F': []byte("FF+F+F+F+F+F-F"),
		},
		angle: math.Pi / 2,
		start: -37 * math.Pi / 360,
		turn:  true,
		base:  3,
		min:   0,
		max:   5,
	},
	"tree": {
		axiom: []byte("A"),
		rules: rewriteSet{
			'A': []byte("B[+A]-A"),
			'B': []byte("BB"),
		},
		angle: math.Pi / 4,
		start: math.Pi / 2,
		turn:  false,
		base:  1.9,
		min:   0,
		max:   9,
	},
	"plant": {
		axiom: []byte("a"),
		rules: rewriteSet{
			'a': []byte("F+[[a]-a]-F[-Fa]+a"),
			'F': []byte("FF"),
		},
		angle: 25.0 / 180.0 * math.Pi,
		start: math.Pi / 4,
		turn:  false,
		base:  2,
		min:   1,
		max:   7,
	},
	"willow": {
		axiom: []byte("a"),
		rules: rewriteSet{
			'a': []byte("bFF[+a]c"),
			'b': []byte("bF"),
			'c': []byte("bFF[-a]a"),
		},
		angle: math.Pi / 6,
		start: 80.0 / 180 * math.Pi,
		turn:  false,
		base:  1.3,
		min:   1,
		max:   12,
	},
	"dragon": {
		axiom: []byte("Fa"),
		rules: rewriteSet{
			'a': []byte("a-bF-"),
			'b': []byte("+Fa+b"),
		},
		angle: math.Pi / 2,
		start: math.Pi / 4,
		turn:  true,
		base:  1.4,
		min:   0,
		max:   15,
	},
	"island": {
		axiom: []byte("F+F+F+F"),
		rules: rewriteSet{
			'F': []byte("F+F-F-FF+F+F-F"),
		},
		angle: math.Pi / 2,
		start: math.Pi / 4,
		turn:  false,
		base:  4,
		min:   0,
		max:   4,
	},
}
