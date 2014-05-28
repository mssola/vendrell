// Copyright (C) 2014 Miquel Sabaté Solà <mikisabate@gmail.com>
// This file is licensed under the MIT license.
// See the LICENSE file.

package app

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"
	"time"
)

type ExtPlayer struct {
	Id         string
	Name       string
	Min        int
	Max        int
	Avg        string
	Ratings    []Rating
	Created_at time.Time
}

func statsQuery(one bool) string {
	q := "select p.id, p.name, min(r.value), max(r.value), avg(r.value),"
	q += " array_agg(r.value) as values, array_agg(r.created_at)"
	q += " from players p, ratings r where r.player_id = p.id"
	if one {
		q += " and p.id = $1"
	}
	q += " group by p.id, p.name order by p.name"
	return q
}

func parseAgg(agg string) []string {
	clean := strings.TrimPrefix(agg, "{")
	clean = strings.TrimRight(clean, "}")
	return strings.Split(clean, ",")
}

func mustAtoi(str, sep string) (int, int, int) {
	// TODO: handle error.
	// TODO: SplitN
	vals := strings.Split(str, sep)
	first, _ := strconv.Atoi(vals[0])
	second, _ := strconv.Atoi(vals[1])
	third, _ := strconv.Atoi(vals[2])
	return first, second, third
}

// PostgreSQL's timestamps have the following format: 2014-05-15 21:41:21.1234
// TODO
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

func getStats(playerId string, one bool) ([]*ExtPlayer, int) {
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
	players := []*ExtPlayer{}
	for rows.Next() {
		var id, name, values, dates string
		var min, max int
		var avg float64

		if ed := rows.Scan(&id, &name, &min, &max, &avg, &values, &dates); ed == nil {
			ratings := parseRatings(values, dates)
			if len(ratings) > rmax {
				rmax = len(ratings)
			}

			p := &ExtPlayer{
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
