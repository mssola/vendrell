// Copyright (C) 2014 Miquel Sabaté Solà <mikisabate@gmail.com>
// This file is licensed under the MIT license.
// See the LICENSE file.

package main

import (
	"github.com/go-martini/martini"
	_ "github.com/lib/pq"
	"github.com/martini-contrib/render"
	"github.com/mssola/vendrell/app"
)

func main() {
	// The gin in your Martini, the clams on your linguini.
	// We keep the spring in Springfield!
	m := martini.New()

	// Let there be middleware.
	// TODO: can be replaced with negroni.
	m.Use(martini.Logger())
	m.Use(martini.Recovery())
	m.Use(martini.Static("public"))

	// TODO: what can I do here ? :/
	m.Use(render.Renderer(render.Options{
		Layout:     "application/layout",
		Directory:  "views",
		Extensions: []string{".tpl"},
		Funcs:      app.ViewHelpers(),
	}))

	// Sessions.
	app.InitSession()

	// Database.
	app.InitDB()
	defer app.CloseDB()

	// Routing.
	// TODO: replace with Gorilla's mux package.
	r := martini.NewRouter()
	r.Get("/", app.RootIndex)
	r.Post("/login", app.Login)
	r.Post("/logout", app.UserLogged, app.Logout)
	r.Post("/users", app.UsersCreate)
	r.Group("/players", func(r martini.Router) {
		r.Get("/new", app.PlayersNew)
		r.Post("", app.UserLogged, app.PlayersCreate)
		r.Get("/:id", app.PlayersShow)
		r.Post("/:id", app.UserLogged, app.PlayersUpdate)
		r.Post("/:id/delete", app.UserLogged, app.PlayersDelete)
		r.Post("/:id/rate", app.PlayersRate)
		r.Get("/:id/rate", app.PlayersRated)
	})
	m.Action(r.Handle)

	// Run, Forrest, run!
	m.Run()
}
