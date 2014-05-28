// Copyright (C) 2014 Miquel Sabaté Solà <mikisabate@gmail.com>
// This file is licensed under the MIT license.
// See the LICENSE file.

package app

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUserLogged(t *testing.T) {
	InitTest()
	defer CloseDB()

	req, err := http.NewRequest("GET", "/", nil)
	assert.Nil(t, err)

	assert.False(t, UserLogged(req, nil))

	s, err := store.Get(req, sessionName)
	assert.Nil(t, err)
	s.Values["userId"] = "1"
	w := httptest.NewRecorder()
	s.Save(req, w)

	assert.False(t, UserLogged(req, nil))

	createUser("user", "1234")
	var user User
	err = Db.SelectOne(&user, "select * from users")
	assert.Nil(t, err)

	s, err = store.Get(req, sessionName)
	assert.Nil(t, err)
	s.Values["userId"] = user.Id
	w = httptest.NewRecorder()
	s.Save(req, w)

	assert.True(t, UserLogged(req, nil))
}
