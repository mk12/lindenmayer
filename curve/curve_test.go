// Copyright 2015 Mitchell Kember. Subject to the MIT License.

package curve

import (
	"math"
	"testing"
)

const epsilon = 1e-10

func samePoint(a, b point) bool {
	return math.Abs(a.x-b.x) < epsilon && math.Abs(a.y-b.y) < epsilon
}

func samePoints(s1, s2 []point) bool {
	if len(s1) != len(s2) {
		return false
	}
	for i := 0; i < len(s1); i++ {
		if !samePoint(s1[i], s2[i]) {
			return false
		}
	}
	return true
}

func TestSVG(t *testing.T) {
	table := []struct {
		name, depth string
	}{
		{"", ""},
		{"hilbert", ""},
		{"", "1"},
		{"nonexistent", "1"},
		{"hilbert", "-1"},
		{"peano", "a"},
	}
	for i, test := range table {
		if _, err := SVG(test.name, test.depth); err == nil {
			t.Errorf("[%d] got nil, want error", i)
		}
	}

	svg, err := SVG("hilbert", "1")
	if err != nil {
		t.Errorf("unexpected error: %q", err)
	}
	if svg[:4] != "<svg" {
		t.Errorf("got %q, want <svg> data", svg)
	}
}

func Test_connectDots(t *testing.T) {
	svg := connectDots([]point{})
	if svg != "" {
		t.Errorf("got %q, want \"\"", svg)
	}

	list := [][]point{
		{{0, 0}},
		{{0.1, 0.1}, {0.5, 0.5}},
		{{1, -2}, {-3, 4}, {5, -6}},
	}
	for i, dots := range list {
		svg := connectDots(dots)
		if svg[:4] != "<svg" {
			t.Errorf("[%d] got %q, want <svg> data", i, svg)
		}
	}
}

func Test_removeLines(t *testing.T) {
	table := []struct {
		n   int
		in  string
		out string
	}{
		{0, "", ""},
		{0, "abc", "abc"},
		{0, "abc\nxyz\n", "abc\nxyz\n"},
		{1, "abc\nxyz", "xyz"},
		{100, "", ""},
		{2, "a\nb\nc\nd\n", "c\nd\n"},
	}
	for i, test := range table {
		out := removeLines(test.in, test.n)
		if out != test.out {
			t.Errorf("[%d] got %q, want %q", i, out, test.out)
		}
	}
}
