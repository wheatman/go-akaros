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
}

void
runtime·sigenable(uint32 sig)
{
	// Do nothing for now
	USED(sig);
}

void
runtime·sigdisable(uint32 sig)
{
	// Do nothing for now
	USED(sig);
}

void
runtime·resetcpuprofiler(int32 hz)
{
	// Do nothing for now
	USED(hz);
}

void
os·sigpipe(void)
{
	// Do nothing for now
}

void
runtime·crash(void)
{
	// Do nothing for now
}

