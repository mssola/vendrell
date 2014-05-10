// Copyright (C) 2014 Miquel Sabaté Solà <mikisabate@gmail.com>
// This file is licensed under the MIT license.
// See the LICENSE file.

package app

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/coopernurse/gorp"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
	"github.com/martini-contrib/sessions"
	"github.com/nu7hatch/gouuid"
)

func PlayersCreate(res http.ResponseWriter, req *http.Request, db gorp.DbMap) {
	// Get a ne uuid.
	id, err := uuid.NewV4()
	if err != nil {
		http.Redirect(res, req, "/", http.StatusNotFound)
		return
	}

	// Try to create a new user and redirect properly.
	p := &Player{
		Id:         id.String(),
		Name:       req.FormValue("name"),
		Created_at: time.Now(),
	}
	db.Insert(p)
	http.Redirect(res, req, "/", http.StatusFound)
}

func PlayersShow(res http.ResponseWriter, req *http.Request, r render.Render,
	params martini.Params, db gorp.DbMap, s sessions.Session) {

	var p Player

	// Get the user to be shown.
	e := db.SelectOne(&p, "select * from players where id=$1", params["id"])
	if e != nil {
		http.Redirect(res, req, "/", http.StatusFound)
		return
	}

	// Prepare parameters and generate the HTML code.
	tpl := map[string]string{"id": p.Id, "name": p.Name}
	id := s.Get("userId")
	if IsUserLogged(id, db) {
		tpl["logged"] = "true"
	}
	r.HTML(200, "players/show", tpl)
}

func PlayersUpdate(res http.ResponseWriter, req *http.Request,
	params martini.Params, db gorp.DbMap) {

	query := "update players set name=$1 where id=$2"
	db.Exec(query, req.FormValue("name"), params["id"])
	http.Redirect(res, req, "/", http.StatusFound)
}

func PlayersDelete(res http.ResponseWriter, req *http.Request,
	params martini.Params, db gorp.DbMap) {

	db.Exec("delete from players where id=$1", params["id"])
	http.Redirect(res, req, "/", http.StatusFound)
}

func fetchRating(rating string) (int, error) {
	r, err := strconv.Atoi(rating)
	if err != nil {
		return 0, err
	}

	if r >= 0 && r <= 10 {
		return r, nil
	}
	return 0, errors.New("Invalid rating!")
}

func PlayersRate(res http.ResponseWriter, req *http.Request,
	params martini.Params, db gorp.DbMap, s sessions.Session) {

	// Get the rating.
	rating, err := fetchRating(req.FormValue("rating"))
	if err != nil {
		http.Redirect(res, req, "/", http.StatusFound)
		return
	}

	// Insert the new rating and redirect the user.
	r := &Rating{
		Value:      rating,
		Player_id:  params["id"],
		Created_at: time.Now(),
	}
	db.Insert(r)
	// TODO: redirect to a proper page.
	http.Redirect(res, req, "/", http.StatusFound)
}
