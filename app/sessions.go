// Copyright (C) 2014 Miquel Sabaté Solà <mikisabate@gmail.com>
// This file is licensed under the MIT license.
// See the LICENSE file.

package app

import (
	"net/http"

	"github.com/martini-contrib/sessions"
	"github.com/mssola/go-utils/security"
)

func IsUserLogged(id interface{}) bool {
	if id == nil {
		return false
	}

	var u User
	e := Db.SelectOne(&u, "select * from users where id=$1", id.(string))
	return e == nil
}

func UserLogged(
	res http.ResponseWriter,
	req *http.Request,
	s sessions.Session,
) {
	id := s.Get("userId")
	if !IsUserLogged(id) {
		http.Redirect(res, req, "/", http.StatusFound)
	}
}

func Login(
	res http.ResponseWriter,
	req *http.Request,
	s sessions.Session,
) {
	var u User

	// Check if the user exists and that the password is spot on.
	n, password := req.FormValue("name"), req.FormValue("password")
	e := Db.SelectOne(&u, "select * from users where name=$1", n)
	if e != nil || !security.PasswordMatch(u.Password_hash, password) {
		http.Redirect(res, req, "/", http.StatusFound)
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
