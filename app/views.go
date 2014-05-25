// Copyright (C) 2014 Miquel Sabaté Solà <mikisabate@gmail.com>
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

// TODO: we can do something nicer with these structs...

type Options struct {
	Id       string
	Name     string
	LoggedIn bool
	Error    bool
	JS       bool
	Stats    *Statistics
}

type ExtendedHome struct {
	Players  []*ExtendedPlayer
	Values   []int
	LoggedIn bool
	JS       bool
}

type Home struct {
	Players  []Player
	LoggedIn bool
	JS       bool
}

const (
	layout   = "application/layout"
	viewsDir = "views"
	viewsExt = "tpl"
)

func view(name string) string {
	return path.Join(viewsDir, name+"."+viewsExt)
}

func render(res http.ResponseWriter, name string, data interface{}) {
	b, e := ioutil.ReadFile(view(layout))
	fmt.Printf("Read file: %v\n", view(layout))
	if e != nil {
		panic("Could not read layout file!")
	}
	t, e := template.New("l").Funcs(LayoutHelpers(name, data)).Parse(string(b))
	if e != nil {
		panic("Could not parse layout file!")
	}
	t.Execute(res, data)
}

func LayoutHelpers(name string, data interface{}) template.FuncMap {
	return template.FuncMap{
		"yield": func() template.HTML {
			var buffer bytes.Buffer

			b, e := ioutil.ReadFile(view(name))
			if e != nil {
				panic(fmt.Sprintf("Could not read: %v", name))
			}
			t := template.New(name).Funcs(newViewHelpers())
			t, e = t.Parse(string(b))
			if e != nil {
				panic(fmt.Sprintf("Could not parse: %v", name))
			}
			t.Execute(&buffer, data)
			return template.HTML(buffer.String())
		},
	}
}

func newViewHelpers() template.FuncMap {
	return template.FuncMap{
		"fmtDate": fmtDate,
		"inc": func(n int) int {
			return n + 1
		},
	}
}

// Returns the view helpers for this application.
// TODO: do not export.
func ViewHelpers() []template.FuncMap {
	return []template.FuncMap{
		{
			"fmtDate": fmtDate,
			"inc": func(n int) int {
				return n + 1
			},
		},
	}
}

// Returns a date formatted in an nicer way.
func fmtDate(t time.Time) string {
	return fmt.Sprintf("%02d/%02d/%04d", t.Day(), t.Month(), t.Year())
}
