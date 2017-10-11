// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

#include "runtime.h"
#include "defs_GOOS_GOARCH.h"
#include "os_GOOS.h"
#include "signal_unix.h"
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
	runtime·ncpu = MAX(__procinfo.max_vcores, 1);
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
	// Initialize signal handling.
	runtime·unblocksignals();
}

// Called from dropm to undo the effect of an minit.
void
runtime·unminit(void)
{
	// Do nothing for now.
}

uintptr
runtime·memlimit(void)
{
	// Do nothing for now
	return 0;
}

/*
 * This assembler routine takes the args from registers, puts them on the stack,
 * and calls sighandler().
 */
#pragma cgo_import_static gcc_sigaction
typedef void (*gcc_call_t)(void *arg);
extern gcc_call_t gcc_sigaction;
extern void runtime·sigtramp(void);
extern SigTab runtime·sigtab[];
static Sigset sigset_none;
static Sigset sigset_all = { ~(uint32)0 };

void
runtime·setsig(int32 i, GoSighandler *fn, bool restart)
{
	USED(restart); // Akaros currently only supports the SA_SIGINFO flag
	Sigaction sa;
	runtime·memclr((byte*)&sa, sizeof sa);

	sa.sa_flags = SA_SIGINFO;
	if(fn == runtime·sighandler)
		fn = (void*)runtime·sigtramp;
	sa.sa_handler = fn;

	SigactionArg sarg;
	sarg.sig = i;
	sarg.act = &sa;
	sarg.oact = nil;
	runtime·asmcgocall(gcc_sigaction, &sarg);
	if (sarg.ret)
		runtime·throw("sigaction failure");
}

GoSighandler*
runtime·getsig(int32 i)
{
	Sigaction sa;
	runtime·memclr((byte*)&sa, sizeof sa);

	SigactionArg sarg;
	sarg.sig = i;
	sarg.act = nil;
	sarg.oact = &sa;
	runtime·asmcgocall(gcc_sigaction, &sarg);
	if (sarg.ret)
		runtime·throw("rt_sigaction read failure");

	if((void*)sa.sa_handler == runtime·sigtramp)
		return runtime·sighandler;
	return (void*)sa.sa_handler;
}

void
runtime·sigpanic(void)
{
	if(!runtime·canpanic(g))
		runtime·throw("unexpected signal during runtime execution");

	switch(g->sig) {
	case SIGBUS:
		if(g->sigcode0 == BUS_ADRERR && g->sigcode1 < 0x1000 || g->paniconfault) {
			if(g->sigpc == 0)
				runtime·panicstring("call of nil func value");
			runtime·panicstring("invalid memory address or nil pointer dereference");
		}
		runtime·printf("unexpected fault address %p\n", g->sigcode1);
		runtime·throw("fault");
	case SIGSEGV:
		if((g->sigcode0 == 0 || g->sigcode0 == SEGV_MAPERR || g->sigcode0 == SEGV_ACCERR) && g->sigcode1 < 0x1000 || g->paniconfault) {
			if(g->sigpc == 0)
				runtime·panicstring("call of nil func value");
			runtime·panicstring("invalid memory address or nil pointer dereference");
		}
		runtime·printf("unexpected fault address %p\n", g->sigcode1);
		runtime·throw("fault");
	case SIGFPE:
		switch(g->sigcode0) {
		case FPE_INTDIV:
			runtime·panicstring("integer divide by zero");
		case FPE_INTOVF:
			runtime·panicstring("integer overflow");
		}
		runtime·panicstring("floating point error");
	}
	runtime·panicstring(runtime·sigtab[g->sig].name);
}

void
runtime·unblocksignals(void)
{
	runtime·sigprocmask(SIG_SETMASK, &sigset_none, nil);
}
