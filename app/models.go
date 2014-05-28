// Copyright (C) 2014 Miquel Sabaté Solà <mikisabate@gmail.com>
// This file is licensed under the MIT license.
// See the LICENSE file.

package app

import (
	"time"

	"github.com/coopernurse/gorp"
	"github.com/mssola/go-utils/db"
	"github.com/mssola/go-utils/misc"
	"github.com/mssola/go-utils/path"
)

// These are the tables stored in the DB.

type User struct {
	Id            string
	Name          string
	Password_hash string
	Created_at    time.Time
}

type Player struct {
	Id         string
	Name       string
	Created_at time.Time
}

type Rating struct {
	Id         int
	Value      int
	Player_id  string
	Created_at time.Time
}

// Global instance that holds a connection to the DB. It gets initialized after
// calling the InitDB function. You have to call CloseDB in order to close the
// connection.
var Db gorp.DbMap

// Initialize the global DB connection.
func InitDB() {
	c := db.Open(db.Options{
		Base:        path.FindRoot("vendrell", "."),
		Relative:    "/db/database.json",
		Environment: misc.EnvOrElse("VENDRELL_ENV", "development"),
		DBMS:        "postgres",
		Heroku:      true,
	})
	Db = gorp.DbMap{Db: c, Dialect: gorp.PostgresDialect{}}
	Db.AddTableWithName(User{}, "users")
	Db.AddTableWithName(Player{}, "players")
	Db.AddTableWithName(Rating{}, "ratings").SetKeys(true, "Id")
}

// Close the global DB connection.
func CloseDB() {
	Db.Db.Close()
}
