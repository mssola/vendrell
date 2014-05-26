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

const (
	layout   = "application/layout"
	viewsDir = "views"
	viewsExt = "tpl"
)

type Options struct {
	Id       string
	One      *ExtPlayer
	Players  []*ExtPlayer
	Values   []int
	LoggedIn bool
	JS       bool
	Error    bool
}

func view(name string) string {
	return path.Join(viewsDir, name+"."+viewsExt)
}

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

func viewHelpers() template.FuncMap {
	return template.FuncMap{
		"fmtDate": func(t time.Time) string {
			return fmt.Sprintf("%02d/%02d/%04d", t.Day(), t.Month(), t.Year())
		},
		"inc": func(n int) int {
			return n + 1
		},
	}
}
