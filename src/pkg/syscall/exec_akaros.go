// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build akaros

package syscall

import (
	"unsafe"
	"runtime/parlib"
)

type SysProcAttr struct {
	Chroot     string      // Chroot.
	Credential *Credential // Credential.
	Ptrace     bool        // Enable tracing.
	Setsid     bool        // Create session.
	Setpgid    bool        // Set process group ID to new pid (SYSV setpgrp)
	Setctty    bool        // Set controlling terminal to fd Ctty (only meaningful if Setsid is set)
	Noctty     bool        // Detach fd 0 from controlling terminal
	Ctty       int         // Controlling TTY fd (Linux only)
	Pdeathsig  Signal      // Signal that the process will get when its parent dies (Linux only)
}

// Fork, dup fd onto 0..len(fd), and exec(argv0, argvv, envv) in child.
// If a dup or exec fails, write the errno error to pipe.
// (Pipe is close-on-exec so if exec succeeds, it will be closed.)
// In the child, this function must not acquire any locks, because
// they might have been locked at the time of the fork.  This means
// no rescheduling, no malloc calls, and no new stack segments.
// The calls to RawSyscall are okay because they are assembly
// functions that do not grow the stack.
func forkAndExecInChild(argv0 *byte, argv0len int, argv, envv []*byte, chroot, dir *byte, attr *ProcAttr, sys *SysProcAttr, pipe int) (pid int, err error) {
	// Declare all variables at top in case any
	// declarations require heap allocation (e.g., err1).
	var (
		r1     uintptr
		err1   error
	)
	// Make sure we aren't passing invalid arguments for Akaros (we should
	// probably support these some day though...)
	if chroot != nil {
		return 0, NewAkaError(EMORON, "Akaros does not support passing 'chroot' to forkAndExecInChild")
	}
	if dir != nil {
		return 0, NewAkaError(EMORON, "Akaros does not support passing 'dir' to forkAndExecInChild")
	}

	// Set up arguments for proc_create
	__cmd := uintptr(unsafe.Pointer(argv0))
	pi, _ := parlib.ProcinfoPackArgs(argv, envv)
    __cmdlen := uintptr(argv0len)
	__pi := uintptr(unsafe.Pointer(&pi))

	// Call proc create.
	r1, _, err1 = RawSyscall6(SYS_PROC_CREATE, __cmd, __cmdlen, __pi, parlib.PROC_DUP_FGRP, 0, 0)
	if err1 != nil {
		return 0, err1
	}
	child := int(r1)

	// Proc create succeeded, now run it!
	r1, _, err1 = RawSyscall(SYS_PROC_RUN, r1, 0, 0)
	if err1 != nil {
		return 0, err1
	}

	// Return the child pid
	return child, nil
}

// Try to open a pipe with O_CLOEXEC set on both file descriptors.
func forkExecPipe(p []int) (err error) {
	err = Pipe(p, O_CLOEXEC)
	return
}
