// Copyright (C) 2014 Miquel Sabaté Solà
// This file is licensed under the MIT license.
// See the LICENSE file.

package lib

import (
	"database/sql"
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/coopernurse/gorp"
	_ "github.com/lib/pq"
)

// Public: get the configuration of the database for the current environment.
//
// file - The name of the JSON file describing the database configuration.
//
// Returns a string that it's supposed to be passed as the second argument
// of the sql.Open function.
func DBConfig(file string) string {
	// Read the contents and unmarshal the thing.
	root := FindRoot("vendrell", ".")
	contents, _ := ioutil.ReadFile(root + "/db/" + file)
	m := map[string]map[string]string{}
	json.Unmarshal(contents, &m)

	// Put it in a fancy string.
	current := m[Env()]
	size := len(current)
	cfg := ""
	i := 0
	for key, value := range current {
		if value != "" {
			cfg += key + "=" + value
			if i != size-1 {
				cfg += " "
			}
		}
		i++
	}
	return strings.TrimSpace(cfg)
}

// Returns a config URL for PostgreSQL when using Heroku.
func herokuUrl(env string) string {
	// Black magic to get the PostgreSQL config.
	rg := "(?i)^postgres://(?:([^:@]+):([^@]*)@)?([^@/:]+):(\\d+)/(.*)$"
	regex := regexp.MustCompile(rg)
	matches := regex.FindStringSubmatch(os.Getenv("DATABASE_URL"))
	if matches == nil {
		log.Fatalf("Wrong URL format!")
	}

	// And now we can build a proper url for PostgreSQL.
	fmt := "user=%s password=%s host=%s port=%s dbname=%s sslmode=%s"
	spec := fmt.Sprintf(fmt, matches[1], matches[2], matches[3], matches[4],
		matches[5], "disable")
	return spec
}

// Public: initializes the global Db variable.
//
// file - The name of the JSON file describing the database configuration.
//
// Returns a new DB instance.
func InitDB(file string) gorp.DbMap {
	// First of all get the URL (depends whether it's on Heroku or not).
	url := ""
	if h := os.Getenv("DATABASE_URL"); h != "" {
		url = herokuUrl(h)
	} else {
		url = DBConfig(file)
	}

	// Connect and return the connection.
	d, err := sql.Open("postgres", url)
	if err != nil {
		panic(err)
	}
	return gorp.DbMap{Db: d, Dialect: gorp.PostgresDialect{}}
}
