// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package syscall

import (
	"runtime"
	"runtime/parlib"
)

func RunWithDeadline(f func(), deadline int64) {
	if (deadline <= 0) {
		f()
	} else {
		var waiter parlib.AlarmWaiter
		runtime.LockOSThread()
	    parlib.AbortSyscallAt(&waiter, deadline)
		f()
		parlib.CancelAbortSyscall(&waiter)
		runtime.UnlockOSThread()
	}
}
