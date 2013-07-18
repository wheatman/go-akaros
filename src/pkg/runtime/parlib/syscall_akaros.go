// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package parlib

/*
#include <sys/syscall.h>
*/
import "C"
import "unsafe"

func Syscall(_num uint32,
             _a0 int, _a1 int,
             _a2 int, _a3 int,
             _a4 int, _a5 int,
             errno_loc *int32, ret *int) {
	__ret := int(C.__ros_syscall(C.uint(_num),
	                             C.long(_a0), C.long(_a1),
	                             C.long(_a2), C.long(_a3),
	                             C.long(_a4), C.long(_a5),
	                             (*C.int)(unsafe.Pointer(errno_loc))))
	if ret != nil {
		*ret = __ret
	}
}

