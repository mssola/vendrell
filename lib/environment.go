// Copyright (C) 2014 Miquel Sabaté Solà
// This file is licensed under the MIT license.
// See the LICENSE file.

package lib

import (
	"os"
	"path/filepath"
)

// Public: Returns the current environment in a string.
func Env() string {
	env := os.Getenv("VENDRELL_ENV")
	if env == "" {
		env = "development"
	}
	return env
}

// Public: get the first absolute path that has the "root" parameter as
// its root from the perspective of the "current" path. The current path
// can be relative (e.g. "." is an accepted value). The returned path
// has no trailing slashes. Note that if no path was found, then
// "/" will be returned.
//
// root    - The root we're looking for.
// current - The current path.
//
// Example:
//      FindRoot("mssola", "/home/mssola/another/mssola/dir")
//              will return "/home/mssola/another/mssola"
//
// Returns a string containing the absolute path matching the confitions.
func FindRoot(root, current string) string {
	current, _ = filepath.Abs(current)
	base := filepath.Base(current)

	for current != "/" && base != root {
		current = filepath.Dir(current)
		base = filepath.Base(current)
	}
	return current
}
