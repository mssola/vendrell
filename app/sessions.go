// Copyright (C) 2014 Miquel Sabaté Solà <mikisabate@gmail.com>
// This file is licensed under the MIT license.
// See the LICENSE file.

package app

import (
	"net/http"

	"github.com/gorilla/sessions"
	"github.com/mssola/go-utils/security"
)

var store *sessions.CookieStore

const sessionName = "vendrell"

func InitSession() {
	store = sessions.NewCookieStore([]byte(security.NewAuthToken()))
	store.Options = &sessions.Options{
		Path:   "/",
		MaxAge: 60 * 60 * 24 * 30 * 12, // A year.
	}
}

func IsUserLogged(id interface{}) bool {
	if id == nil {
		return false
	}

	var u User
	e := Db.SelectOne(&u, "select * from users where id=$1", id.(string))
	return e == nil
}

func UserLogged(res http.ResponseWriter, req *http.Request) {
	s, _ := store.Get(req, sessionName)

	if !IsUserLogged(s.Values["userId"]) {
		http.Redirect(res, req, "/", http.StatusFound)
	}
}

func Login(res http.ResponseWriter, req *http.Request) {
	var u User

	// Check if the user exists and that the password is spot on.
	n, password := req.FormValue("name"), req.FormValue("password")
	e := Db.SelectOne(&u, "select * from users where name=$1", n)
	if e != nil || !security.PasswordMatch(u.Password_hash, password) {
		http.Redirect(res, req, "/", http.StatusFound)
		return
	}

	// It's ok to login this user.
	s, _ := store.Get(req, sessionName)
	s.Values["userId"] = u.Id
	s.Save(req, res)
	http.Redirect(res, req, "/", http.StatusFound)
}

func Logout(res http.ResponseWriter, req *http.Request) {
	s, _ := store.Get(req, sessionName)
	delete(s.Values, "userId")
	s.Save(req, res)

	http.Redirect(res, req, "/", http.StatusFound)
}
