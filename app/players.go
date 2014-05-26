// Copyright (C) 2014 Miquel Sabaté Solà <mikisabate@gmail.com>
// This file is licensed under the MIT license.
// See the LICENSE file.

package app

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/nu7hatch/gouuid"
)

// TODO: move this into app/stats.go ?

// TODO: remove ?
type Statistics struct {
	Ratings []Rating
	Min     int
	Max     int
	Avg     string
}

func statsQuery(one bool) string {
	q := "select p.id, p.name, min(r.value), max(r.value), avg(r.value),"
	q += " array_agg(r.value) as values, array_agg(r.created_at)"
	q += " from players p, ratings r where r.player_id = p.id"
	if one {
		q += " and p.id = $1"
	}
	q += " group by p.id, p.name"
	return q
}

// TODO: remove
const psqlFmt = "2014-05-15 21:41:21"

func mustAtoi(str, sep string) (int, int, int) {
	// TODO: handle error.
	// TODO: SplitN
	vals := strings.Split(str, sep)
	first, _ := strconv.Atoi(vals[0])
	second, _ := strconv.Atoi(vals[1])
	third, _ := strconv.Atoi(vals[2])
	return first, second, third
}

func parseDate(complete string) time.Time {
	// TODO: handle error
	// TODO SPlitN

	// TODO: I'm sure there's a better way...
	complete = strings.TrimPrefix(complete, "\"")
	complete = strings.TrimSuffix(complete, "\"")

	spd := strings.Split(complete, " ")
	y, mo, d := mustAtoi(spd[0], "-")
	hour := strings.Split(spd[1], ".")
	h, m, s := mustAtoi(hour[0], ":")
	n, _ := strconv.Atoi(hour[1])
	return time.Date(y, time.Month(mo), d, h, m, s, n, time.UTC)
}

func parseRatings(values, dates string) []Rating {
	ratings := []Rating{}
	vls := parseAgg(values)
	d := parseAgg(dates)

	for i := 0; i < len(vls); i += 1 {
		v, _ := strconv.Atoi(vls[i])
		t := parseDate(d[i])

		nr := Rating{
			Id:         0, // We don't care.
			Value:      v,
			Player_id:  "", // We don't care.
			Created_at: t,
		}
		ratings = append(ratings, nr)
	}
	return ratings
}

func newGetStats(playerId string, one bool) ([]*NewPlayer, int) {
	var rows *sql.Rows

	// Prepare the query.
	q := statsQuery(one)
	if one {
		rows, _ = Db.Db.Query(q, playerId)
	} else {
		rows, _ = Db.Db.Query(q)
	}

	// And fetch the players and their ratings.
	rmax := 0
	players := []*NewPlayer{}
	for rows.Next() {
		var id, name, values, dates string
		var min, max int
		var avg float64

		if ed := rows.Scan(&id, &name, &min, &max, &avg, &values, &dates); ed == nil {
			ratings := parseRatings(values, dates)
			if len(ratings) > rmax {
				rmax = len(ratings)
			}

			p := &NewPlayer{
				Id:      id,
				Name:    name,
				Min:     min,
				Max:     max,
				Avg:     fmt.Sprintf("%.2f", avg),
				Ratings: ratings,
			}
			players = append(players, p)
		}
	}
	return players, rmax
}

// TODO
func getStats(id string) (*Statistics, error) {
	s := &Statistics{}

	query := "select * from ratings where player_id=$1"
	if _, e := Db.Select(&s.Ratings, query, id); e != nil {
		return nil, e
	}

	count := 0.0
	for _, v := range s.Ratings {
		if v.Value < s.Min {
			s.Min = v.Value
		} else if v.Value > s.Max {
			s.Max = v.Value
		}
		count += float64(v.Value)
	}
	avg := count / float64(len(s.Ratings))
	s.Avg = fmt.Sprintf("%0.2f", avg)
	return s, nil
}

func PlayersNew(res http.ResponseWriter, req *http.Request) {
	o := &Options{LoggedIn: true}
	render(res, "players/new", o)
}

func PlayersCreate(res http.ResponseWriter, req *http.Request) {
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
	Db.Insert(p)
	http.Redirect(res, req, "/", http.StatusFound)
}

func PlayersShow(res http.ResponseWriter, req *http.Request) {
	// Get the user to be shown.
	params := mux.Vars(req)
	players, _ := newGetStats(params["id"], true)

	// Let's make sure that the user exists.
	if len(players) == 0 {
		http.Redirect(res, req, "/", http.StatusFound)
		return
	}

	// Prepare parameters and generate the HTML code.
	o := &NewOptions{One: players[0]}
	s, _ := store.Get(req, sessionName)
	id := s.Values["userId"]
	if IsUserLogged(id) {
		o.LoggedIn = true
		o.JS = true
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
