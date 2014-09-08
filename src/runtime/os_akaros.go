// Copyright 2018 The Go Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package runtime

import "unsafe"

func futex(addr unsafe.Pointer, op int32, val uint32, ts, addr2 unsafe.Pointer, val3 uint32) int32
func clone(flags int32, stk, mm, gg, fn unsafe.Pointer) int32
func sigaction(sig int32, new, old unsafe.Pointer)
func setitimer(mode int32, new, old unsafe.Pointer)
func sigprocmask(sig int32, new, old unsafe.Pointer) int32
func getrlimit(kind int32, limit unsafe.Pointer) int32
func netpollinit()
func netpollopen(fd uintptr, pd *pollDesc) int32
func netpollclose(fd uintptr) int32
func netpollarm(pd *pollDesc, mode int)

const stackSystem = 0

func os_sigpipe() {
	gothrow("too many writes on closed pipe")
}

//TODO(wheatman)this is minimal just to get it compiling
func sigpanic() {
        gothrow("fault");
}

