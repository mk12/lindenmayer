package main

import (
	"reflect"
	"testing"
)

func Test_pathSegments(t *testing.T) {
	table := []struct {
		path string
		segs []string
	}{
		{"", []string{}},
		{"/", []string{}},
		{"//", []string{""}},
		{"foo/bar/baz", []string{"foo", "bar", "baz"}},
		{"/foo/bar/baz", []string{"foo", "bar", "baz"}},
		{"foo/bar/baz/", []string{"foo", "bar", "baz"}},
		{"/foo/bar/baz/", []string{"foo", "bar", "baz"}},
	}
	for i, test := range table {
		result := pathSegments(test.path)
		if !reflect.DeepEqual(result, test.segs) {
			t.Errorf("[%d] got %#v, want %#v", i, result, test.segs)
		}
	}
}
