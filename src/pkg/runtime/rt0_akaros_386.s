// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

#include "../../cmd/ld/textflag.h"

/* In akaros we ALWAYS link using the cross compiler linker, so there is no
 * need to implement _rt0_GOARCH_akaros() as our entry point.  We do need it
 * defined however, to make gc happy.
 */
TEXT _rt0_386_akaros(SB),NOSPLIT,$0

/* The main function called out to from libc */
TEXT main(SB),NOSPLIT,$0
	JMP	runtime·rt0_go(SB)
