// Copyright (C) 2014 Miquel Sabaté Solà <mikisabate@gmail.com>
// This file is licensed under the MIT license.
// See the LICENSE file.

package app

import (
	"encoding/csv"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"

	"github.com/nu7hatch/gouuid"
	"github.com/stretchr/testify/assert"
)

func TestCreateUserPage(t *testing.T) {
	InitTest()
	defer CloseDB()

	req, err := http.NewRequest("GET", "/", nil)
	assert.Nil(t, err)
	w := httptest.NewRecorder()
	RootIndex(w, req)

	assert.Equal(t, w.Code, 200)
	assert.Contains(t, w.Body.String(), "<h1>Crear usuari</h1>")
}

func TestLoginPage(t *testing.T) {
	InitTest()
	defer CloseDB()

	createUser("user", "1111")
	req, err := http.NewRequest("GET", "/", nil)
	assert.Nil(t, err)
	w := httptest.NewRecorder()
	RootIndex(w, req)

	assert.Equal(t, w.Code, 200)
	assert.Contains(t, w.Body.String(), "<h1>Login</h1>")
}

func TestRootPage(t *testing.T) {
	InitTest()
	defer CloseDB()

	createUser("user", "1111")
	req, err := http.NewRequest("GET", "/", nil)
	assert.Nil(t, err)
	w := httptest.NewRecorder()
	login(w, req)
	RootIndex(w, req)

	assert.Equal(t, w.Code, 200)
	assert.Contains(t, w.Body.String(), "<a href=\"/\">Inici</a>")
}

func createPlayer(name string, ratings []int) {
	id, err := uuid.NewV4()
	if err != nil {
		panic("could not create uuid")
	}

	p := &Player{
		Id:         id.String(),
		Name:       name,
		Created_at: time.Now(),
	}
	Db.Insert(p)

	for _, v := range ratings {
		r := &Rating{
			Value:      v,
			Player_id:  id.String(),
			Created_at: time.Now(),
		}
		Db.Insert(r)
	}
}

func testCSV(t *testing.T, r *csv.Reader, fields int, values ...interface{}) {
	r.FieldsPerRecord = fields
	row, err := r.Read()
	if err != nil {
		t.Fatalf("oops...")
	}

	i := 0
	for _, v := range values {
		if !reflect.DeepEqual(v, row[i]) {
			t.Fatalf("Values are not the same")
		}
		i += 1
	}
}

func TestRootCsv(t *testing.T) {
	InitTest()
	defer CloseDB()

	// Let's create a couple of players and some ratings.
	createPlayer("one", []int{1, 2, 3})
	createPlayer("another", []int{0, 3, 8})
	createPlayer("notrated", []int{})

	// Perform the request.
	req, err := http.NewRequest("GET", "/", nil)
	assert.Nil(t, err)
	w := httptest.NewRecorder()
	RootCsv(w, req)

	// HTTP
	assert.Equal(t, w.Code, 200)
	header := w.Header()
	assert.Equal(t, header["Content-Type"][0], "text/csv")
	assert.Equal(t, header["Content-Disposition"][0],
		"attachment;filename=data.csv")

	// CSV
	dt := fmtDate(time.Now())
	r := csv.NewReader(w.Body)

	testCSV(t, r, 1, "another")
	testCSV(t, r, 3, "0", "8", "3.67")
	testCSV(t, r, 3, "0", "3", "8")
	testCSV(t, r, 3, dt, dt, dt)

	testCSV(t, r, 1, "notrated")
	testCSV(t, r, 1, "No ha puntuat cap entrenament")

	testCSV(t, r, 1, "one")
	testCSV(t, r, 3, "1", "3", "2.00")
	testCSV(t, r, 3, "1", "2", "3")
	testCSV(t, r, 3, dt, dt, dt)
}
