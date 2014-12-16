package data

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

// Parse attempts to parse s in several ways. First it attempts to open a file
// locally with the name s. If that fails then it tries to get s over http. If
// that fails it treats s as input and attempts to parse it.
func Parse(s string) (interface{}, error) {
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
