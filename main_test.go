// Copyright 2015 Mitchell Kember. Subject to the MIT License.

package main

import (
	"bytes"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
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

func Test_curveOptions(t *testing.T) {
	table := []struct {
		path  string
		name  string
		depth string
	}{
		{"", "hilbert", "1"},
		{"/", "hilbert", "1"},
		{"//", "", "1"},
		{"///", "", ""},
		{"/hilbert", "hilbert", "1"},
		{"/hilbert/2", "hilbert", "2"},
		{"/peano/", "peano", "1"},
		{"/peano/12/", "peano", "12"},
	}
	for i, test := range table {
		name, depth := curveOptions(test.path)
		if name != test.name || depth != test.depth {
			t.Errorf("[%d] got (%q, %q), want (%q, %q)",
				i, name, depth, test.name, test.depth)
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
	if !bytes.Contains(body, []byte("Space-filling Curves")) {
		t.Error("index page should include title 'Space-filling Curves'")
	}
	if !bytes.Contains(buf.Bytes(), []byte("hilbert")) {
		t.Error("should log about rendering hilbert curve")
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
