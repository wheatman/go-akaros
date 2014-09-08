// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

#include "runtime.h"
#include "../cmd/ld/textflag.h"

// memset aligned words.
#pragma textflag NOSPLIT
static inline void *
memsetw(intgo* _v, intgo c, uintptr n)
{
    intgo *start, *end, *v;

    start = _v;
    end = _v + n/sizeof(intgo);
    v = _v;
    c = c & 0xff;
    c = c | c<<8;
    c = c | c<<16;
    #ifdef _64BIT
    c = c | c<<32;
    #endif

    while(v < end - (8-1))
    {
        v[3] = v[2] = v[1] = v[0] = c;
        v += 4;
        v[3] = v[2] = v[1] = v[0] = c;
        v += 4;
    }

    while(v < end)
      *v++ = c;

    return start;
}

#pragma textflag NOSPLIT
static inline void *
memset(byte* v, uint8 c, uintptr n)
{
    byte *p;
    uintgo n0;

	if (n == 0) return nil;

	p = v;

    while (n > 0 && ((uintptr)p & (sizeof(intgo)-1))) {
        *p++ = c;
        n--;
    }

    if (n >= sizeof(intgo)) {
        n0 = n / sizeof(intgo) * sizeof(intgo);
        memsetw((intgo*)p, c, n0);
        n -= n0;
        p += n0;
    }

    while (n > 0) {
        *p++ = c;
        n--;
    }

	return v;
}

// Runtime functions themselves
#pragma textflag NOSPLIT
void runtime·memclr(byte* p, uintptr n)
{
	memset(p, 0, n);
}

