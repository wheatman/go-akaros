#include "textflag.h"
#include "../runtime/zasm_akaros_amd64.h"

TEXT ·call(SB), NOSPLIT,$0-40
	CALL    runtime·entersyscall(SB)
	MOVQ    SP, R13 // save the stack pointer, so that it can be restored at the end
	MOVQ    SP, R14 // save the stack pointer, if we need to change stacks
	// first figure out what stack we will want to change to and do everything but change the stack
        // Figure out if we need to switch to m->g0 stack.
        get_tls(CX)
        MOVQ    g(CX), BP
        MOVQ    g_m(BP), BP
        MOVQ    m_g0(BP), SI
        MOVQ    g(CX), DI
        CMPQ    SI, DI
        JEQ     nosave
        MOVQ    m_gsignal(BP), SI
        CMPQ    SI, DI
        JEQ     nosave 
        MOVQ    m_g0(BP), SI
	MOVQ    (g_sched+gobuf_sp)(SI), R14
	// at this point r14 has the stack we will want to switch to
	// it will either be unchanged or the g0 stack
	// but we can unconditionally set SP to R14 after getting the arguments
nosave:
	MOVQ    8(SP), AX // the function pointer we want to call
	MOVQ    16(SP), R11 // the location of the slice which holds the arguments
	MOVQ    24(SP), R10 // the the number of arguments
	//MOVQ    32(SP), XX // the capacity of the slice
	//TODO this could be speed up with a PC reletive jump
	CMPQ	R10, $7
	JGE	stack_args
	CMPQ	R10, $6
	JE	six_args
	CMPQ	R10, $5
	JE	five_args
	CMPQ	R10, $4
	JE	four_args
	CMPQ	R10, $3
	JE	three_args
	CMPQ	R10, $2
	JE	two_args
	CMPQ	R10, $1
	JE	one_args
	CMPQ	R10, $0
	JE	zero_args
	JNE	0x0(PC)
stack_args:
	MOVQ	-8(R11)(R10*8), R9
	SUBQ	$8, R14  // 
	MOVQ	R9, 0x0(R14)
	SUBQ	$1, R10
	CMPQ	R10, $7
	JGE	stack_args

six_args:
	MOVQ	40(R11), R9
five_args:
	MOVQ	32(R11), R8
four_args:
	MOVQ	24(R11), CX
three_args:
	MOVQ	16(R11), DX
two_args:
	MOVQ	8(R11), SI
one_args:
	MOVQ	0x0(R11), DI
zero_args:
	MOVQ	R14, SP // change the stack pointer if needed
	CALL	AX
	MOVQ	R13, SP // restore the stack pointer
	// we need room for our 6 aruments and one extra for the return address
	MOVQ	AX, 40(SP)
	CALL    runtime·exitsyscall(SB)
	RET


TEXT ·call1(SB), NOSPLIT,$0-24
	CALL    runtime·entersyscall(SB)
	MOVQ    SP, R13 // save the stack pointer, so that it can be restored at the end
	MOVQ    SP, R14 // save the stack pointer, if we need to change stacks
	// first figure out what stack we will want to change to and do everything but change the stack
        // Figure out if we need to switch to m->g0 stack.
        get_tls(CX)
        MOVQ    g(CX), BP
        MOVQ    g_m(BP), BP
        MOVQ    m_g0(BP), SI
        MOVQ    g(CX), DI
        CMPQ    SI, DI
        JEQ     nosave1
        MOVQ    m_gsignal(BP), SI
        CMPQ    SI, DI
        JEQ     nosave1
        MOVQ    m_g0(BP), SI
	MOVQ    (g_sched+gobuf_sp)(SI), R14
	// at this point r14 has the stack we will want to switch to
	// it will either be unchanged or the g0 stack
	// but we can unconditionally set SP to R14 after getting the arguments
nosave1:
	MOVQ    8(SP), AX // the function pointer we want to call
	MOVQ    16(SP), DI // the argument
	MOVQ	R14, SP // change the stack pointer if needed
	CALL	AX
	MOVQ	R13, SP // restore the stack pointer
	// we need room for our 6 aruments and one extra for the return address
	MOVQ	AX, 24(SP)
	CALL    runtime·exitsyscall(SB)
	RET
