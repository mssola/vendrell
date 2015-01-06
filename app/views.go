// Copyright (C) 2014-2015 Miquel Sabaté Solà <mikisabate@gmail.com>
// This file is licensed under the MIT license.
// See the LICENSE file.

package app

import (
	"bytes"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"path"
	"time"
)

var (
	// The directory where all the views are being stored.
	viewsDir = "views"
)

const (
	// The path to the layout file.
	layout = "application/layout"

	// The extension of views.
	viewsExt = "tpl"
)

// This struct holds all the data that can be passed to a view.
type Options struct {
	// The id of the current user.
	Id string

	// The current player.
	One *ExtPlayer

	// All the players to be displayed.
	Players []*ExtPlayer

	// The maximum number of values to be displayed for a set of players.
	Values []int

	// The URL of "Download CSV". It might change depending where we are.
	Download string

	// Set to true if the current user is logged in.
	LoggedIn bool

	// Set to true if the views has to include Javascript.
	JS bool

	// Set to true if an error has happenned.
	Error bool
}

// Returns the path to be used to open the view with the given name.
func view(name string) string {
	return path.Join(viewsDir, name+"."+viewsExt)
}

// Render the view with the given name after evaluating the passed data. The
// rendered view will be written to the given writer.
func render(res http.ResponseWriter, name string, data interface{}) {
	b, e := ioutil.ReadFile(view(layout))
	if e != nil {
		panic("Could not read layout file!")
	}
	t, e := template.New("l").Funcs(layoutHelpers(name, data)).Parse(string(b))
	if e != nil {
		panic("Could not parse layout file!")
	}
	t.Execute(res, data)
}

// Returns all the helpers used by the layout template. Right now only the
// "yield" helpers has been implemented.
func layoutHelpers(name string, data interface{}) template.FuncMap {
	return template.FuncMap{
		"yield": func() template.HTML {
			var buffer bytes.Buffer

			b, e := ioutil.ReadFile(view(name))
			if e != nil {
				r := fmt.Sprintf("Could not read: %v => %v", name, e)
				panic(r)
			}
			t := template.New(name).Funcs(viewHelpers())
			t, e = t.Parse(string(b))
			if e != nil {
				r := fmt.Sprintf("Could not parse: %v => %v", name, e)
				panic(r)
			}
			t.Execute(&buffer, data)
			return template.HTML(buffer.String())
		},
	}
}

// Returns all the helpers available to any view. We have the following
// helpers: fmtDate and inc. The inc helper just increases the given integer
// value by one. The fmtDate helper executes the fmtDate function.
func viewHelpers() template.FuncMap {
	return template.FuncMap{
		"fmtDate": fmtDate,
		"inc": func(n int) int {
			return n + 1
		},
	}
}

// Returns a string with the given time formatted as expected by the view.
func fmtDate(t time.Time) string {
	return fmt.Sprintf("%02d/%02d/%04d", t.Day(), t.Month(), t.Year())
}
