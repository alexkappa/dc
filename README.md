# dc

Quickly render mustache templates

# Usage

Call `dc` with a template `-t` and some data `-d`. Both the template and the data arguments can be one of the following forms:

- File `-t example.html.mustache` `-d data.json`
- URL: `-t 'http://example.com/example.html.mustache'` `-d 'http://jsonip.com'`
- String: `-t '<p>{{foo}}</p>'` `-d '{"foo": "bar"}'`

## Examples

```bash
$ ./dc -t '<div>{{#dicks}}<p>{{.}}</p>{{/dicks}}</div>' -d 'http://dicks-api.herokuapp.com/dicks/5'
```

```bash
$ ./dc -t template.html.mustache -d products.json
```

A useful feature of `dc` is the built in web server. Using the `-s` option `dc` will start a web server listening on `http://localhost:8080` and opens up a browser window with your rendered creative. From here on out you can freely refresh the page and the template will be re-rendered.

```bash
$ ./dc -s -t template.html.mustache -d products.json
```

Access with recommender data.

```bash
$ cat template.html.mustache
<ul>
    {{#products}}
    <li>
        <a href="{{deeplink}}">{{from_airport}}-{{to_airport}}</a>
    </li>
    {{/products}}
</ul>

$ cat products.json
{
  "products": [
    {
      "from_airport": "TRN"
      "to_airport": "PMO",
      "deeplink": "/Search.aspx?culture=it-IT&from=TRN&to=PMO&departuredate=2014-08-20&returndate=2014-08-22",
    }
  ]
}

$ dc -t template.html.mustache -d products.json
<ul>
    <li>
        <a href="/Search.aspx?culture=it-IT&amp;bookingtype=flight&amp;from=TRN&amp;to=PMO&amp;triptype=roundtrip&amp;departuredate=2014-08-20&amp;returndate=2014-08-22&amp;adults=1&amp;children=0&amp;infants=0&amp;showNewSearch=false&amp;utm_campaign=retargeting">TRN-PMO</a>
    </li>
</ul>
```

# Mustache

You can read the documentation on mustache logic-less templating [here](http://mustache.github.io/) or play around with their [demo](http://mustache.github.io/#demo).