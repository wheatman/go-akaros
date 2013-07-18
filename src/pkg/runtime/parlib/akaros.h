// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

#ifndef PARLIB_AKAROS_H
#include <parlib/syscall_akaros.h>
#include <parlib/types_akaros.h>

#define __procinfo (*(Procinfo*)UINFO)
#define akaros_syscall(num, a0, a1, a2, a3, a4, a5, perrno, pret) \
	runtime∕parlib·Syscall((uint32)(num),                         \
	                       (intgo)(a0), (intgo)(a1),              \
	                       (intgo)(a2), (intgo)(a3),              \
	                       (intgo)(a4), (intgo)(a5),              \
	                       (perrno), (pret))
extern void runtime∕parlib·Syscall(uint32 _num,
                                   intgo _a0, intgo _a1,
                                   intgo _a2, intgo _a3,
                                   intgo _a4, intgo _a5,
                                   int32 *errno_loc, intgo* ret);
#endif  // PARLIB_AKAROS_H
