// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package syscall

import (
	"runtime"
	"usys"
)

func RunWithDeadline(f func(), deadline int64) {
	if deadline <= 0 {
		f()
	} else {
		runtime.LockOSThread()
		handle := usys.Call1(usys.USYS_ABORT_SYSCALL_AT_ABS_UNIX, uintptr(deadline))
		f()
		usys.Call1(usys.USYS_UNSET_ALARM, uintptr(handle))
		runtime.UnlockOSThread()
	}
}
