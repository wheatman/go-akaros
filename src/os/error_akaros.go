// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package os

import "syscall"

func isExist(err error) bool {
	switch pe := err.(type) {
	case nil:
		return false
	case *PathError:
		err = pe.Err
	case *LinkError:
		err = pe.Err
	}
	switch pe := err.(type) {
		case *syscall.AkaError:
			return pe.Errno() == syscall.EEXIST
		default:
	       return pe == ErrExist
	}
}

func isNotExist(err error) bool {
	switch pe := err.(type) {
	case nil:
		return false
	case *PathError:
		err = pe.Err
	case *LinkError:
		err = pe.Err
	}
	switch pe := err.(type) {
		case *syscall.AkaError:
			return pe.Errno() == syscall.ENOENT
		default:
	       return pe == ErrNotExist
	}
}

func isPermission(err error) bool {
	switch pe := err.(type) {
	case nil:
		return false
	case *PathError:
		err = pe.Err
	case *LinkError:
		err = pe.Err
	}
	switch pe := err.(type) {
		case *syscall.AkaError:
			return pe.Errno() == syscall.EACCES || pe.Errno() == syscall.EPERM
		default:
	       return pe == ErrPermission
	}
}

