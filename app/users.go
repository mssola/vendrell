// Copyright (C) 2014 Miquel Sabaté Solà <mikisabate@gmail.com>
// This file is licensed under the MIT license.
// See the LICENSE file.

package app

import (
	"net/http"
	"time"

	"github.com/coopernurse/gorp"
	"github.com/mssola/go-utils/security"
	"github.com/nu7hatch/gouuid"
)

func UsersCreate(res http.ResponseWriter, req *http.Request, db gorp.DbMap) {
	// Only one user is allowed in this application.
	count, err := db.SelectInt("select count(*) from users")
	if err != nil || count > 0 {
		http.Redirect(res, req, "/", http.StatusFound)
		return
	}

	// Create the user and redirect.
	uuid, err := uuid.NewV4()
	if err != nil {
		http.Redirect(res, req, "/", http.StatusFound)
		return
	}
	u := &User{
		Id:            uuid.String(),
		Name:          req.FormValue("name"),
		Password_hash: security.PasswordSalt(req.FormValue("password")),
		Created_at:    time.Now(),
	}
	db.Insert(u)
	http.Redirect(res, req, "/", http.StatusFound)
}
