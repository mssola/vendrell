// Copyright (C) 2014 Miquel Sabaté Solà <mikisabate@gmail.com>
// This file is licensed under the MIT license.
// See the LICENSE file.

package main

import (
	"github.com/coopernurse/gorp"
	"github.com/go-martini/martini"
	_ "github.com/lib/pq"
	"github.com/martini-contrib/render"
	"github.com/martini-contrib/sessions"
	"github.com/mssola/go-utils/db"
	"github.com/mssola/go-utils/misc"
	"github.com/mssola/go-utils/path"
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
	}))
	store := sessions.NewCookieStore([]byte(security.NewAuthToken()))
	store.Options(sessions.Options{
		MaxAge: 60 * 60 * 24 * 30 * 12, // A year.
	})
	m.Use(sessions.Sessions("vendrell", store))

	// Database.
	c := db.Open(db.Options{
		Base:        path.FindRoot("vendrell", "."),
		Relative:    "/db/database.json",
		Environment: misc.EnvOrElse("VENDRELL_ENV", "development"),
		DBMS:        "postgres",
		Heroku:      true,
	})
	d := gorp.DbMap{Db: c, Dialect: gorp.PostgresDialect{}}
	d.AddTableWithName(app.User{}, "users")
	d.AddTableWithName(app.Player{}, "players")
	d.AddTableWithName(app.Rating{}, "ratings").SetKeys(true, "Id")
	m.Map(d)
	defer d.Db.Close()

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
