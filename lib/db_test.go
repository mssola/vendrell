// Copyright (C) 2014 Miquel Sabaté Solà
// This file is licensed under the MIT license.
// See the LICENSE file.

package lib

import (
	"os"
	"strings"
	"testing"
)

func mapify(str string) map[string]string {
	res := make(map[string]string)

	s := strings.Split(str, " ")
	for _, v := range s {
		pair := strings.Split(v, "=")
		res[pair[0]] = pair[1]
	}
	return res
}

func TestDBConfig(t *testing.T) {
	// Implicit environment.
	res := mapify(DBConfig("test.json"))
	if res["user"] != "postgres" {
		t.Errorf("Wrong user")
	}
	if res["dbname"] != "vendrell-development" {
		t.Errorf("Wrong dbname")
	}

	// Setting the VENDRELL_ENV environment variable.
	os.Setenv("VENDRELL_ENV", "test")
	res = mapify(DBConfig("test.json"))
	if res["user"] != "postgres" {
		t.Errorf("Wrong user")
	}
	if res["dbname"] != "vendrell-test" {
		t.Errorf("Wrong dbname")
	}

	// Wrong value for LEAKY_ENV.
	os.Setenv("VENDRELL_ENV", "invented")
	str := DBConfig("test.json")
	if str != "" {
		t.Errorf("It should be empty")
	}
}

func TestInitDB(t *testing.T) {
	db := InitDB("test.json")
	if db.Db == nil {
		t.Errorf("Db was not properly initialized")
	}
	db.Db.Close()
}
