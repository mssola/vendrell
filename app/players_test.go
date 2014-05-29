// Copyright (C) 2014 Miquel Sabaté Solà <mikisabate@gmail.com>
// This file is licensed under the MIT license.
// See the LICENSE file.

package app

import (
	"encoding/csv"
	"fmt"
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
	InitTest()
	defer CloseDB()

	m := mux.NewRouter()
	m.HandleFunc("/{id}", PlayersShow)

	// No players.
	req, err := http.NewRequest("GET", "/1", nil)
	assert.Nil(t, err)
	w := httptest.NewRecorder()
	m.ServeHTTP(w, req)
	assert.Equal(t, w.Code, 302)
	errorUrl := []string{"/"}
	assert.Equal(t, w.Header()["Location"], errorUrl)

	// Create some players.
	createPlayer("one", []int{})
	createPlayer("another", []int{1, 2, 3})
	var one, another Player
	err = Db.SelectOne(&one, "select * from players where name=$1", "one")
	assert.Nil(t, err)
	err = Db.SelectOne(&another, "select * from players where name=$1", "another")
	assert.Nil(t, err)
	createUser("user", "1111")

	// Login and perform a couple of requests.
	req, err = http.NewRequest("GET", "/"+one.Id, nil)
	assert.Nil(t, err)
	w = httptest.NewRecorder()
	login(w, req)
	m.ServeHTTP(w, req)
	assert.Equal(t, w.Code, 200)
	assert.Contains(t, w.Body.String(), "<span class=\"empty\">Aquest"+
		" jugador encara no ha valorat cap entrenament.</span>")

	// And now the other player.
	req, err = http.NewRequest("GET", "/"+another.Id, nil)
	assert.Nil(t, err)
	w = httptest.NewRecorder()
	login(w, req)
	m.ServeHTTP(w, req)
	assert.Equal(t, w.Code, 200)
	assert.Contains(t, w.Body.String(),
		`
    <table>
        <tr>
            <th>Mínim</th>
            <th>Màxim</th>
            <th>Mitjana</th>
        </tr>
        <tr>
            <td>1</td>
            <td>3</td>
            <td>2.00</td>
        </tr>
    </table>`)
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
	InitTest()
	defer CloseDB()

	// Invalid Rating (it's not a number).
	m := mux.NewRouter()
	m.HandleFunc("/{id}", PlayersRate)
	r, err := http.NewRequest("POST", "/1", nil)
	assert.Nil(t, err)
	w := httptest.NewRecorder()
	m.ServeHTTP(w, r)
	assert.Equal(t, w.Code, 302)
	errorUrl := []string{"/players/1/rate?error=true"}
	assert.Equal(t, w.Header()["Location"], errorUrl)

	// Invalid Rating (out of range number).
	r, err = http.NewRequest("POST", "/1", nil)
	assert.Nil(t, err)
	param := make(url.Values)
	param["rating"] = []string{"11"}
	r.PostForm = param
	w = httptest.NewRecorder()
	m.ServeHTTP(w, r)
	assert.Equal(t, w.Code, 302)
	assert.Equal(t, w.Header()["Location"], errorUrl)

	// Invalid id parameter.
	r, err = http.NewRequest("POST", "/1", nil)
	assert.Nil(t, err)
	param["rating"] = []string{"5"}
	r.PostForm = param
	w = httptest.NewRecorder()
	m.ServeHTTP(w, r)
	assert.Equal(t, w.Code, 302)
	assert.Equal(t, w.Header()["Location"], errorUrl)

	// Ok.
	createPlayer("user", []int{1, 2, 3})
	var p Player
	err = Db.SelectOne(&p, "select * from players where name=$1", "user")
	assert.Nil(t, err)

	r, err = http.NewRequest("POST", "/"+p.Id, nil)
	assert.Nil(t, err)
	r.PostForm = param
	w = httptest.NewRecorder()
	m.ServeHTTP(w, r)
	assert.Equal(t, w.Code, 302)
	okUrl := []string{fmt.Sprintf("/players/%v/rate", p.Id)}
	assert.Equal(t, w.Header()["Location"], okUrl)
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
