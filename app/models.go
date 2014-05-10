// Copyright (C) 2014 Miquel Sabaté Solà <mikisabate@gmail.com>
// This file is licensed under the MIT license.
// See the LICENSE file.

package app

import (
	"time"
)

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
