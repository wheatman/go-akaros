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
#include "gcc_akaros.h"
*/
import "C"

const (
	MAX_VCORES = C.MAX_VCORES
)
const (
	FUTEX_WAIT = C.FUTEX_WAIT
	FUTEX_WAKE = C.FUTEX_WAKE
)
type Vcore C.struct_vcore
type Pcore C.struct_pcore
type ProcinfoType C.procinfo_t
type Timespec C.struct_timespec
type Timeval C.struct_timeval
type Ucq C.struct_ucq
type EventQueue C.struct_event_queue
type EventMbox C.struct_event_mbox
type SyscallArg C.gcc_syscall_arg_t
type FutexArg C.gcc_futex_arg_t

