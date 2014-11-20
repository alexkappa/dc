# dc

The Yieldr Dynamic Creative rendering tool is intended for use by developers implementing dynamic creatives and wish to test using product information.

# Usage

Call `dc` with a template `-t` and some data `-d`. Both the template and the data arguments can be one of the following forms:

- File `-t example.html.mustache` `-d data.json`
- URL: `-t 'http://example.com/example.html.mustache'` `-d 'http://jsonip.com'`
- String: `-t '<p>{{foo}}</p>'` `-d '{"foo": "bar"}'`

```bash
$ ./dc -t '<div>{{#dicks}}<p>{{.}}</p>{{/dicks}}</div>' -d 'http://dicks-api.herokuapp.com/dicks/5'
<div><p>8==D</p><p>8======D</p><p>8====D</p><p>8===D</p><p>8==D</p></div>
```

