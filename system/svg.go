// Copyright 2015 Mitchell Kember. Subject to the MIT License.

package system

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
)

// Options are a collection of settings for rendering systems as SVG.
type Options struct {
	Depth     int     // Depth of recursion.
	Thickness float64 // Stroke thickness.
	Color     string  // Stroke color.
	Precision int     // Float precision (number of decimal places).
}

// SVG renders the system as a curve, returning a string of SVG data.
func (s *System) SVG(opts *Options) string {
	segments := s.render(opts.Depth)
	view := calcViewBox(segments, opts)

	var buf bytes.Buffer
	fmt.Fprintf(
		&buf,
		"<svg xmlns='http://www.w3.org/2000/svg' viewBox='%s %s %s %s'>",
		precise(view.x, opts.Precision),
		precise(view.y, opts.Precision),
		precise(view.w, opts.Precision),
		precise(view.h, opts.Precision),
	)

	thickness := opts.Thickness * view.w / stepFactor
	polyTag := fmt.Sprintf(
		"<polyline fill='none' stroke-linecap='square' stroke='%s' "+
			"stroke-width='%s' points='",
		opts.Color,
		precise(thickness, opts.Precision))

	for _, points := range segments {
		buf.WriteString(polyTag)
		for _, pt := range points {
			x := precise(pt.x, opts.Precision)
			y := precise(pt.y, opts.Precision)
			fmt.Fprintf(&buf, "%s,%s ", x, y)
		}
		buf.WriteString("'/>")
	}
	buf.WriteString("</svg>")

	return buf.String()
}

// A viewBox describes the boundaries of an SVG image.
type viewBox struct {
	x, y, w, h float64
}

// paddingFactor is used to add additional padding to the viewBox.
const paddingFactor = 1.1

// calcViewBox returns a viewBox large enough to contain all the points.
func calcViewBox(segments [][]vector, opts *Options) viewBox {
	var xMin, xMax, yMin, yMax float64
	for _, points := range segments {
		for _, pt := range points {
			if pt.x < xMin {
				xMin = pt.x
			} else if pt.x > xMax {
				xMax = pt.x
			}
			if pt.y < yMin {
				yMin = pt.y
			} else if pt.y > yMax {
				yMax = pt.y
			}
		}
	}

	edge := paddingFactor * opts.Thickness / 2
	xMin -= edge
	yMin -= edge
	xMax += edge
	yMax += edge

	width := xMax - xMin
	height := yMax - yMin

	if width < height {
		diff := height - width
		width += diff
		xMin -= diff / 2
	}
	return viewBox{xMin, yMin, width, height}
}

// precise returns a float formatted as a string with the given number of
// decimal places, omitting trailing zeros.
func precise(f float64, precision int) string {
	str := strconv.FormatFloat(f, 'f', precision, 64)
	str = strings.TrimRight(str, "0")
	str = strings.TrimRight(str, ".")
	return str
}
