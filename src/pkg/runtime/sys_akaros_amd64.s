// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//
// System calls and other sys stuff for Akaros amd64
//

#include "zasm_GOOS_GOARCH.h"
#include "../cmd/ld/textflag.h"

// Do nothing for now
TEXT runtime·settls(SB), NOSPLIT, $0
	RET

TEXT sigtramp_real(SB),NOSPLIT,$40
    get_tls(BX)

    // check that g exists
    MOVQ    g(BX), R10
    CMPQ    R10, $0
    JNE     4(PC)
	// The sig_hand function is actually declared at the top of
	// runtime/parlib/signal.go inside of a cgo function! This is our normal
	// way of popping into go from the signals generated by parlib, so we might
	// as well reuse it...
    MOVQ    $sig_hand(SB), AX
    CALL    AX
    RET

    // save g
    MOVQ    R10, 32(SP)

    // g = m->gsignal
    MOVQ    g_m(R10), BP
    MOVQ    m_gsignal(BP), BP
    MOVQ    BP, g(BX)

    MOVQ    DI, 0(SP)
    MOVQ    SI, 8(SP)
    MOVQ    DX, 16(SP)
    MOVQ    R10, 24(SP)

    CALL    runtime·sighandler(SB)

    // restore g
    get_tls(BX)
    MOVQ    32(SP), R10
    MOVQ    R10, g(BX)
	RET

TEXT runtime·sigtramp(SB),NOSPLIT,$0
	// Follow the fucking calling convention!
	PUSHQ	BX
	PUSHQ	BP
	PUSHQ	R12
	PUSHQ	R13
	PUSHQ	R14
	PUSHQ	R15
	CALL	sigtramp_real(SB)
	POPQ	R15
	POPQ	R14
	POPQ	R13
	POPQ	R12
	POPQ	BP
	POPQ	BX
    RET

// This is the default handler assigned to our array of signal handlers in
// runtime/parlib/signal.go.  We want the default handler to be the same
// whether it is initiated by a user (i.e. a kill call), or by the kernel (i.e.
// a trap such as SIGSEGV occurred).  Defining this function and assigning as
// the default function called from go code ensures that this happens.
TEXT runtime∕parlib·defaultSighandler(SB),NOSPLIT,$0-8
    MOVQ    8(SP), DI
    CALL    sigtramp_real(SB)
	RET


