// Copyright 2015 Mitchell Kember. Subject to the MIT License.

package curve

import "testing"

func Test_drawHilbert(t *testing.T) {
	table := []struct {
		depth int
		dots  []point
	}{
		{-2, []point{}},
		{-1, []point{}},
		{0, []point{{0.5, 0.5}}},
		{1, []point{{0.25, 0.75}, {0.25, 0.25}, {0.75, 0.25}, {0.75, 0.75}}},
		{2, []point{
			{0.125, 0.875}, {0.375, 0.875}, {0.375, 0.625}, {0.125, 0.625},
			{0.125, 0.375}, {0.125, 0.125}, {0.375, 0.125}, {0.375, 0.375},
			{0.625, 0.375}, {0.625, 0.125}, {0.875, 0.125}, {0.875, 0.375},
			{0.875, 0.625}, {0.625, 0.625}, {0.625, 0.875}, {0.875, 0.875},
		}},
	}
	for i, test := range table {
		dots := drawHilbert(test.depth)
		if !samePoints(dots, test.dots) {
			t.Errorf("[%d] got %v, want %v", i, dots, test.dots)
		}
	}
}

func Test_point_toQuad(t *testing.T) {
	table := []struct {
		in    point
		quad  quadrant
		level int
		out   point
	}{
		{point{0, 0}, qTL, 0, point{-1, -1}},
		{point{0, 0}, qTR, 0, point{1, -1}},
		{point{0, 0}, qBL, 0, point{-1, 1}},
		{point{0, 0}, qBR, 0, point{1, 1}},
		{point{1, 1}, qBR, 1, point{1.5, 1.5}},
		{point{-1, -1}, qTR, 2, point{-0.75, -1.25}},
		{point{0.5, 0.5}, qTL, 4, point{0.4375, 0.4375}},
	}
	for i, test := range table {
		out := test.in.toQuad(test.quad, test.level)
		if !samePoint(out, test.out) {
			t.Errorf("[%d] got %v, want %v", i, out, test.out)
		}
	}
}

func Test_hilbert(t *testing.T) {
}
