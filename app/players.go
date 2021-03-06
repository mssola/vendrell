// Copyright (C) 2014-2015 Miquel Sabaté Solà <mikisabate@gmail.com>
// This file is licensed under the MIT license.
// See the LICENSE file.

package app

import (
	"encoding/csv"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/nu7hatch/gouuid"
)

func PlayersNew(res http.ResponseWriter, req *http.Request) {
	o := &Options{LoggedIn: true}
	render(res, "players/new", o)
}

func PlayersCreate(res http.ResponseWriter, req *http.Request) {
	// Get a new uuid.
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
	Db.Insert(p)
	http.Redirect(res, req, "/", http.StatusFound)
}

func PlayersShow(res http.ResponseWriter, req *http.Request) {
	// Get the user to be shown.
	params := mux.Vars(req)
	players, _ := getStats(params["id"], true)

	// Let's make sure that the user exists.
	if len(players) == 0 {
		http.Redirect(res, req, "/", http.StatusFound)
		return
	}

	// Prepare parameters and generate the HTML code.
	o := &Options{One: players[0]}
	s, _ := store.Get(req, sessionName)
	id := s.Values["userId"]
	if IsUserLogged(id) {
		o.LoggedIn = true
		o.JS = true
		o.Download = "/players/" + params["id"] + "/csv"
	}
	render(res, "players/show", o)
}

func PlayersUpdate(res http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	query := "update players set name=$1 where id=$2"
	Db.Exec(query, req.FormValue("name"), params["id"])
	http.Redirect(res, req, "/", http.StatusFound)
}

func PlayersDelete(res http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	Db.Exec("delete from players where id=$1 and name=$2",
		params["id"], req.FormValue("name"))
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

func PlayersRate(res http.ResponseWriter, req *http.Request) {
	// Get the rating.
	params := mux.Vars(req)
	rating, err := fetchRating(req.FormValue("rating"))
	if err != nil {
		url := fmt.Sprintf("/players/%v/rate?error=true", params["id"])
		http.Redirect(res, req, url, http.StatusFound)
		return
	}

	// Insert the new rating.
	r := &Rating{
		Value:      rating,
		Player_id:  params["id"],
		Created_at: time.Now(),
	}
	e := Db.Insert(r)

	// Redirect.
	url := fmt.Sprintf("/players/%v/rate", params["id"])
	if e != nil {
		url += "?error=true"
	}
	http.Redirect(res, req, url, http.StatusFound)
}

func PlayersRated(res http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	p := &Options{Id: params["id"]}
	if req.FormValue("error") == "true" {
		p.Error = true
	}
	render(res, "players/rated", p)
}

func PlayersCsv(res http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	players, _ := getStats(params["id"], true)

	// Let's make sure that the user exists.
	if len(players) == 0 {
		http.Redirect(res, req, "/", http.StatusFound)
		return
	}

	// Write the CSV.
	writeCsv(res, players[0].Name, players)
}

func writeCsv(res http.ResponseWriter, name string, players []*ExtPlayer) {
	// Set the headers for CSV.
	res.Header().Set("Content-Type", "text/csv")
	cd := "attachment;filename=" + name + ".csv"
	res.Header().Set("Content-Disposition", cd)

	// Finally write the data.
	w := csv.NewWriter(res)
	for _, v := range players {
		min, max := strconv.Itoa(v.Min), strconv.Itoa(v.Max)
		w.Write([]string{v.Name})
		if v.Avg == "" {
			w.Write([]string{"No ha puntuat cap entrenament"})
		} else {
			w.Write([]string{min, max, v.Avg})

			ratings, dates := []string{}, []string{}
			for _, r := range v.Ratings {
				ratings = append(ratings, strconv.Itoa(r.Value))
				dates = append(dates, fmtDate(r.Created_at))
			}
			w.Write(ratings)
			w.Write(dates)
		}
		w.Write([]string{}) // Extra line.
	}
	w.Flush()
}
