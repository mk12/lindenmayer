// Copyright 2015 Mitchell Kember. Subject to the MIT License.

package main

import (
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/mk12/lindenmayer/system"
)

// parameters are a collection of accepted URL parameters.
type parameters struct {
	name      string
	depth     string
	thickness string
	color     string
	precision string
	onlySVG   bool
}

// pageData contains information used to render templates.
type pageData struct {
	Name       string
	Query      string
	Thickness  string
	Color      string
	Depth      int
	MaxDepth   int
	StepFactor float64
	PadFactor  float64
	SVG        template.HTML
	Systems    []string
}

// modifiedTime is the time this program was last changed in a way that affects
// the client. It should be updated to the current time to prevent browsers from
// using stale cached pages.
var modifiedTime = time.Date(2015, time.September, 8, 18, 5, 0, 0, time.UTC)

// systemNames contains the names of the systems shown in the sidebar.
var systemNames = []string{
	"koch", "hilbert", "peano", "gosper", "sierpinski", "rings", "tree",
	"plant", "willow", "dragon", "island",
}

// Limitations on parameters.
const (
	minimumDepth     = 0
	minimumPrecision = 1
	maximumPrecision = 15
)

// Default parameter values.
var defaultParams = parameters{
	name:      "koch",
	depth:     "2",
	thickness: "3",
	color:     "black",
	precision: "3",
	onlySVG:   false,
}

// Compile all templates on startup.
var templateSet *template.Template

func init() {
	add := func(a, b int) int { return a + b }
	capitalize := func(s string) string {
		runes := []rune(s)
		runes[0] = unicode.ToUpper(runes[0])
		return string(runes)
	}
	funcs := template.FuncMap{"add": add, "capitalize": capitalize}

	paths := []string{"header", "footer", "index", "404"}
	for i, name := range paths {
		paths[i] = "templates/" + name + ".html"
	}
	templateSet = template.Must(
		template.New("main").Funcs(funcs).ParseFiles(paths...))
}

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

// splitPath returns a slice of URL path segments by splitting on slashes. It
// also returns the extension, if present. The presence of a leading or trailing
// slash does not affect the result.
func splitPath(path string) (segs []string, ext string) {
	segs = strings.Split(path, "/")
	if len(segs) > 0 && segs[0] == "" {
		segs = segs[1:]
	}
	if s := len(segs); s > 0 && segs[s-1] == "" {
		segs = segs[:s-1]
	}
	if s := len(segs); s > 0 {
		last := segs[s-1]
		i := strings.LastIndexByte(last, '.')
		if i != -1 {
			if i < len(last)-1 {
				ext = last[i+1:]
			}
			if name := last[:i]; name != "" {
				segs[s-1] = name
			} else {
				segs = segs[:s-1]
			}
		}
	}
	return
}

// parseParams parses the parameters from the URL, and uses defaultParams if
// there are missing values.
func parseParams(href *url.URL) parameters {
	params := defaultParams

	path, ext := splitPath(href.Path)
	if len(path) >= 1 {
		params.name = path[0]
	}
	if len(path) >= 2 {
		params.depth = path[1]
	}
	if ext == "svg" {
		params.onlySVG = true
	}

	query := href.Query()
	if thickness := query.Get("t"); thickness != "" {
		params.thickness = thickness
	}
	if color := query.Get("c"); color != "" {
		params.color = color
	}
	if precision := query.Get("p"); precision != "" {
		params.precision = precision
	}

	return params
}

// renderSVG generates the SVG data for the specified curve.
func renderSVG(params parameters) (string, error) {
	sys := system.Named(params.name)
	if sys == nil {
		return "", fmt.Errorf("no system named %q", params.name)
	}

	depth, err := strconv.Atoi(params.depth)
	if err != nil {
		return "", err
	}
	if depth < minimumDepth || depth > sys.MaxDepth() {
		return "", fmt.Errorf("invalid depth %d", depth)
	}

	thickness, err := strconv.ParseFloat(params.thickness, 64)
	if err != nil {
		return "", err
	}
	if thickness <= 0 {
		return "", fmt.Errorf("invalid thickness %f", thickness)
	}

	if params.color == "" {
		return "", fmt.Errorf("invalid color \"\"")
	}

	precision, err := strconv.Atoi(params.precision)
	if err != nil {
		return "", err
	}
	if precision < minimumPrecision || precision > maximumPrecision {
		return "", fmt.Errorf("invalid precision %d", precision)
	}

	opts := system.Options{
		Depth:     depth,
		Thickness: thickness,
		Color:     params.color,
		Precision: precision,
	}
	return sys.SVG(&opts), nil
}

// mainHandler responds to an HTTP request.
func mainHandler(w http.ResponseWriter, req *http.Request) {
	log.SetPrefix(fmt.Sprintf("[%s %s] ", req.Method, req.URL))
	if req.Method != "GET" {
		fail(w, "invalid request method")
		return
	}

	ifModSince := req.Header.Get("If-Modified-Since")
	cacheTime, err := time.Parse(http.TimeFormat, ifModSince)
	if err == nil && modifiedTime.Unix() <= cacheTime.Unix() {
		w.WriteHeader(http.StatusNotModified)
		return
	}
	w.Header().Set("Last-Modified", modifiedTime.Format(http.TimeFormat))

	params := parseParams(req.URL)
	log.Printf("Rendering %+v\n", params)
	svg, err := renderSVG(params)
	if err != nil {
		fail(w, err.Error())
		return
	}

	if params.onlySVG {
		w.Header().Set("Content-Type", "image/svg+xml")
		io.WriteString(w, svg)
	} else {
		depth, _ := strconv.Atoi(params.depth)
		max := system.Named(params.name).MaxDepth()
		query := "?" + req.URL.RawQuery
		if query == "?" {
			query = ""
		}
		page := pageData{
			Name:       params.name,
			Thickness:  params.thickness,
			Color:      params.color,
			Query:      query,
			Depth:      depth,
			MaxDepth:   max,
			StepFactor: system.StepFactor,
			PadFactor:  system.PadFactor,
			SVG:        template.HTML(svg),
			Systems:    systemNames,
		}
		display(w, "index", page)
	}
}

func main() {
	log.SetFlags(0) // don't show timestamps in logs

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	port = ":" + port

	staticHandler := http.StripPrefix(
		"/static/",
		http.FileServer(http.Dir("static")),
	)
	http.HandleFunc("/", mainHandler)
	http.Handle("/static/", staticHandler)

	log.Println("=> Serving on http://localhost" + port)
	log.Println("=> Ctrl-C to shutdown server")
	err := http.ListenAndServe(port, nil)
	if err != nil {
		log.Fatalln("ListenAndServe:", err)
	}
}
