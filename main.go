package main

import (
	"flag"
	"fmt"
	"github.com/ajstarks/svgo"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"
)

// A Renderer is a function that renders a space-filling curve on an SVG object
// to a certain depth of recursion.
type Renderer func(depth int, s *svg.SVG)

// allRenderes maps from URL handles to Renderer functions.
var allRenderers = map[string]Renderer{
	"hilbert": renderHilbert,
	"peano":   renderPeano,
}

// Default parameter values.
var (
	defaultRender = renderHilbert
	defaultDepth  = 1
)

// Compile templates on startup.
var templateSet = template.Must(template.ParseFiles(
	"templates/header.html",
	"templates/footer.html",
	"templates/index.html",
	"templates/404.html",
))

// display executes the template with the given name, and passes along data.
func display(w http.ResponseWriter, templateName string, data interface{}) {
	err := templateSet.ExecuteTemplate(w, templateName, data)
	if err != nil {
		log.Println("500 Internal Server Error:", err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

// pathSegments returns a slice of URL path segments by splitting on slashes.
// The presence of a leading or trailing slash does not matter.
func pathSegments(path string) []string {
	segs := strings.Split(path, "/")
	if len(segs) > 0 && segs[0] == "" {
		segs = segs[1:]
	}
	if s := len(segs); s > 0 && segs[s-1] == "" {
		segs = segs[:s-1]
	}
	return segs
}

// notFound responds with a 404 Not Found status after logging a message.
func notFound(w http.ResponseWriter, reason string) {
	log.Printf("404 Not Found (%s)\n", reason)
	w.WriteHeader(http.StatusNotFound)
	display(w, "404", nil)
}

// mainHandler responds to an HTTP request.
func mainHandler(w http.ResponseWriter, req *http.Request) {
	if req.Method != "GET" {
		notFound(w, "invalid request method")
		return
	}

	log.SetPrefix(fmt.Sprintf("[GET %s] ", req.URL.Path))
	segments := pathSegments(req.URL.Path)
	render := defaultRender
	depth := defaultDepth

	if len(segments) >= 1 {
		if r, ok := allRenderers[segments[0]]; ok {
			render = r
		} else {
			log.Printf("segs: %#v\n", segments)
			notFound(w, "invalid renderer type")
			return
		}
	}
	if len(segments) >= 2 {
		if d, err := strconv.Atoi(segments[1]); err == nil && d >= 0 {
			depth = d
		} else {
			notFound(w, "invalid depth")
			return
		}
	}

	w.Header().Set("Content-Type", "image/svg+xml")
	s := svg.New(w)
	s.Start(600, 600)
	render(depth, s)
	s.End()
}

func main() {
	log.SetFlags(0) // don't show timestamps in logs

	port := flag.Int("port", 8080, "localhost port to serve from")
	flag.Parse()
	portStr := fmt.Sprintf(":%d", *port)

	staticHandler := http.FileServer(http.Dir("static"))
	staticHandler = http.StripPrefix("/static/", staticHandler)
	http.HandleFunc("/", mainHandler)
	http.Handle("/static/", staticHandler)

	log.Println("=> Serving on http://localhost" + portStr)
	log.Println("=> Ctrl-C to shutdown server")
	err := http.ListenAndServe(portStr, nil)
	if err != nil {
		log.Fatalln("ListenAndServe:", err)
	}
}
