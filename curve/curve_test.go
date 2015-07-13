// Copyright 2015 Mitchell Kember. Subject to the MIT License.

package curve

import (
	"strings"
	"testing"
)

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
			t.Errorf("[%d] want error, got nil", i)
		}
	}

	svg, err := SVG("hilbert", "1")
	if err != nil {
		t.Errorf("unexpected error: %q", err)
	}
	if !strings.Contains(svg, "<svg") {
		t.Errorf("got %q, want <svg> data", svg)
	}
}

func Test_connectDots(t *testing.T) {
	svg := connectDots([]point{})
	if svg != "" {
		t.Errorf("want \"\", got %q", svg)
	}

	list := [][]point{
		{{0, 0}},
		{{0.1, 0.1}, {0.5, 0.5}},
		{{1, -2}, {-3, 4}, {5, -6}},
	}
	for i, dots := range list {
		svg := connectDots(dots)
		if !strings.Contains(svg, "<svg") {
			t.Errorf("[%d] got %q, want <svg> data", i, svg)
		}
	}
}
