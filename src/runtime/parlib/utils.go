// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// +build akaros

package parlib

import (
	"errors"
)

// Get the strlen of a c style byte array
func Cstrlen(str []byte) int {
    var i int = 0
    for (str[i] != 0) { i++ }
    return i
}

// StringByteSlice is deprecated. Use ByteSliceFromString instead.
// If s contains a NUL byte this function panics instead of
// returning an error.
func StringByteSlice(s string) []byte {
	a, err := ByteSliceFromString(s)
	if err != nil {
		panic("syscall: string with NUL passed to StringByteSlice")
	}
	return a
}

// ByteSliceFromString returns a NUL-terminated slice of bytes
// containing the text of s. If s contains a NUL byte at any
// location, it returns (nil, error).
func ByteSliceFromString(s string) ([]byte, error) {
	for i := 0; i < len(s); i++ {
		if s[i] == 0 {
			return nil, errors.New("String contains a NULL byte.")
		}
	}
	a := make([]byte, len(s)+1)
	copy(a, s)
	return a, nil
}

// StringBytePtr is deprecated. Use BytePtrFromString instead.
// If s contains a NUL byte this function panics instead of
// returning an error.
func StringBytePtr(s string) *byte { return &StringByteSlice(s)[0] }

// BytePtrFromString returns a pointer to a NUL-terminated array of
// bytes containing the text of s. If s contains a NUL byte at any
// location, it returns (nil, EINVAL).
func BytePtrFromString(s string) (*byte, error) {
	a, err := ByteSliceFromString(s)
	if err != nil {
		return nil, err
	}
	return &a[0], nil
}

// StringSlicePtr is deprecated. Use SlicePtrFromStrings instead.
// If any string contains a NUL byte this function panics instead
// of returning an error.
func StringSlicePtr(ss []string) []*byte {
    bb := make([]*byte, len(ss)+1)
    for i := 0; i < len(ss); i++ {
        bb[i] = StringBytePtr(ss[i])
    }
    bb[len(ss)] = nil
    return bb
}

// SlicePtrFromStrings converts a slice of strings to a slice of
// pointers to NUL-terminated byte slices. If any string contains
// a NUL byte, it returns (nil, EINVAL).
func SlicePtrFromStrings(ss []string) ([]*byte, error) {
    var err error
    bb := make([]*byte, len(ss)+1)
    for i := 0; i < len(ss); i++ {
        bb[i], err = BytePtrFromString(ss[i])
        if err != nil {
            return nil, err
        }
    }
    bb[len(ss)] = nil
    return bb, nil
}

