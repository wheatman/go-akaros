// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

#include <parlib/ztypes_akaros.h>

// I wanted to pass UINFO through cgo -cdefs, but it turns all #defines into
// enums, which can only be 32-bit
#ifdef _64BIT
//#define UINFO 0x7f7fffe00000ULL
static const uint64 UINFO = 0x7f7fffe00000ULL;
#else
//#define UINFO 0x7f800000
static const uint32 UINFO = 0x7f800000UL;
#endif

#define MIN(a, b)	((a < b) ? a : b)
#define MAX(a, b)	((a > b) ? a : b)

#define __procinfo (*(ProcinfoType*)UINFO)
