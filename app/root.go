// Copyright (C) 2014 Miquel Sabaté Solà <mikisabate@gmail.com>
// This file is licensed under the MIT license.
// See the LICENSE file.

package app

import (
	"fmt"
	"net/http"
	"strings"
)

func parseAgg(agg string) []string {
	clean := strings.TrimPrefix(agg, "{")
	clean = strings.TrimRight(clean, "}")
	return strings.Split(clean, ",")
}

func homePage(res http.ResponseWriter) {
	players := []*ExtendedPlayer{}
	o := &ExtendedHome{LoggedIn: true}

	// TODO: put this in a less scary way.
	q := "select p.id, p.name, min(r.value), max(r.value), avg(r.value),"
	q += " array_agg(r.value) as values, array_agg(r.created_at)"
	q += " from players p, ratings r where r.player_id = p.id"
	q += " group by p.id, p.name"

	rows, _ := Db.Db.Query(q)
	rmax := 0
	for rows.Next() {
		var id, name, values, dates string
		var min, max int
		var avg float64

		if ed := rows.Scan(&id, &name, &min, &max, &avg, &values, &dates); ed == nil {
			vls := parseAgg(values)
			if len(vls) > rmax {
				rmax = len(vls)
			}

			p := &ExtendedPlayer{
				Id:     id,
				Name:   name,
				Min:    min,
				Max:    max,
				Avg:    fmt.Sprintf("%.2f", avg),
				Values: vls,
			}
			players = append(players, p)
		}
	}

	o.Values = make([]int, rmax)
	o.Players = players
	render(res, "root/home", o)
}

func RootIndex(res http.ResponseWriter, req *http.Request) {
	s, _ := store.Get(req, sessionName)
	id := s.Values["userId"]

	if id == nil {
		o := &Options{LoggedIn: false}
		count, err := Db.SelectInt("select count(*) from users")
		if err == nil && count == 0 {
			render(res, "users/new", o)
		} else {
			render(res, "root/index", o)
		}
	} else {
		homePage(res)
	}
}
