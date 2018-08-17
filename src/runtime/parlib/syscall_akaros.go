// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package parlib

/*
#include <parlib/alarm.h>
#include <parlib/uthread.h>
#include <sys/syscall.h>

*/
import "C"

type SyscallType C.struct_syscall

const (
	MAX_ERRSTR_LEN = C.MAX_ERRSTR_LEN
)
