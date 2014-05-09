// Copyright (C) 2014 Miquel Sabaté Solà <mikisabate@gmail.com>
// This file is licensed under the MIT license.
// See the LICENSE file.

package app

import (
	"net/http"

	"github.com/coopernurse/gorp"
	"github.com/martini-contrib/sessions"
	"github.com/mssola/leaky/lib"
)

func Login(res http.ResponseWriter, req *http.Request, db gorp.DbMap, s sessions.Session) {
	var u User

	// Check if the user exists and that the password is spot on.
	n, password := req.FormValue("name"), req.FormValue("password")
	e := db.SelectOne(&u, "select * from users where name=$1", n)
	if e != nil || !lib.PasswordMatch(u.Password_hash, password) {
		http.Redirect(res, req, "/", http.StatusNotFound)
		return
	}

	// It's ok to login this user.
	s.Set("userId", u.Id)
	http.Redirect(res, req, "/", http.StatusFound)
}

func Logout(res http.ResponseWriter, req *http.Request, s sessions.Session) {
	s.Delete("userId")
	http.Redirect(res, req, "/", http.StatusFound)
}
