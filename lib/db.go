// Copyright (C) 2014 Miquel Sabaté Solà
// This file is licensed under the MIT license.
// See the LICENSE file.

package lib

import (
	"database/sql"
	"encoding/json"
	"io/ioutil"
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

// Public: initializes the global Db variable.
//
// file - The name of the JSON file describing the database configuration.
//
// Returns a new DB instance.
func InitDB(file string) gorp.DbMap {
	d, err := sql.Open("postgres", DBConfig(file))
	if err != nil {
		panic(err)
	}
	return gorp.DbMap{Db: d, Dialect: gorp.PostgresDialect{}}
}
