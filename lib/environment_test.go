// Copyright (C) 2014 Miquel Sabaté Solà
// This file is licensed under the MIT license.
// See the LICENSE file.

package lib

import (
	"os"
	"testing"
)

func TestEnv(t *testing.T) {
	os.Clearenv()
	if env := Env(); env != "development" {
		t.Errorf("Expected 'development'")
	}
	os.Setenv("VENDRELL_ENV", "production")
	if env := Env(); env != "production" {
		t.Errorf("Expected 'production'")
	}
	os.Setenv("VENDRELL_ENV", "")
	if env := Env(); env != "development" {
		t.Errorf("Expected 'development'")
	}
}

func TestFindRoot(t *testing.T) {
	abs := FindRoot("mssola", "/home/mssola/lala")
	if abs != "/home/mssola" {
		t.Errorf("Expected '/home/mssola'")
	}
	abs = FindRoot("home", "/home/mssola/lala")
	if abs != "/home" {
		t.Errorf("Expected '/home'")
	}
	abs = FindRoot("/", "/home/mssola/lala")
	if abs != "/" {
		t.Errorf("Expected '/'")
	}
}
