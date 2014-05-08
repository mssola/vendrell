// Copyright (C) 2014 Miquel Sabaté Solà <mikisabate@gmail.com>
// This file is licensed under the MIT license.
// See the LICENSE file.

package main

import (
	"github.com/codegangsta/martini"
	"github.com/martini-contrib/render"
	"github.com/mssola/vendrell/app"
)

func main() {
	m := martini.New()

	// Let there be middleware.
	m.Use(martini.Logger())
	m.Use(martini.Recovery())
	m.Use(render.Renderer(render.Options{
		Directory:  "views",
		Extensions: []string{".tpl"},
	}))

	// Routing.
	r := martini.NewRouter()
	r.Get("/", app.RootIndex)
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
