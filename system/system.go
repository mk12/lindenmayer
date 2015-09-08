// Copyright 2015 Mitchell Kember. Subject to the MIT License.

// Package system implements SVG rendering of Lindenmayer systems
package system

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"log"
	"os"
	"unicode"
)

// A rewriteSet is a set of rules for replacing variables with sequences.
type rewriteSet map[byte][]byte

// A System is a representation of a Lindenmayer system (L-system).
type System struct {
	axiom []byte     // Initial state of the system.
	rules rewriteSet // Rules for rewriting variables.
	angle float64    // Turn angle, in radians.
	start float64    // Initial direction, in radians (standard position).
	turn  bool       // Turn initial direction (used for dragon).
	base  float64    // Base b in y ~ b^x where y is size and x is depth.
	min   int        // Number of initial depths to skip.
	max   int        // Maximum depth of recursion.
}

// cacheRender is like render but caches results to disk for better performance.
func (s *System) cacheRender(depth int) [][]vector {
	name := s.name()
	if name == "" {
		return s.render(depth)
	}

	filename := fmt.Sprintf("/app/cache/%s-%d.gob", name, depth)
	_, err := os.Stat(filename)
	if err == nil {
		log.Println("Cache hit:", filename)
		file, err := os.Open(filename)
		if err != nil {
			log.Println("File open error:", err)
			return s.render(depth)
		}
		defer file.Close()

		var segments [][]vector
		err = gob.NewDecoder(file).Decode(&segments)
		if err != nil {
			log.Println("Gob decode error:", err)
			return s.render(depth)
		}
		return segments
	}
	if os.IsNotExist(err) {
		log.Println("Cache miss:", filename)
		segments := s.render(depth)
		file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
		if err != nil {
			log.Println("File create error:", err)
			return segments
		}
		defer file.Close()

		gob.NewEncoder(file).Encode(segments)
		return segments
	}

	log.Println("File stat error:", err)
	return s.render(depth)
}

// render draws the curve for the system at the given depth of recursion. It
// returns a list of polygonal lines as a 2D slice of points.
func (s *System) render(depth int) [][]vector {
	var buf bytes.Buffer
	s.expand(s.axiom, depth, &buf)
	return s.execute(buf.Bytes(), depth)
}

// expand recursively expands seq using the system's rules, stopping when the
// level reaches zero. It writes the results into buf.
func (s *System) expand(seq []byte, level int, buf *bytes.Buffer) {
	if level == 0 {
		buf.Write(seq)
		return
	}

	for _, sym := range seq {
		if replacement, ok := s.rules[sym]; ok {
			s.expand(replacement, level-1, buf)
		} else {
			buf.WriteByte(sym)
		}
	}
}

// execute performs the instructions in seq, ignoring symbols with no special
// meaning. It returns a list of polygonal lines, each one being a sequence of
// points that the were visited.
func (s *System) execute(seq []byte, depth int) [][]vector {
	var segments [][]vector
	var stack []turtle
	turt := s.initialTurtle(depth)
	points := []vector{turt.pos}

	for _, sym := range seq {
		if unicode.IsUpper(rune(sym)) {
			turt.advance()
			points = append(points, turt.pos)
			continue
		}
		switch sym {
		case '+':
			turt.turnCCW()
		case '-':
			turt.turnCW()
		case '[':
			stack = append(stack, turt)
		case ']':
			last := len(stack) - 1
			turt = stack[last]
			stack = stack[:last]
			if len(points) > 0 {
				segments = append(segments, points)
			}
			points = []vector{turt.pos}
		}
	}

	if len(points) > 0 {
		segments = append(segments, points)
	}
	return segments
}
