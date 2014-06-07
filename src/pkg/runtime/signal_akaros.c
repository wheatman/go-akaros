// Copyright 2012 The Go Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

#include "runtime.h"
#include "defs_GOOS_GOARCH.h"
#include "os_GOOS.h"
#include "signal_unix.h"

extern SigTab runtime·sigtab[];

void
runtime·initsig(void)
{
	int32 i;
	SigTab *t;

	// First call: basic setup.
	for(i = 0; i<NSIG; i++) {
		t = &runtime·sigtab[i];
		if((t->flags == 0) || (t->flags & SigDefault))
			continue;

		// For some signals, we respect an inherited SIG_IGN handler
		// rather than insist on installing our own default handler.
		// Even these signals can be fetched using the os/signal package.
		switch(i) {
		case SIGHUP:
		case SIGINT:
			if(runtime·getsig(i) == SIG_IGN) {
				t->flags = SigNotify | SigIgnored;
				continue;
			}
		}

		t->flags |= SigHandling;
		runtime·setsig(i, runtime·sighandler, true);
	}
}

void
runtime·sigenable(uint32 sig)
{
	SigTab *t;

	if(sig >= NSIG)
		return;

	t = &runtime·sigtab[sig];
	if((t->flags & SigNotify) && !(t->flags & SigHandling)) {
		t->flags |= SigHandling;
		if(runtime·getsig(sig) == SIG_IGN)
			t->flags |= SigIgnored;
		runtime·setsig(sig, runtime·sighandler, true);
	}
}

void
runtime·sigdisable(uint32 sig)
{
	SigTab *t;

	if(sig >= NSIG)
		return;

	t = &runtime·sigtab[sig];
	if((t->flags & SigNotify) && (t->flags & SigHandling)) {
		t->flags &= ~SigHandling;
		if(t->flags & SigIgnored)
			runtime·setsig(sig, SIG_IGN, true);
		else
			runtime·setsig(sig, SIG_DFL, true);
	}
}

void
runtime·resetcpuprofiler(int32 hz)
{
	if (hz == 0) {
		runtime·disable_profalarm();
	} else {
		runtime·enable_profalarm(1000000 / hz);
	}
	m->profilehz = hz;
	return;
}

void
os·sigpipe(void)
{
	runtime·setsig(SIGPIPE, SIG_DFL, false);
	runtime·raise(SIGPIPE);
}

void
runtime·crash(void)
{
	runtime·unblocksignals();
	runtime·setsig(SIGABRT, SIG_DFL, false);
	runtime·raise(SIGABRT);
}
