// Copyright (C) 2014 Miquel Sabaté Solà <mikisabate@gmail.com>
// This file is licensed under the MIT license.
// See the LICENSE file.

package main

import (
	"fmt"

	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"github.com/mssola/go-utils/misc"
	"github.com/mssola/vendrell/app"
)

func main() {
	// Because Martini was too mainstream :P
	n := negroni.Classic()

	// Sessions.
	app.InitSession()

	// Database.
	app.InitDB()
	defer app.CloseDB()

	// Routing.
	r := mux.NewRouter()
	r.HandleFunc("/", app.RootIndex).Methods("GET")
	r.HandleFunc("/login", app.Login).Methods("POST")
	r.HandleFunc("/logout", app.Logout).Methods("POST").
		MatcherFunc(app.UserLogged)
	r.HandleFunc("/users", app.UsersCreate).Methods("POST")
	r.HandleFunc("/players", app.PlayersCreate).Methods("POST").
		MatcherFunc(app.UserLogged)
	r.HandleFunc("/players", app.PlayersIndex).Methods("GET")
	s := r.PathPrefix("/players").Subrouter()
	s.HandleFunc("/new", app.PlayersNew).Methods("GET").
		MatcherFunc(app.UserLogged)
	s.HandleFunc("/{id}", app.PlayersShow).Methods("GET")
	s.HandleFunc("/{id}", app.PlayersUpdate).Methods("POST").
		MatcherFunc(app.UserLogged)
	s.HandleFunc("/{id}/delete", app.PlayersDelete).Methods("POST").
		MatcherFunc(app.UserLogged)
	s.HandleFunc("/{id}/rate", app.PlayersRate).Methods("POST")
	s.HandleFunc("/{id}/rate", app.PlayersRated).Methods("GET")
	s.HandleFunc("/{id}/csv", app.PlayersCsv).Methods("GET").
		MatcherFunc(app.UserLogged)
	r.HandleFunc("/csv", app.RootCsv).Methods("GET").
		MatcherFunc(app.UserLogged)
	n.UseHandler(r)

	// Run, Forrest, run!
	port := fmt.Sprintf(":%v", misc.EnvOrElse("PORT", "3000"))
	n.Run(port)
}
