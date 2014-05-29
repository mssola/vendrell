// Copyright (C) 2014 Miquel Sabaté Solà <mikisabate@gmail.com>
// This file is licensed under the MIT license.
// See the LICENSE file.

package app

import (
	"encoding/csv"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func TestPlayersNew(t *testing.T) {
	InitTest()
	defer CloseDB()

	req, err := http.NewRequest("GET", "/", nil)
	assert.Nil(t, err)
	w := httptest.NewRecorder()
	PlayersNew(w, req)

	assert.Equal(t, w.Code, 200)
	assert.Contains(t, w.Body.String(),
		"<a href=\"/players/new\">Crear jugador</a>")
}

func TestPlayersCreate(t *testing.T) {
	InitTest()
	defer CloseDB()

	param := make(url.Values)

	// "name" is not set.
	req, err := http.NewRequest("POST", "/", nil)
	assert.Nil(t, err)
	req.PostForm = param
	w := httptest.NewRecorder()
	PlayersCreate(w, req)

	assert.Equal(t, w.Code, 302)
	count, err := Db.SelectInt("select count(*) from players")
	assert.Equal(t, count, 0)

	// Creating a player.
	param["name"] = []string{"player"}
	req, err = http.NewRequest("POST", "/", nil)
	assert.Nil(t, err)
	req.PostForm = param
	w = httptest.NewRecorder()
	PlayersCreate(w, req)

	assert.Equal(t, w.Code, 302)
	count, err = Db.SelectInt("select count(*) from players")
	assert.Equal(t, count, 1)

	// You can't create the same player.
	req, err = http.NewRequest("POST", "/", nil)
	assert.Nil(t, err)
	req.PostForm = param
	w = httptest.NewRecorder()
	PlayersCreate(w, req)

	assert.Equal(t, w.Code, 302)
	count, err = Db.SelectInt("select count(*) from players")
	assert.Equal(t, count, 1)
}

func TestPlayersShow(t *testing.T) {
	// TODO
}

func TestPlayersUpdate(t *testing.T) {
	InitTest()
	defer CloseDB()

	createPlayer("one", []int{1, 2, 3})
	var p Player
	err := Db.SelectOne(&p, "select * from players where name=$1", "one")
	assert.Nil(t, err)

	param := make(url.Values)
	param["name"] = []string{"another"}

	req, err := http.NewRequest("POST", "/"+p.Id, nil)
	assert.Nil(t, err)
	req.PostForm = param
	w := httptest.NewRecorder()

	m := mux.NewRouter()
	m.HandleFunc("/{id}", PlayersUpdate)
	m.ServeHTTP(w, req)

	assert.Equal(t, w.Code, 302)
	err = Db.SelectOne(&p, "select * from players")
	assert.Equal(t, p.Name, "another")
}

func TestPlayersDelete(t *testing.T) {
	InitTest()
	defer CloseDB()

	createPlayer("one", []int{1, 2, 3})
	var p Player
	err := Db.SelectOne(&p, "select * from players where name=$1", "one")
	assert.Nil(t, err)

	// Nothing has to happen since we picked the wrong guy.
	param := make(url.Values)
	param["name"] = []string{"another"}

	req, err := http.NewRequest("POST", "/"+p.Id, nil)
	assert.Nil(t, err)
	req.PostForm = param
	w := httptest.NewRecorder()

	m := mux.NewRouter()
	m.HandleFunc("/{id}", PlayersDelete)
	m.ServeHTTP(w, req)

	assert.Equal(t, w.Code, 302)
	err = Db.SelectOne(&p, "select * from players")
	assert.Equal(t, p.Name, "one")

	// Now we pick the right one.
	param["name"] = []string{"one"}
	req, err = http.NewRequest("POST", "/"+p.Id, nil)
	assert.Nil(t, err)
	req.PostForm = param
	w = httptest.NewRecorder()
	m.ServeHTTP(w, req)

	assert.Equal(t, w.Code, 302)
	count, err := Db.SelectInt("select count(*) from players")
	assert.Equal(t, count, 0)
}

func TestPlayersRate(t *testing.T) {
	// TODO
}

func TestPlayersRated(t *testing.T) {
	InitTest()
	defer CloseDB()

	m := mux.NewRouter()
	m.HandleFunc("/{id}", PlayersRated)

	// Ok.
	r, err := http.NewRequest("GET", "/1", nil)
	assert.Nil(t, err)
	w := httptest.NewRecorder()
	m.ServeHTTP(w, r)
	assert.Equal(t, w.Code, 200)
	assert.Contains(t, w.Body.String(), "<h1>Ho tenim !</h1>")

	// Error.
	param := make(url.Values)
	param["error"] = []string{"true"}
	r, err = http.NewRequest("GET", "/1", nil)
	assert.Nil(t, err)
	r.PostForm = param
	w = httptest.NewRecorder()
	m.ServeHTTP(w, r)
	assert.Equal(t, w.Code, 200)
	assert.Contains(t, w.Body.String(), "<h1>Error !</h1>")
}

func TestPlayersCsv(t *testing.T) {
	InitTest()
	defer CloseDB()

	// Someone that doesn't exist.
	m := mux.NewRouter()
	m.HandleFunc("/{id}", PlayersCsv)
	r, err := http.NewRequest("GET", "/1", nil)
	assert.Nil(t, err)
	w := httptest.NewRecorder()
	m.ServeHTTP(w, r)

	assert.Equal(t, w.Code, 302)
	header := w.Header()
	assert.Equal(t, header["Location"], []string{"/"})

	// Let's create a couple of players and some ratings.
	createPlayer("one", []int{1, 2, 3})
	createPlayer("another", []int{0, 3, 8})

	var p Player
	err = Db.SelectOne(&p, "select * from players where name=$1", "one")
	assert.Nil(t, err)

	// Perform the request.
	r, err = http.NewRequest("GET", "/"+p.Id, nil)
	assert.Nil(t, err)
	w = httptest.NewRecorder()
	m.ServeHTTP(w, r)

	// HTTP
	assert.Equal(t, w.Code, 200)
	header = w.Header()
	assert.Equal(t, header["Content-Type"][0], "text/csv")
	assert.Equal(t, header["Content-Disposition"][0],
		"attachment;filename=one.csv")

	dt := fmtDate(time.Now())

	// CSV
	re := csv.NewReader(w.Body)
	testCSV(t, re, 4, "one", "1", "3", "2.00")
	testCSV(t, re, 7, "one", "1", dt, "2", dt, "3", dt)
}
