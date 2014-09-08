// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// StartProcess and Exec

package syscall

import (
	"unsafe"
	"runtime/parlib"
)

// ProcAttr holds attributes that will be applied to a new process started
// by StartProcess.
type ProcAttr struct {
	Dir   string    // Current working directory.
	Env   []string  // Environment.
	Files []uintptr // File descriptors.
	Sys   *SysProcAttr // System specific attrs
}
var zeroProcAttr ProcAttr

// Undefined on Akaros
type SysProcAttr struct {}
var zeroSysProcAttr SysProcAttr

// SlicePtrFromStrings converts a slice of strings to a slice of
// pointers to NUL-terminated byte slices. If any string contains
// a NUL byte, it returns (nil, EINVAL).
func SlicePtrFromStrings(ss []string) ([]*byte, error) {
	var err error
	bb := make([]*byte, len(ss)+1)
	for i := 0; i < len(ss); i++ {
		bb[i], err = BytePtrFromString(ss[i])
		if err != nil {
			return nil, err
		}
	}
	bb[len(ss)] = nil
	return bb, nil
}

func StartProcess(argv0 string, argv []string, attr *ProcAttr) (pid int, handle uintptr, err error) {
	if attr == nil {
		attr = &zeroProcAttr
	}
	sys := attr.Sys
	if sys == nil {
		sys = &zeroSysProcAttr
	}

	// Convert args to C form.
	argv0p, err := ByteSliceFromString(argv0)
	if err != nil {
		return 0, 0, err
	}
	argvp, err := SlicePtrFromStrings(argv)
	if err != nil {
		return 0, 0, err
	}
	envvp, err := SlicePtrFromStrings(attr.Env)
	if err != nil {
		return 0, 0, err
	}

	var dir []byte
	if attr.Dir != "" {
		dir, err = ByteSliceFromString(attr.Dir)
		if err != nil {
			return 0, 0, err
		}
	}

	// Kick off child.
	pid, err = startProcess(argv0p, argvp, envvp, dir, attr.Files)

	// Return the pid and the error if there was one
	return pid, 0, err
}

func startProcess(argv0 []byte, argv, envv []*byte, dir []byte, files []uintptr) (pid int, err error) {
	var r1 uintptr

	// Adjust argv0 to prepend 'dir' if argv0 is a relative path
	if argv0[0] != '/' {
		if len(dir) > 0 {
			argv0 = append(dir[:len(dir)-1], append([]byte{'/'}, argv0...)...)
		}
	}

	// Call proc create to create a child.
	cmd := uintptr(unsafe.Pointer(&argv0[0]))
	cmdlen := uintptr(len(argv0))
	sd, err := parlib.SerializeArgvEnvp(argv, envv)
	if err != nil {
		return 0, err
	}
	sdbuf := uintptr(unsafe.Pointer(&sd.Buf[0]))
	sdlen := uintptr(sd.Len)
	r1, _, err = RawSyscall6(SYS_PROC_CREATE, cmd, cmdlen, sdbuf, sdlen, 0, 0)
	parlib.FreeSerializedData(sd)
	if err != nil {
		return 0, err
	}
	child := r1

	// Dup the fd map properly into the child
	__cfdm := make([]Childfdmap_t, len(files))
	for i, f := range(files) {
		__cfdm[i].Parentfd = uint32(f)
		__cfdm[i].Childfd = uint32(i)
		__cfdm[i].Ok = int32(-1)
	}
	cfdm := uintptr(unsafe.Pointer(&__cfdm[0]))
	cfdmlen := uintptr(len(__cfdm))
	r1, _, err = RawSyscall(SYS_DUP_FDS_TO, child, cfdm, cfdmlen)
	if err != nil {
		return 0, err
	}

	// If 'dir' passed in, set the pwd of the child
	if len(dir) > 0 {
		pwd := uintptr(unsafe.Pointer(&dir[0]))
		pwdlen := uintptr(len(dir))
		r1, _, err = RawSyscall(SYS_CHDIR, child, pwd, pwdlen)
		if err != nil {
			return 0, err
		}
	}

	// Now run the child!
	r1, _, err = RawSyscall(SYS_PROC_RUN, child, 0, 0)
	if err != nil {
		return 0, err
	}

	// Return the child pid
	return int(child), nil
}

// Ordinary exec.
func Exec(argv0 string, argv []string, envv []string) (err error) {
	// Convert args to C form.
	argv0p, err := ByteSliceFromString(argv0)
	if err != nil {
		return err
	}
	argvp, err := SlicePtrFromStrings(argv)
	if err != nil {
		return err
	}
	envvp, err := SlicePtrFromStrings(envv)
	if err != nil {
		return err
	}

	// exec to new cmd
	cmd := uintptr(unsafe.Pointer(&argv0p[0]))
	cmdlen := uintptr(len(argv0))
	sd, err := parlib.SerializeArgvEnvp(argvp, envvp)
	if err != nil {
		return err
	}
	sdbuf := uintptr(unsafe.Pointer(&sd.Buf[0]))
	sdlen := uintptr(sd.Len)
	_, _, err = RawSyscall6(SYS_EXEC, cmd, cmdlen, sdbuf, sdlen, 0, 0)
	parlib.FreeSerializedData(sd)
	return err
}

func CloseOnExec(fd int) { fcntl(fd, F_SETFD, FD_CLOEXEC) }
