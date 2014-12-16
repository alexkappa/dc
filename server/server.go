package server

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/alexkappa/dc/data"
	"github.com/alexkappa/dc/fileutil"
	"github.com/alexkappa/dc/flag"
	"github.com/alexkappa/dc/template"

	"github.com/gorilla/mux"
	"github.com/skratchdot/open-golang/open"
)

// Server is able to listen to http connections and render mustache templates.
type Server struct {
	template,
	data,
	static,
	addr string
	isDir bool

	*mux.Router
}

// ServeTemplate renders a single mustache template as defined by command line
// arguments using the -t option.
func (s *Server) ServeTemplate(w http.ResponseWriter, r *http.Request) {
	t, err := template.Parse(filepath.Join(s.template))
	if err != nil {
		fmt.Fprintf(w, "Unable to parse template. %s", err)
		return
	}
	d, err := data.Parse(s.data)
	if err != nil {
		fmt.Fprintf(w, "Unable to parse data. %s", err)
		return
	}
	w.Header().Add("Content-Type", "text/html")
	fmt.Fprintf(w, "%s", t.Render(d))
}

// ServeDirectory renders a single mustache template which is to be found inside
// a specific directory. This directory is specified using command line argumens
// using the -t option.
//
// Usually you wouldn't navigate to this page yourself, rather this endpoint
// will be called from the preview page.
func (s *Server) ServeDirectory(w http.ResponseWriter, r *http.Request) {
	name := mux.Vars(r)["name"]
	if !s.isDir {
		fmt.Fprintf(w, "This endpoint is not allowed")
	}
	t, err := template.Parse(filepath.Join(s.template, name))
	if err != nil {
		fmt.Fprintf(w, "Unable to parse template. %s", err)
		return
	}
	d, err := data.Parse(s.data)
	if err != nil {
		fmt.Fprintf(w, "Unable to parse data. %s", err)
		return
	}
	w.Header().Add("Content-Type", "text/html")
	fmt.Fprintf(w, "%s", t.Render(d))
}

// ServePreview renders the preview page with all templates inside the template
// directory as iframes. The template directory is specified using command line
// arguments with the -t option.
func (s *Server) ServePreview(w http.ResponseWriter, r *http.Request) {
	// render a preview html which contains several iframes, each one with the
	// name of every template under the directory we're in.
	files, err := ioutil.ReadDir(s.template)
	if err != nil {
		fmt.Fprintf(w, "Unable to read template directory. %s", err)
		return
	}
	context := Context{
		Frames: make([]Frame, len(files)),
	}
	for i, file := range files {
		context.Frames[i].Name = file.Name()
	}
	fmt.Fprintf(w, "%s", template.Render(template.PreviewHTML, context))
}

// URL returns the address which the server listens on for HTTP connections.
func (s *Server) URL() string {
	return "http://" + s.addr
}

// Open opens a new browser window at the servers URL. The path depends on
// whether the template specified is a
func (s *Server) Open() error {
	if s.isDir {
		return open.Run(s.URL() + "/p")
	}
	return open.Run(s.URL() + "/t")

}

// Listen listens for HTTP requests on the address specified by the -a command
// line option.
func (s *Server) Listen() error {
	return http.ListenAndServe(s.addr, s)
}

// New creates a new Server instance configured with the supplied command line
// arguments.
func New(f flag.Flag) *Server {
	s := &Server{
		template: f.Template,
		data:     f.Data,
		static:   f.Static,
		addr:     f.Addr,
		isDir:    fileutil.IsDir(f.Template),
	}

	r := mux.NewRouter()
	r.HandleFunc("/t", s.ServeTemplate)         // single template rendering
	r.HandleFunc("/d/{name}", s.ServeDirectory) // template directory
	r.HandleFunc("/p", s.ServePreview)          // preview

	r.NotFoundHandler = http.FileServer(http.Dir(fileutil.Abs(s.static)))

	s.Router = r

	return s
}

// Context is the data structure that populates the preview template.
type Context struct {
	Frames []Frame
}

// Frame is used within Context and it represents a single <iframe> element.
type Frame struct {
	Name string
}

// URL returns the address for rendering the particular iframe.
func (f Frame) URL() string {
	return "/d/" + f.Name
}

// Width returns the width of the iframe based on f.Name.
func (f Frame) Width() string {
	return strings.Split(f.Name[:strings.Index(f.Name, ".")], "x")[0]
}

// Height returns the height of the iframe based on f.Name.
func (f Frame) Height() string {
	return strings.Split(f.Name[:strings.Index(f.Name, ".")], "x")[1]
}
