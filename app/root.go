// Copyright (C) 2014 Miquel Sabaté Solà <mikisabate@gmail.com>
// This file is licensed under the MIT license.
// See the LICENSE file.

package app

import (
	"github.com/martini-contrib/render"
	"github.com/martini-contrib/sessions"
)

func RootIndex(r render.Render, s sessions.Session) {
	id := s.Get("userId")
	if id == nil {
		r.HTML(200, "root/index", "mssola")
	} else {
		r.HTML(200, "root/home", "mssola")
	}
}
