package data

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

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
