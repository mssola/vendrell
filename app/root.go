// Copyright (C) 2014 Miquel Sabaté Solà <mikisabate@gmail.com>
// This file is licensed under the MIT license.
// See the LICENSE file.

package app

import (
	"github.com/coopernurse/gorp"
	"github.com/martini-contrib/render"
	"github.com/martini-contrib/sessions"
)

// TODO: move this to another place.
type Options struct {
	Id       string
	Name     string
	LoggedIn bool
	Error    bool
}

type Home struct {
	Players  []Player
	LoggedIn bool
}

func homePage(db gorp.DbMap, r render.Render, s sessions.Session) {
	var players []Player
	o := &Home{LoggedIn: true}

	_, err := db.Select(&players, "select * from players order by name")
	if err != nil {
		// TODO
		return
	}
	o.Players = players
	r.HTML(200, "root/home", o)
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
		homePage(db, r, s)
	}
}
