package console

import (
	"fmt"

	"github.com/alexkappa/dc/data"
	"github.com/alexkappa/dc/flag"
	"github.com/alexkappa/dc/template"
)

// Console can render creatives to stdout.
type Console struct {
	template,
	data string
}

// Render renders a mustache template to stdout.
func (c *Console) Render() error {
	t, err := template.Parse(c.template)
	if err != nil {
		return err
	}
	d, err := data.Parse(c.data)
	if err != nil {
		return err
	}
	fmt.Println(t.Render(d))
	return nil
}

// New returns a new instance of Console configured by f.
func New(f flag.Flag) *Console {
	return &Console{f.Template, f.Data}
}
