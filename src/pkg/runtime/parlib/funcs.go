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
*/
import "C"
import "unsafe"

func Max_vcores(n *uint32) {
	*n = uint32(C.max_vcores())
}

func Futex(uaddr *int32, op int32, val int32,
           timeout *Timespec, uaddr2 *int32, val3 int32, ret *int32) {
	// For now, just ignore the timeout since Akaros doesn't support it
	timeout = nil
	__ret := int32(C.futex((*C.int)(unsafe.Pointer(uaddr)),
	                       C.int(op), C.int(val),
	                       (*C.struct_timespec)(unsafe.Pointer(timeout)),
	                       (*C.int)(unsafe.Pointer(uaddr2)), C.int(val3)))
	if ret != nil {
		*ret = __ret
	}
}
func Set_tls_desc(addr *uintptr, vcoreid int32) {
	C.set_tls_desc(unsafe.Pointer(addr), C.int(vcoreid))
}
