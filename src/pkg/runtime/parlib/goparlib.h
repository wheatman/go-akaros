// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

#ifndef GOPARLIB_H
#include <parlib/akaros.h> // Eventually change to <parlib/GOOS.h>

extern void runtime∕parlib·Max_vcores(uint32 *n);
extern void runtime∕parlib·Futex(int32 *uaddr, int32 op, int32 val,
                                 const Timespec *timeout,
                                 int32 *uaddr2, int32 val3, int32 *ret);

extern void runtime∕parlib·Checkpoint(uint32 n);
#endif  // GOPARLIB_H
