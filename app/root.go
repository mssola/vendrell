// Copyright (C) 2014 Miquel Sabaté Solà <mikisabate@gmail.com>
// This file is licensed under the MIT license.
// See the LICENSE file.

package app

import (
	"github.com/coopernurse/gorp"
	"github.com/martini-contrib/render"
	"github.com/martini-contrib/sessions"
)

type Options struct {
	LoggedIn bool
}

func RootIndex(db gorp.DbMap, r render.Render, s sessions.Session) {
	id := s.Get("userId")
	if id == nil {
		o := &Options{LoggedIn: false}
		count, err := db.SelectInt("select count(*) from users")
		if err == nil && count == 0 {
			r.HTML(200, "users/new", o)
		} else {
			r.HTML(200, "root/index", o)
		}
	} else {
		o := &Options{LoggedIn: true}
		r.HTML(200, "root/home", o)
	}
}
