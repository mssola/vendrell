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

// A Player with some useful data.
type ExtPlayer struct {
	// The Id of the player (uuid).
	Id string

	// The name of the player.
	Name string

	// The minimum rating from the player.
	Min int

	// The maximum rating from the player.
	Max int

	// The average rating from the player.
	Avg string

	// All the ratings from this player.
	Ratings []Rating

	// The date of creation of this player.
	Created_at time.Time
}

// Returns a string containing the creepy SQL statement to be executed in
// order to obtain the stats of a player.
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

// Returns the inner values of a PostgreSQL's array_agg.
func parseAgg(agg string) []string {
	clean := strings.TrimPrefix(agg, "{")
	clean = strings.TrimRight(clean, "}")
	return strings.Split(clean, ",")
}

// Returns three integers that are inside the given string. This given string
// is basically a set of integers separated by the given sep value.
func mustAtoi(str, sep string) (int, int, int) {
	vals := strings.SplitN(str, sep, 3)
	first, _ := strconv.Atoi(vals[0])
	second, _ := strconv.Atoi(vals[1])
	third, _ := strconv.Atoi(vals[2])
	return first, second, third
}

// Parse a date from a PostgreSQL's timestamp. PostgreSQL's timestamps have the
// following format: 2014-05-15 21:41:21.1234
func parseDate(complete string) time.Time {
	complete = strings.TrimPrefix(complete, "\"")
	complete = strings.TrimSuffix(complete, "\"")

	spd := strings.SplitN(complete, " ", 2)
	y, mo, d := mustAtoi(spd[0], "-")
	hour := strings.SplitN(spd[1], ".", 2)
	h, m, s := mustAtoi(hour[0], ":")
	n, _ := strconv.Atoi(hour[1])
	return time.Date(y, time.Month(mo), d, h, m, s, n, time.UTC)
}

// Returns a list of ratings that can be extracted from the given values and
// dates. These values and dates are in PostgreSQL's array_agg format.
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

// This function fetches the statistics of the required players. If the "one"
// parameter is set to true, then it's interpreted that the player to be
// fetched has the ID given by the playerId parameter. Otherwise, if the "one"
// parameter is set to false, then the playerId parameter will be ignored and
// all the users that have rated a practice at least once will be fetched.
//
// Returns the list of the fetched players plus an integer value. This integer
// value represents the maximum number of practices that the selected players
// have rated.
func getStats(playerId string, one bool) ([]*ExtPlayer, int) {
	var rows *sql.Rows
	var err error

	// Prepare the query.
	q := statsQuery(one)
	if one {
		rows, err = Db.Db.Query(q, playerId)
	} else {
		rows, err = Db.Db.Query(q)
	}
	if err != nil {
		return []*ExtPlayer{}, 0
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
