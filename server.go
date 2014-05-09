// Copyright (C) 2014 Miquel Sabaté Solà <mikisabate@gmail.com>
// This file is licensed under the MIT license.
// See the LICENSE file.

package main

import (
	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
	"github.com/martini-contrib/sessions"
	"github.com/mssola/vendrell/app"
	"github.com/mssola/vendrell/lib"
)

func main() {
	// The gin in your Martini, the clams on your linguini.
	// We keep the spring in Springfield!
	m := martini.New()

	// Let there be middleware.
	m.Use(martini.Logger())
	m.Use(martini.Recovery())
	m.Use(render.Renderer(render.Options{
		Layout:     "application/layout",
		Directory:  "views",
		Extensions: []string{".tpl"},
	}))
	store := sessions.NewCookieStore([]byte(lib.NewAuthToken()))
	m.Use(sessions.Sessions("vendrell", store))

	// Database.
	db := lib.InitDB("database.json")
	m.Map(db)
	defer db.Db.Close()

	// Routing.
	r := martini.NewRouter()
	r.Get("/", app.RootIndex)
	r.Post("/login", app.Login)
	r.Post("/logout", app.Logout)
	r.Group("/players", func(r martini.Router) {
		r.Get("/new", app.PlayersNew)
		r.Post("/", app.PlayersCreate)
		r.Get("/:id", app.PlayersShow)
		r.Put("/:id", app.PlayersUpdate)
		r.Delete("/:id", app.PlayersDelete)
	})
	m.Action(r.Handle)

	// Run, Forrest, run!
	m.Run()
}
