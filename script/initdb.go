// Copyright (C) 2014 Miquel Sabaté Solà
// This file is licensed under the MIT license.
// See the LICENSE file.

package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/coopernurse/gorp"
	_ "github.com/lib/pq"
	"github.com/mssola/vendrell/app"
	"github.com/mssola/vendrell/lib"
	uuid "github.com/nu7hatch/gouuid"
)

func handleArgs() map[string]string {
	// Getting all the extra arguments.
	args := make(map[string]string)
	for k, v := range os.Args {
		if k > 0 {
			got := strings.Split(v, "=")
			args[got[0]] = got[1]
		}
	}

	// Check that we've defined all the arguments.
	failed := false
	for _, v := range []string{"name", "password"} {
		if _, ok := args[v]; !ok {
			fmt.Printf("Missing argument \"%v\"\n", v)
			failed = true
		}
	}
	if failed {
		os.Exit(1)
	}
	return args
}

func main() {
	args := handleArgs()
	cfg := lib.DBConfig("database.json")
	db, err := sql.Open("postgres", cfg)
	if err != nil {
		log.Fatalln("oops", err)
	}
	dbmap := &gorp.DbMap{Db: db, Dialect: gorp.PostgresDialect{}}
	defer dbmap.Db.Close()

	dbmap.AddTableWithName(app.User{}, "users")
	err = dbmap.Insert(newUser(args))
	if err != nil {
		log.Fatalln("oops", err)
	}
	fmt.Printf("User '%v' created successfully!\n", args["name"])
}

func newUser(args map[string]string) *app.User {
	uuid, err := uuid.NewV4()
	if err != nil {
		log.Fatalln("oops", err)
	}
	return &app.User{
		Id:            uuid.String(),
		Name:          args["name"],
		Password_hash: lib.PasswordSalt(args["password"]),
		Created_at:    time.Now(),
	}
}
