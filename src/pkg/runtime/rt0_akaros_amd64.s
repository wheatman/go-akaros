// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

#include "zasm_GOOS_GOARCH.h"
#include "funcdata.h"
#include "../../cmd/ld/textflag.h"

/* In akaros we ALWAYS link using the cross compiler linker, so there is no
 * need to implement _rt0_GOARCH_akaros() as our entry point.  We do need it
 * defined however, to make gc happy.
 */
TEXT _rt0_amd64_akaros(SB),NOSPLIT,$0

/* The main function called out to from libc */
TEXT main(SB),NOSPLIT,$-8
	// copy arguments forward on an even stack
	MOVQ	DI, AX		// argc
	MOVQ	SI, BX		// argv
	SUBQ	$(4*8+7), SP		// 2args 2auto
	ANDQ	$~15, SP
	MOVQ	AX, 16(SP)
	MOVQ	BX, 24(SP)

	// find out information about the processor we're on
	MOVQ	$0, AX
	CPUID
	CMPQ	AX, $0
	JE	nocpuinfo
	MOVQ	$1, AX
	CPUID
	MOVL	CX, runtime·cpuid_ecx(SB)
	MOVL	DX, runtime·cpuid_edx(SB)
nocpuinfo:	
	
	// if there is an _cgo_init, call it.
	MOVQ	_cgo_init(SB), AX
	MOVQ	$runtime·g0(SB), DI
	MOVQ	$setmg_gcc<>(SB), SI
	CALL	AX
	// update stackguard after _cgo_init
	MOVQ	$runtime·g0(SB), CX
	MOVQ	g_stackguard0(CX), AX
	MOVQ	AX, g_stackguard(CX)

	// set the per-goroutine and per-mach "registers"
	get_tls(BX)
	LEAQ	runtime·g0(SB), CX
	MOVQ	CX, g(BX)
	LEAQ	runtime·m0(SB), AX
	MOVQ	AX, m(BX)

	// save m->g0 = g0
	MOVQ	CX, m_g0(AX)

	CLD				// convention is D is always left cleared
	CALL	runtime·check(SB)

	MOVL	16(SP), AX		// copy argc
	MOVL	AX, 0(SP)
	MOVQ	24(SP), AX		// copy argv
	MOVQ	AX, 8(SP)

	CALL	runtime·args(SB)
	CALL	runtime·osinit(SB)
	CALL	runtime·hashinit(SB)
	CALL	runtime·schedinit(SB)


	CALL	runtime·mstart(SB)
// TODO:
//	

	CALL	main·main(SB)		// entry
	MOVL	$0xf1, 0xf1  // crash



	// create a new goroutine to start program
	PUSHQ	$runtime·main·f(SB)		// entry
	PUSHQ	$0			// arg size
	ARGSIZE(16)
	CALL	runtime·newproc(SB)
	ARGSIZE(-1)
	POPQ	AX
	POPQ	AX

	// start this M
	CALL	runtime·mstart(SB)

	MOVL	$0xf1, 0xf1  // crash
	RET

// void setmg_gcc(M*, G*); set m and g called from gcc.
TEXT setmg_gcc<>(SB),NOSPLIT,$0
	get_tls(AX)
	MOVQ	DI, m(AX)
	MOVQ	SI, g(AX)
	RET

TEXT setmg_ken(SB),NOSPLIT,$0-16
	MOVQ	8(SP), DI
	MOVQ	16(SP), SI
	get_tls(AX)
	MOVQ	DI, m(AX)
	MOVQ	SI, g(AX)
	RET
