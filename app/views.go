// Copyright (C) 2014 Miquel Sabaté Solà <mikisabate@gmail.com>
// This file is licensed under the MIT license.
// See the LICENSE file.

package app

import (
	"fmt"
	"html/template"
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

type Home struct {
	Players  []Player
	LoggedIn bool
	JS       bool
}

// Returns the view helpers for this application.
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
