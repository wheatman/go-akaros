// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build akaros
// we're always root on akaros for now.
package user

func current() (*User, error) {
	return lookupAkaros(0, "", false)
}

func lookup(username string) (*User, error) {
	return lookupAkaros(-1, username, true)
}

func lookupId(uid string) (*User, error) {
	return lookupAkaros(0, "", false)
}

func lookupAkaros(uid int, username string, lookupByName bool) (*User, error) {
	u := &User{
		Uid:      "0",
		Gid:      "0",
		Username: "root",
		Name:     "nanwan",
		HomeDir:  "/",
	}
	return u, nil
}
