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
func forkAndExecInChild(argv0 []byte, argv, envv []*byte, chroot, dir []byte, attr *ProcAttr, sys *SysProcAttr) (pid int, err error) {
	// Declare all variables at top in case any
	// declarations require heap allocation (e.g., err1).
	var (
		r1     uintptr
		err1   error
	)

	// Make sure we aren't passing invalid arguments for Akaros (we should
	// probably support these some day though...)
	if len(chroot) > 0 {
		return 0, NewAkaError(EMORON, "Akaros does not support passing 'chroot' to forkAndExecInChild")
	}

	// Adjust argv0 to prepend 'dir' if argv0 is a relative path
	if argv0[0] != '/' {
		if len(dir) > 0 {
			argv0 = append(dir[:len(dir)-1], append([]byte{'/'}, argv0...)...)
		}
	}

	// Call proc create to create a child.
	cmd := uintptr(unsafe.Pointer(&argv0[0]))
	cmdlen := uintptr(len(argv0))
	__pi, _ := parlib.ProcinfoPackArgs(argv, envv)
	pi := uintptr(unsafe.Pointer(&__pi))
	r1, _, err1 = RawSyscall6(SYS_PROC_CREATE, cmd, cmdlen, pi, 0, 0, 0)
	if err1 != nil {
		return 0, err1
	}
	child := r1

	// Dup the fd map properly into the child
	__cfdm := make([]Childfdmap_t, len(attr.Files))
	for i, f := range(attr.Files) {
		__cfdm[i].Parentfd = uint32(f)
		__cfdm[i].Childfd = uint32(i)
		__cfdm[i].Ok = int32(-1)
	}
	cfdm := uintptr(unsafe.Pointer(&__cfdm[0]))
	cfdmlen := uintptr(len(__cfdm))
	r1, _, err1 = RawSyscall(SYS_DUP_FDS_TO, child, cfdm, cfdmlen)
	if err1 != nil {
		return 0, err1
	}

	// If 'dir' passed in, set the pwd of the child
	if len(dir) > 0 {
		pwd := uintptr(unsafe.Pointer(&dir[0]))
		pwdlen := uintptr(len(dir))
		r1, _, err1 = RawSyscall(SYS_CHDIR, child, pwd, pwdlen)
		if err1 != nil {
			return 0, err1
		}
	}

	// Now run the child!
	r1, _, err1 = RawSyscall(SYS_PROC_RUN, child, 0, 0)
	if err1 != nil {
		return 0, err1
	}

	// Return the child pid
	return int(child), nil
}

// Try to open a pipe with O_CLOEXEC set on both file descriptors.
func forkExecPipe(p []int) (err error) {
	err = Pipe(p, O_CLOEXEC)
	return
}
