// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

#include "runtime.h"
#include "defs_GOOS_GOARCH.h"
#include "os_GOOS.h"
#include "../../cmd/ld/textflag.h"

void runtime·setldt(void)
{
	// Do nothing for now
}

#pragma textflag NOSPLIT
void
runtime·futexsleep(uint32 *addr, uint32 val, int64 ns)
{
	Timespec ts;
	
	if(ns < 0) {
		runtime·futex(addr, FUTEX_WAIT, val, nil, nil, 0);
		return;
	}
	// NOTE: tv_nsec is int64 on amd64, so this assumes a little-endian system.
	ts.tv_nsec = 0;
	ts.tv_sec = runtime·timediv(ns, 1000000000LL, (int32*)&ts.tv_nsec);
	runtime·futex(addr, FUTEX_WAIT, val, &ts, nil, 0);
}

#pragma textflag NOSPLIT
void
runtime·futexwakeup(uint32 *addr, uint32 cnt)
{
	int32 ret = runtime·futex(addr, FUTEX_WAKE, cnt, nil, nil, 0);
	if(ret >= 0)
		return;

	runtime·printf("futexwakeup addr=%p returned %d\n", addr, ret);
	runtime·throw("runtime.futexwakeup");
}

void
runtime·newosproc(M *mp, void *stk)
{
	// Unimplemented for now...
	USED(mp, stk);
	runtime·printf("runtime: failed to create new OS thread (have %d already)\n",
	               runtime·mcount());
	runtime·throw("runtime.newosproc");
}

void
runtime·osinit(void)
{
	runtime·ncpu = MIN(__procinfo.max_vcores, MAX_VCORES);
}

void
runtime·get_random_data(byte **rnd, int32 *rnd_len)
{
	// TODO: revisit and do something similar to Linux with #c/random
	*rnd = nil;
	*rnd_len = 0;
}

void
runtime·goenvs(void)
{
	runtime·goenvs_unix();
}

// Called to initialize a new m (including the bootstrap m).
// Called on the parent thread (main thread in case of bootstrap), can allocate memory.
void
runtime·mpreinit(M *mp)
{
	mp->gsignal = runtime·malg(32*1024);	// OS X wants >=8K, Akaros >=2K
}

// Called to initialize a new m (including the bootstrap m).
// Called on the new thread, can not allocate memory.
void
runtime·minit(void)
{
	// Do nothing for now.
}

// Called from dropm to undo the effect of an minit.
void
runtime·unminit(void)
{
	// Do nothing for now.
}

void
runtime·sigpanic(void)
{
	// Do nothing for now.
}

uintptr
runtime·memlimit(void)
{
	// Do nothing for now
	return 0;
}

