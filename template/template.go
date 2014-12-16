package template

import (
	"io/ioutil"
	"net/http"
	"os"

	"github.com/hoisie/mustache"
)

func Parse(s string) (*mustache.Template, error) {
	if f, err := os.Open(s); err == nil {
		defer f.Close()
		b, err := ioutil.ReadAll(f)
		if err != nil {
			return nil, err
		}
		return mustache.ParseString(string(b))
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

func Render(s string, c interface{}) string {
	return mustache.Render(s, c)
}

const PreviewHTML = `<!DOCTYPE html>
<html>
<head>
	<title>Yieldr Creative Preview</title>
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
    {{#Frames}}
    	<h2>{{Name}}</h2>
		<iframe src="{{URL}}" width="{{Width}}px" height="{{Height}}px" scrolling="no" frameborder="0">
		</iframe>
    {{/Frames}}
</body>
</html>`
