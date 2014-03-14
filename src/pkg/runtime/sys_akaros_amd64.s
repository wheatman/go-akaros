// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//
// System calls and other sys stuff for Akaros amd64
//

#include "zasm_GOOS_GOARCH.h"
#include "../../cmd/ld/textflag.h"

// Do nothing for now
TEXT runtime·settls(SB), NOSPLIT, $0
	RET

TEXT sigtramp_real(SB),NOSPLIT,$64
    get_tls(BX)

    // check that m exists
    MOVQ    m(BX), BP
    CMPQ    BP, $0
    JNE     4(PC)
    MOVQ    $sig_hand(SB), AX
    CALL    AX
    RET

    // save g
    MOVQ    g(BX), R10
    MOVQ    R10, 40(SP)

    // g = m->gsignal
    MOVQ    m_gsignal(BP), BP
    MOVQ    BP, g(BX)

    MOVQ    DI, 0(SP)
    MOVQ    SI, 8(SP)
    MOVQ    DX, 16(SP)
    MOVQ    R10, 24(SP)

    CALL    runtime·sighandler(SB)

    // restore g
    get_tls(BX)
    MOVQ    40(SP), R10
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

TEXT runtime∕parlib·defaultSighandler(SB),NOSPLIT,$0
    MOVQ    8(SP), DI
    CALL    sigtramp_real(SB)
	RET


