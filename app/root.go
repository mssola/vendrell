// Copyright (C) 2014-2015 Miquel Sabaté Solà <mikisabate@gmail.com>
// This file is licensed under the MIT license.
// See the LICENSE file.

package app

import (
	"net/http"
)

// Renders the home page for a logged in user.
func homePage(res http.ResponseWriter) {
	players, rmax := getStats("", false)

	o := &Options{
		LoggedIn: true,
		Values:   make([]int, rmax),
		Players:  players,
		Download: "/csv",
	}
	render(res, "root/home", o)
}

// Renders the root page. It has three different options:
//
//  1. If there's no user, it renders the "Create user" page.
//  2. If the current user is not logged in, it render the "Login" page.
//  3. If the current user is logged in, then it call homePage.
func RootIndex(res http.ResponseWriter, req *http.Request) {
	s, _ := store.Get(req, sessionName)
	id := s.Values["userId"]

	if id == nil {
		o := &Options{LoggedIn: false}
		count, err := Db.SelectInt("select count(*) from users")
		if err == nil && count == 0 {
			render(res, "users/new", o)
		} else {
			render(res, "root/index", o)
		}
	} else {
		homePage(res)
	}
}

func RootCsv(res http.ResponseWriter, req *http.Request) {
	players, _ := getStats("", false)
	writeCsv(res, "data", players)
}
