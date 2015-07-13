// Copyright 2015 Mitchell Kember. Subject to the MIT License.

package curve

import (
	"bytes"
	"fmt"
	"github.com/ajstarks/svgo"
	"strconv"
	"strings"
)

// A point is a pair of coordinates in the unit square. Both x and y should be
// on the interval [0,1]. The origin is in the top-left corner.
type point struct {
	x, y float64
}

// A drawFn is a function that draws a space-filling curve to a certain depth
// of recursion, returning a sequence of points.
type drawFn func(depth int) []point

// allDrawFns maps from names (URL handles) to curve drawing functions.
var allDrawFns = map[string]drawFn{
	"hilbert": drawHilbert,
	"peano":   drawPeano,
}

// SVG renders a space-filling curve with the given name and depth, and returns
// it as a string of HTML. Returns the empty string if there is an error.
func SVG(name, depth string) (html string, err error) {
	n, err := strconv.Atoi(depth)
	if err != nil {
		return
	}
	if n < 0 {
		err = fmt.Errorf("invalid depth %d", depth)
		return
	}
	draw, ok := allDrawFns[name]
	if !ok {
		err = fmt.Errorf("invalid curve '%s'", name)
		return
	}

	html = connectDots(draw(n))
	return
}

// connectDots draws an SVG polygonal line by connecting the dots in the given
// sequence, and returns it as a string of HTML.
func connectDots(dots []point) string {
	count := len(dots)
	if count == 0 {
		return ""
	}

	// TODO: change the 500 (use viewPort, not width & height)
	xs := make([]int, count)
	ys := make([]int, count)
	for i, dot := range dots {
		xs[i] = int(dot.x * 500)
		ys[i] = int(dot.y * 500)
	}

	var buf bytes.Buffer
	s := svg.New(&buf)
	s.Start(500, 500)
	s.Polyline(xs, ys, "fill:none;stroke:black")
	s.End()
	return removeLines(buf.String(), 2)
}

// removeLines returns a string with the first n lines removed.
func removeLines(s string, n int) string {
	for i := 0; i < n; i++ {
		index := strings.IndexByte(s, '\n')
		if index == -1 {
			return s
		}
		if index+1 >= len(s) {
			return ""
		}
		s = s[index+1:]
	}
	return s
}
