// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// +build akaros

package parlib

/*
#include <parlib.h>
#include <uthread.h>
#include <vcore.h>
#include <mcs.h>
#include <futex.h>
#include <ros/memlayout.h>
*/
import "C"

const (
	FUTEX_WAIT = C.FUTEX_WAIT
	FUTEX_WAKE = C.FUTEX_WAKE
)
const (
	UINFO = C.UINFO
)
type Vcore C.struct_vcore
type Pcore C.struct_pcore
type Procinfo C.procinfo_t
type Timespec C.struct_timespec
type Timeval C.struct_timeval

