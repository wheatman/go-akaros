// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

#include "runtime.h"
#include "defs_GOOS_GOARCH.h"

void
runtime·netpollinit(void)
{
}

int32
runtime·netpollopen(uintptr fd, PollDesc *pd)
{
	USED(fd);
	USED(pd);
	return 0;
}

int32
runtime·netpollclose(uintptr fd)
{
	USED(fd);
	return 0;
}

// polls for ready network connections
// returns list of goroutines that become runnable
G*
runtime·netpoll(bool block)
{
	USED(block);
	return nil;
}
