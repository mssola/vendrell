// Copyright (C) 2014 Miquel Sabaté Solà <mikisabate@gmail.com>
// This file is licensed under the MIT license.
// See the LICENSE file.

package main

import (
	"github.com/go-martini/martini"
	_ "github.com/lib/pq"
	"github.com/martini-contrib/render"
	"github.com/martini-contrib/sessions"
	"github.com/mssola/go-utils/security"
	"github.com/mssola/vendrell/app"
)

func main() {
	// The gin in your Martini, the clams on your linguini.
	// We keep the spring in Springfield!
	m := martini.New()

	// Let there be middleware.
	m.Use(martini.Logger())
	m.Use(martini.Recovery())
	m.Use(martini.Static("public"))
	m.Use(render.Renderer(render.Options{
		Layout:     "application/layout",
		Directory:  "views",
		Extensions: []string{".tpl"},
		Funcs:      app.ViewHelpers(),
	}))
	store := sessions.NewCookieStore([]byte(security.NewAuthToken()))
	store.Options(sessions.Options{
		MaxAge: 60 * 60 * 24 * 30 * 12, // A year.
	})
	m.Use(sessions.Sessions("vendrell", store))

	// Database.
	app.InitDB()
	defer app.CloseDB()

	// Routing.
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
