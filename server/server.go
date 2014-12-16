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

func (s *Server) URL() string {
	if s.isDir {
		return "http://" + s.addr + "/p"
	}
	return "http://" + s.addr + "/t"
}

func (s *Server) Open() error {
	return open.Run(s.URL())
}

func (s *Server) Listen() error {
	return http.ListenAndServe(s.addr, s)
}

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

type Context struct {
	Frames []Frame
}

type Frame struct {
	Name string
}

func (f Frame) URL() string {
	return "/d/" + f.Name
}

func (f Frame) Width() string {
	return strings.Split(f.Name[:strings.Index(f.Name, ".")], "x")[0]
}

func (f Frame) Height() string {
	return strings.Split(f.Name[:strings.Index(f.Name, ".")], "x")[1]
}
