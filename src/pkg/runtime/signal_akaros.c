// Copyright 2012 The Go Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

#include "runtime.h"
#include "defs_GOOS_GOARCH.h"
#include "os_GOOS.h"

void
runtime·initsig(void)
{
	// Do nothing for now
	runtime·printf("In runtime.initsig\n");
}

void
runtime·sigenable(uint32 sig)
{
	// Do nothing for now
	runtime·printf("In runtime.sigenable: %d\n", sig);
	USED(sig);
}

void
runtime·sigdisable(uint32 sig)
{
	// Do nothing for now
	runtime·printf("In runtime.sigdisable: %d\n", sig);
	USED(sig);
}

void
runtime·resetcpuprofiler(int32 hz)
{
	// Do nothing for now
	runtime·printf("In runtime.resetcpuprofiler: %d\n", hz);
	USED(hz);
}

void
os·sigpipe(void)
{
	runtime·printf("In os.pipe\n");
	// Do nothing for now
}

void
runtime·crash(void)
{
	// Do nothing for now
	runtime·printf("In runtime.crash\n");
}

