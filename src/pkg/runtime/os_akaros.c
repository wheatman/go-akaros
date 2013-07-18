// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

#include "runtime.h"
#include "defs_GOOS_GOARCH.h"
#include "os_GOOS.h"

void runtime·setldt(void)
{
	// Do nothing for now
}

void
runtime·futexsleep(uint32 *addr, uint32 val, int64 ns)
{
	Timespec ts, *tsp;
	int64 secs;

	if(ns < 0)
		tsp = nil;
	else {
		secs = ns/1000000000LL;
		// Avoid overflow
		if(secs > 1LL<<30)
			secs = 1LL<<30;
		ts.tv_sec = secs;
		ts.tv_nsec = ns%1000000000LL;
		tsp = &ts;
	}
	runtime∕parlib·Futex((int32*)addr, FUTEX_WAIT, val, tsp, nil, 0, nil);
}

void
runtime·futexwakeup(uint32 *addr, uint32 cnt)
{
	int32 ret;
	runtime∕parlib·Futex((int32*)addr, FUTEX_WAKE, cnt, nil, nil, 0, &ret);
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
	uint32 n;
	runtime∕parlib·Max_vcores(&n);
	runtime·ncpu = n;
}

void
runtime·get_random_data(byte **rnd, int32 *rnd_len)
{
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

void
runtime·setprof(bool on)
{
	USED(on);
}

#pragma dataflag 16  // no pointers
static int8 badsignal[] = "runtime: signal received on thread not created by Go!\n";

// This runs on a foreign stack, without an m or a g.  No stack split.
#pragma textflag 7
void
runtime·badsignal(int32 sig)
{
	// Think of a better way to do this with a symbol table to print the actual
	// name of the event received
	USED(sig);
	runtime·write(2, badsignal, sizeof badsignal - 1);
	runtime·exit(1);
}

