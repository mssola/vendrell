// Copyright (C) 2014 Miquel Sabaté Solà <mikisabate@gmail.com>
// This file is licensed under the MIT license.
// See the LICENSE file.

package app

import (
	"github.com/coopernurse/gorp"
	"github.com/martini-contrib/render"
	"github.com/martini-contrib/sessions"
)

func RootIndex(db gorp.DbMap, r render.Render, s sessions.Session) {
	id := s.Get("userId")
	if id == nil {
		count, err := db.SelectInt("select count(*) from users")
		if err == nil && count == 0 {
			r.HTML(200, "users/new", "mssola")
		} else {
			r.HTML(200, "root/index", "mssola")
		}
	} else {
		r.HTML(200, "root/home", "mssola")
	}
}
