package main

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/hoisie/mustache"
	"github.com/skratchdot/open-golang/open"
)

const (
	File = iota
	URL
	Text
)

var (
	html,
	data,
	addr,
	static string
	serve bool
)

func init() {
	flag.StringVar(&html, "t", "", "The creative template.")
	flag.StringVar(&data, "d", "", "The data to populate the template with.")
	flag.StringVar(&static, "r", "", "Static file root")
	flag.BoolVar(&serve, "s", false, "Serve the template using HTTP.")
	flag.StringVar(&addr, "h", "localhost:8080", "Server address.")
	flag.Parse()
}

func main() {
	if html == "" || data == "" {
		fmt.Fprintf(os.Stderr, "Usage %s:\n", os.Args[0])
		flag.PrintDefaults()
		os.Exit(1)
	}
	if serve {
		go open.Run("http://" + addr + "/render")
		http.ListenAndServe(addr, handle(html, data, static))
	} else {
		t, err := parseTemplate(html)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Unable to parse template. %s\n", err)
			os.Exit(1)
		}
		d, err := parseData(data)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Unable to parse data. %s\n", err)
			os.Exit(1)
		}
		fmt.Println(t.Render(d))
	}
}

func parseTemplate(s string) (*mustache.Template, error) {
	if f, err := os.Open(s); err == nil {
		defer f.Close()
		finfo, err := f.Stat()
		if err != nil {
			return nil, err
		}
		if finfo.IsDir() {
			return parseTemplateDir(s)
		}
		return mustache.ParseFile(s)
	} else if r, err := http.Get(s); err == nil {
		b, err := ioutil.ReadAll(r.Body)
		if err != nil {
			return nil, err
		}
		r.Body.Close()
		return mustache.ParseString(string(b))
	}
	return mustache.ParseString(s)
}

func parseTemplateDir(dir string) (*mustache.Template, error) {
	d, err := os.Open(dir)
	if err != nil {
		fmt.Println(dir, err)
		return nil, err
	}
	defer d.Close()
	files, err := d.Readdirnames(0)
	if err != nil {
		return nil, err
	}
	var t, p string
	for _, file := range files {
		f, err := os.Open(filepath.Join(dir, file))
		if err != nil {
			return nil, err
		}
		b, err := ioutil.ReadAll(f)
		if err != nil {
			return nil, err
		}
		dim := strings.Split(file[:strings.Index(file, ".")], "x")
		if len(dim) != 2 {
			return nil, errors.New("file names should be in the form {width}x{height}.html")
		}
		r := strings.NewReplacer(
			"{{name}}", file,
			"{{src}}", "data:text/html;base64,"+base64.StdEncoding.EncodeToString(b),
			"{{width}}", dim[0],
			"{{height}}", dim[1])
		p += r.Replace(frame)
	}
	t = strings.Replace(layout, "{{frames}}", p, 1)
	return mustache.ParseString(t)
}

func parseData(s string) (v interface{}, err error) {
	if f, err := os.Open(s); err == nil {
		b, err := ioutil.ReadAll(f)
		if err != nil {
			return nil, err
		}
		defer f.Close()
		return unmarshal(b)
	} else if r, err := http.Get(s); err == nil {
		b, err := ioutil.ReadAll(r.Body)
		if err != nil {
			return nil, err
		}
		r.Body.Close()
		return unmarshal(b)
	}
	return unmarshal([]byte(s))
}

func unmarshal(b []byte) (v interface{}, err error) {
	if err = json.Unmarshal(b, &v); err == nil {
		return
	}
	return nil, fmt.Errorf("Invalid format. %s", err)
}

func handle(html, data, static string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, "/render") {
			t, err := parseTemplate(html)
			if err != nil {
				fmt.Fprintf(w, "Unable to parse template. %s\n", err)
				return
			}
			d, err := parseData(data)
			if err != nil {
				fmt.Fprintf(w, "Unable to parse data. %s\n", err)
				return
			}
			w.Header().Add("Content-Type", "text/html")
			fmt.Fprintf(w, "%s", t.Render(d))
		} else {
			if !filepath.IsAbs(static) {
				wd, err := os.Getwd()
				if err != nil {
					fmt.Fprintf(w, "Unable to get the current working directory: %s", err)
				}
				static = filepath.Join(wd, static)
			}
			http.FileServer(http.Dir(static)).ServeHTTP(w, r)
		}
	}
}

const layout = `<!DOCTYPE html>
<html>
<head>
    <style>
        * {
            font-family: "Avant Garde", Avantgarde, "Century Gothic", CenturyGothic, "AppleGothic", sans-serif;
        }
        h2 {
            margin: 5px 0;
        }
        .note {
            display: block;
            text-align: center;
            background: #f0f0f0;
            line-height: 36px;
        }
        div {
            float: left;
            margin: 15px;
        }
        .long {
            height: 200px;
        }
        .square {
            width: 50%;
            float: left;
            margin-top: -5%;
        }
    </style>
</head>
<body>
    <p class="note">Please clean cache browser before reviewing it. Reload holding shift key</p>
    {{frames}}
</body>
</html>`

const frame = `<div>
	<h2>{{name}}</h2>
	<iframe src="{{src}}" width="{{width}}px" height="{{height}}px" scrolling="no" frameborder="0">
	</iframe>
</div>`
