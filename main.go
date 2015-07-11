package main

import (
	"flag"
	"github.com/ajstarks/svgo"
	"log"
	"strconv"
	"strings"
)

// A Renderer is a function that renders a space-filling curve on an SVG object
// to a certain depth of recursion.
type Renderer func(depth int, s *svg.SVG)

var allRenderers = map[string]Renderer{
	"hilbert": renderHilbert,
	"peano":   renderPeano,
}

// Default parameter values.
const (
	defaultRender = renderHilbert
	defaultDepth  = 1
)

// Dimensions of the SVG element.
const (
	WIDTH  = 600
	HEIGHT = 600
)

// pathSegments returns a slice of URL path segments by splitting on slashes.
// The presence of a leading or trailing slash does not matter.
func pathSegments(path string) []string {
	path = strings.trimPrefix(path, "/")
	path = strings.trimSuffix(path, "/")
	return strings.Split(path, "/")
}

// application responds to an HTTP request.
func application(w http.ResponseWriter, req *http.Request) {
	if req.Method != "GET" {
		return notFound(w)
	}

	log.SetPrefix(fmt.Sprintf("[GET %s] ", req.URL.Path))
	segments := pathSegments(req.URL.Path)
	render := defaultRender
	depth := defaultDepth

	if len(segments) >= 1 {
		if r, ok := allRenderers[segments[0]]; ok {
			render = r
		} else {
			return notFound(w)
		}
	}
	if len(segments) >= 2 {
		if d, err := strconv.Atoi(segments[1]); err == nil && d >= 0 {
			depth = d
		} else {
			return notFound(w)
		}
	}

	w.Header().Set("Content-Type", "image/svg+xml")
	s := svg.New(w)
	s.Start(WIDTH, HEIGHT)
	render(depth, s)
	s.End()
}

func main() {
	log.SetFlags(0)

	port := flag.Int("port", 8080, "localhost port to serve from")
	flag.Parse()
	portStr := fmt.Sprintf(":%d", *port)

	log.Println("=> Serving on http://localhost" + portStr)
	log.Println("=> Ctrl-C to shutdown server")
	http.HandleFunc("/", application)
	err := http.ListenAndServe(portStr, nil)
	if err != nil {
		log.Fatalln("ListenAndServe:", err)
	}
}
