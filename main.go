// Copyright 2015 Mitchell Kember. Subject to the MIT License.

package main

import (
	"flag"
	"fmt"
	"github.com/mk12/sfc/curve"
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"
)

// Default curve options.
const (
	defaultCurveName  = "hilbert"
	defaultCurveDepth = "1"
)

// Compile all templates on startup.
var templateSet = template.Must(template.ParseFiles(
	"templates/header.html",
	"templates/footer.html",
	"templates/index.html",
	"templates/404.html",
))

// display finds a template by name and executes it with the given data.
func display(w http.ResponseWriter, templateName string, data interface{}) {
	err := templateSet.ExecuteTemplate(w, templateName, data)
	if err != nil {
		log.Println("500 Internal Server Error:", err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

// fail logs a failure message and responds with a 404 Not Found status.
func fail(w http.ResponseWriter, reason string) {
	log.Println("404 Not Found:", reason)
	w.WriteHeader(http.StatusNotFound)
	display(w, "404", nil)
}

// splitPath returns a slice of URL path segments by splitting on slashes. The
// presence of a leading or trailing slash does not affect the result.
func splitPath(path string) []string {
	segs := strings.Split(path, "/")
	if len(segs) > 0 && segs[0] == "" {
		segs = segs[1:]
	}
	if s := len(segs); s > 0 && segs[s-1] == "" {
		segs = segs[:s-1]
	}
	return segs
}

// curveOptions returns the options (name and depth) for the curve as specified
// in the path. Returns default values if one or both are missing.
func curveOptions(path string) (name string, depth string) {
	args := splitPath(path)
	if len(args) >= 1 {
		name = args[0]
	} else {
		name = defaultCurveName
	}
	if len(args) >= 2 {
		depth = args[1]
	} else {
		depth = defaultCurveDepth
	}
	return
}

// mainHandler responds to an HTTP request. If there are no errors, it renders
// the index page with the desired space-filling curve embedded as an SVG.
func mainHandler(w http.ResponseWriter, req *http.Request) {
	log.SetPrefix(fmt.Sprintf("[%s %s] ", req.Method, req.URL.Path))
	if req.Method != "GET" {
		fail(w, "invalid request method")
		return
	}

	name, depth := curveOptions(req.URL.Path)
	svg, err := curve.SVG(name, depth)
	if err != nil {
		fail(w, err.Error())
		return
	}
	log.Printf("Rendering curve '%s' at depth %s\n", name, depth)
	display(w, "index", template.HTML(svg))
}

func main() {
	log.SetFlags(0) // don't show timestamps in logs

	port := flag.Int("port", 8080, "localhost port to serve from")
	flag.Parse()
	if flag.NArg() > 0 {
		flag.Usage()
		os.Exit(1)
	}

	staticHandler := http.StripPrefix(
		"/static/",
		http.FileServer(http.Dir("static")),
	)
	http.HandleFunc("/", mainHandler)
	http.Handle("/static/", staticHandler)

	portStr := fmt.Sprintf(":%d", *port)
	log.Println("=> Serving on http://localhost" + portStr)
	log.Println("=> Ctrl-C to shutdown server")
	err := http.ListenAndServe(portStr, nil)
	if err != nil {
		log.Fatalln("ListenAndServe:", err)
	}
}
