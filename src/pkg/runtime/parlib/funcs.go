// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// +build akaros

package parlib

/*
#cgo akaros LDFLAGS: -lpthread -lbenchutil -lm

#include <parlib.h>
#include <uthread.h>
#include <vcore.h>
#include <mcs.h>
#include <futex.h>
#include <sys/stat.h>

int go_parlib_errno(void)
{
	return current_uthread->err_no;
}

char *go_parlib_errstr(void)
{
	return current_uthread->err_str;
}

*/
import "C"
import (
	"unsafe"
	"errors"
)

var Procinfo *ProcinfoType = (*ProcinfoType)(unsafe.Pointer(uintptr(C.UINFO)))

func Futex(uaddr *int32, op int32, val int32,
           timeout *Timespec, uaddr2 *int32, val3 int32) (ret int32) {
	// For now, akaros futexes don't support uaddr2 or val3, so we
	// just 0 them out.
	uaddr2 = nil;
	val3 = 0;
	// Also, the minimum timout is 1ms, so up it to that if it's too small
	if (timeout != nil) {
		if (timeout.Sec == 0) {
			if (timeout.Nsec < 1000000) {
				timeout.Nsec = 1000000;
			}
		}
	}
	return int32(C.futex((*C.int)(unsafe.Pointer(uaddr)),
	                     C.int(op), C.int(val),
	                     (*C.struct_timespec)(unsafe.Pointer(timeout)),
	                     (*C.int)(unsafe.Pointer(uaddr2)), C.int(val3)))
}

func ProcinfoPackArgs(argv []*byte, envp []*byte) (pi ProcinfoType, err error) {
	p_pi := (*_Ctype_struct_procinfo)(unsafe.Pointer(&pi))
    p_argv := (**_Ctype_char)(unsafe.Pointer(&argv[0]))
    p_envp := (**_Ctype_char)(unsafe.Pointer(&envp[0]))

	__err := C.procinfo_pack_args(p_pi, p_argv, p_envp)
	if __err == -1 {
		err = nil
	} else {
		err = errors.New("ProcinfoPackArgs: error packing argv and envp")
	}
	return pi, err
}

func Errno() (int) {
	return int(C.go_parlib_errno())
}

func Errstr() (string) {
	return C.GoString(C.go_parlib_errstr())
}

func Chmod(path string, mode uint32) (int) {
	return int(C.chmod(C.CString(path), C.__mode_t(mode)))
}

func Fchmod(fd int, mode uint32) (int) {
	return int(C.fchmod(C.int(fd), C.__mode_t(mode)))
}

