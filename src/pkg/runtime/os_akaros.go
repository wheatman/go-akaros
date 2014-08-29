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

const stackSystem = 0
