package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

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
	addr string
	serve bool
)

func init() {
	flag.StringVar(&html, "t", "", "The creative template.")
	flag.StringVar(&data, "d", "", "The data to populate the template with.")
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
		go open.Run("http://" + addr)
		http.ListenAndServe(addr, handle(html, data))
	} else {
		t, err := parseTemplate(html)
		if err != nil {
			fmt.Fprintf(os.Stdout, "Unable to parse template. %s\n", err)
			os.Exit(1)
		}
		d, err := parseData(data)
		if err != nil {
			fmt.Fprintf(os.Stdout, "Unable to parse data. %s\n", err)
			os.Exit(1)
		}
		fmt.Println(t.Render(d))
	}
}

func parseTemplate(s string) (*mustache.Template, error) {
	if f, err := os.Open(s); err == nil {
		f.Close()
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

func handle(html, data string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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
		fmt.Fprintf(w, "%s", t.Render(d))
	}
}
