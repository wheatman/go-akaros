// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

#include <sys/syscall.h>
#include <time.h>

typedef struct syscall gcc_syscall_arg_t;

typedef struct gcc_futex_arg {
	int *uaddr;
	int op;
	int val;
	struct timespec *timeout;
	int *uaddr2;
	int val3;
	int retval;	
} gcc_futex_arg_t;

