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
	store.Options(sessions.Options{
		MaxAge: 60 * 60 * 24 * 30 * 12, // A year.
	})
	m.Use(sessions.Sessions("vendrell", store))

	// Database.
	db := lib.InitDB("database.json")
	db.AddTableWithName(app.User{}, "users")
	db.AddTableWithName(app.Player{}, "players")
	db.AddTableWithName(app.Rating{}, "ratings").SetKeys(true, "Id")
	m.Map(db)
	defer db.Db.Close()

	// Routing.
	r := martini.NewRouter()
	r.Get("/", app.RootIndex)
	r.Post("/login", app.Login)
	r.Post("/logout", app.UserLogged, app.Logout)
	r.Post("/users", app.UsersCreate)
	r.Group("/players", func(r martini.Router) {
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
