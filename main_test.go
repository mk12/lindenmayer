// Copyright 2015 Mitchell Kember. Subject to the MIT License.

package main

import (
	"bytes"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"testing"
)

func Test_display200(t *testing.T) {
	var buf bytes.Buffer
	log.SetOutput(&buf)
	w := httptest.NewRecorder()

	display(w, "index", "HELLO")
	if w.Code != http.StatusOK {
		t.Errorf("got %d, want http.StatusOK", w.Code)
	}
	if !bytes.Contains(w.Body.Bytes(), []byte("HELLO")) {
		t.Error("index page should include string 'HELLO'")
	}
	if buf.Len() != 0 {
		t.Errorf("should not log anything: %q", buf.String())
	}
}

func Test_display500(t *testing.T) {
	var buf bytes.Buffer
	log.SetOutput(&buf)
	w := httptest.NewRecorder()

	display(w, "nonexistent", nil)
	if w.Code != http.StatusInternalServerError {
		t.Errorf("got %d, want http.StatusInternalServerError", w.Code)
	}
	if w.Body.Len() != 0 {
		t.Error("should not render a body: %q", w.Body.String())
	}
	if !bytes.Contains(buf.Bytes(), []byte("500")) {
		t.Error("should have logged the 500")
	}
	if !bytes.Contains(buf.Bytes(), []byte("nonexistent")) {
		t.Error("should complain about the 'nonexistent' template")
	}
}

func Test_fail(t *testing.T) {
	var buf bytes.Buffer
	log.SetOutput(&buf)
	w := httptest.NewRecorder()

	fail(w, "REASON")
	if w.Code != http.StatusNotFound {
		t.Errorf("got %d, want http.StatusNotFound", w.Code)
	}
	if !bytes.Contains(w.Body.Bytes(), []byte("404")) {
		t.Error("404 page should include string '404'")
	}
	if !bytes.Contains(buf.Bytes(), []byte("REASON")) {
		t.Error("should have logged the failure reason")
	}
}

func Test_splitPath(t *testing.T) {
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
		result := splitPath(test.path)
		if !reflect.DeepEqual(result, test.segs) {
			t.Errorf("[%d] got %#v, want %#v", i, result, test.segs)
		}
	}
}

func Test_parseParams(t *testing.T) {
	d := defaultParams
	table := []struct {
		rawurl string
		params parameters
	}{
		{"/", d},
		{"/abc", parameters{
			name:      "abc",
			depth:     d.depth,
			thickness: d.thickness,
			color:     d.color,
			precision: d.precision,
		}},
		{"/abc/3", parameters{
			name:      "abc",
			depth:     "3",
			thickness: d.thickness,
			color:     d.color,
			precision: d.precision,
		}},
		{"/?t=4&c=red&p=9", parameters{
			name:      d.name,
			depth:     d.depth,
			thickness: "4",
			color:     "red",
			precision: "9",
		}},
	}
	for i, test := range table {
		href, err := url.Parse(test.rawurl)
		if err != nil {
			t.Errorf("[%d] could not parse %q", i, test.rawurl)
			continue
		}

		params := parseParams(href)
		if !reflect.DeepEqual(params, test.params) {
			t.Errorf("[%d] got %#v, want %#v", i, params, test.params)
		}
	}
}

func Test_renderSVG(t *testing.T) {
	table := []struct {
		params parameters
		good   bool
	}{
		{defaultParams, true},
		{parameters{"koch", "1", "1", "black", "1"}, true},
		{parameters{"", "1", "1", "black", "1"}, false},
		{parameters{"koch", "", "1", "black", "1"}, false},
		{parameters{"koch", "1", "1", "", "1"}, false},
		{parameters{"koch", "1", "", "black", "1"}, false},
		{parameters{"koch", "1", "1", "black", ""}, false},
		{parameters{"!@#$", "1", "1", "black", "1"}, false},
		{parameters{"koch", "-1", "1", "black", "1"}, false},
		{parameters{"koch", "1", "-1", "red", "1"}, false},
		{parameters{"koch", "1", "1", "black", "-1"}, false},
	}
	for i, test := range table {
		svg, err := renderSVG(test.params)
		if test.good {
			if err != nil {
				t.Errorf("[%d] unexpected error: %q", i, err)
			} else if svg == "" {
				t.Errorf("[%d] got %q, want SVG data", i, svg)
			}
		} else {
			if err == nil {
				t.Errorf("[%d] expected error", i)
			}
		}
	}
}

func Test_mainHandlerGet(t *testing.T) {
	var buf bytes.Buffer
	log.SetOutput(&buf)
	server := httptest.NewServer(http.HandlerFunc(mainHandler))
	defer server.Close()

	resp, err := http.Get(server.URL)
	if err != nil {
		t.Errorf("unexpected error: %q", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("unexpected error: %q", err)
	}
	if !bytes.Contains(body, []byte("Lindenmayer")) {
		t.Error("index page should include title 'Lindenmayer'")
	}
	if !bytes.Contains(buf.Bytes(), []byte("koch")) {
		t.Error("should log about rendering koch curve")
	}
}

func Test_mainHandlerPost(t *testing.T) {
	var buf bytes.Buffer
	log.SetOutput(&buf)
	server := httptest.NewServer(http.HandlerFunc(mainHandler))
	defer server.Close()

	resp, err := http.Post(server.URL, "application/json", nil)
	if err != nil {
		t.Errorf("unexpected error: %q", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("unexpected error: %q", err)
	}
	if !bytes.Contains(body, []byte("404 Not Found")) {
		t.Error("should render a 404 page")
	}
	if !bytes.Contains(buf.Bytes(), []byte("POST")) {
		t.Error("should complain about POST request")
	}
}
