// Copyright 2015 Mitchell Kember. Subject to the MIT License.

package system

import (
	"bytes"
	"fmt"
	"math"
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
	segments := s.cacheRender(s.min + opts.Depth)
	view := calcViewBox(segments, opts)
	thickness := opts.Thickness * math.Max(view.w, view.h) / StepFactor
	p := fmtPrecision(opts.Precision)

	var buf bytes.Buffer
	fmt.Fprintf(
		&buf,
		"<svg xmlns='http://www.w3.org/2000/svg' viewBox='%s %s %s %s'>"+
			"<defs><style>polyline { fill: none; stroke-linecap: square; "+
			"stroke-width: %s; stroke: %s; }</style></defs>",
		p(view.x), p(view.y), p(view.w), p(view.h), p(thickness), opts.Color,
	)

	for _, points := range segments {
		buf.WriteString("<polyline points='")
		for _, pt := range points {
			fmt.Fprintf(&buf, "%s,%s ", p(pt.X), p(pt.Y))
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

// PadFactor is used to add additional padding to the viewBox.
const PadFactor = 0.8

// calcViewBox returns a viewBox large enough to contain all the points.
func calcViewBox(segments [][]vector, opts *Options) viewBox {
	var xMin, xMax, yMin, yMax float64
	for _, points := range segments {
		for _, pt := range points {
			if pt.X < xMin {
				xMin = pt.X
			} else if pt.X > xMax {
				xMax = pt.X
			}
			if pt.Y < yMin {
				yMin = pt.Y
			} else if pt.Y > yMax {
				yMax = pt.Y
			}
		}
	}

	edge := PadFactor * opts.Thickness
	xMin -= edge
	yMin -= edge
	xMax += edge
	yMax += edge

	width := xMax - xMin
	height := yMax - yMin

	return viewBox{xMin, yMin, width, height}
}

// fmtPrecision returns a function that formats floats as strings with the given
// number of decimal places, omitting trailing zeros.
func fmtPrecision(precision int) func(float64) string {
	return func(f float64) string {
		str := strconv.FormatFloat(f, 'f', precision, 64)
		str = strings.TrimRight(str, "0")
		str = strings.TrimRight(str, ".")
		return str
	}
}
