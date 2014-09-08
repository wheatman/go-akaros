// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// +build akaros

package parlib

/*
#cgo akaros LDFLAGS: -lpthread -lbenchutil -lm

#define _LARGEFILE64_SOURCE

#include <futex.h>
#include <unistd.h>
#include <malloc.h>
#include <parlib/parlib.h>
#include <parlib/uthread.h>
#include <parlib/vcore.h>
#include <parlib/mcs.h>
#include <parlib/serialize.h>
#include <sys/stat.h>
#include <sys/types.h>

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

func SerializeArgvEnvp(argv []*byte, envp []*byte) (sd *SerializedData, err error) {
    p_argv := (**_Ctype_char)(unsafe.Pointer(&argv[0]))
    p_envp := (**_Ctype_char)(unsafe.Pointer(&envp[0]))

	__sd := C.serialize_argv_envp(p_argv, p_envp)
	if __sd == nil {
		err = errors.New("SerializeArgvEnvp: error packing argv and envp")
	} else {
		sd = (*SerializedData)(unsafe.Pointer(__sd))
	}
	return sd, err
}

func FreeSerializedData(sd *SerializedData) {
	C.free(unsafe.Pointer(sd))
}

func Errno() (int) {
	return int(C.go_parlib_errno())
}

func Errstr() (string) {
	return C.GoString(C.go_parlib_errstr())
}
