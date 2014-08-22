// Copyright 2013 The Go Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

#define ROS_HW_CTX              1                                                        |#endif
#define ROS_SW_CTX              2

#define SIG_CODE0(info, ctxt) ((info)->si_code)
#define SIG_CODE1(info, ctxt) (*((uintptr*)(&((info)->si_addr))))

// Registers in both the HW and SW trapframes
// Safe to use as an l-value
#define SIG_GS(info, ctxt) \
	(*((uint64*)((((struct UserContext*)(ctxt))->type == 1) ? \
		&(((HwTrapframe*)(&((((UserContext*)(ctxt))->tf)[0])))->tf_gsbase) : \
		&(((SwTrapframe*)(&((((UserContext*)(ctxt))->tf)[0])))->tf_gsbase))))
#define SIG_FS(info, ctxt) \
	(*((uint64*)((((struct UserContext*)(ctxt))->type == 1) ? \
		&(((HwTrapframe*)(&((((UserContext*)(ctxt))->tf)[0])))->tf_fsbase) : \
		&(((SwTrapframe*)(&((((UserContext*)(ctxt))->tf)[0])))->tf_fsbase))))
#define SIG_RBX(info, ctxt) \
	(*((uint64*)((((struct UserContext*)(ctxt))->type == 1) ? \
		&(((HwTrapframe*)(&((((UserContext*)(ctxt))->tf)[0])))->tf_rbx) : \
		&(((SwTrapframe*)(&((((UserContext*)(ctxt))->tf)[0])))->tf_rbx))))
#define SIG_RBP(info, ctxt) \
	(*((uint64*)((((struct UserContext*)(ctxt))->type == 1) ? \
		&(((HwTrapframe*)(&((((UserContext*)(ctxt))->tf)[0])))->tf_rbp) : \
		&(((SwTrapframe*)(&((((UserContext*)(ctxt))->tf)[0])))->tf_rbp))))
#define SIG_R12(info, ctxt) \
	(*((uint64*)((((struct UserContext*)(ctxt))->type == 1) ? \
		&(((HwTrapframe*)(&((((UserContext*)(ctxt))->tf)[0])))->tf_r12) : \
		&(((SwTrapframe*)(&((((UserContext*)(ctxt))->tf)[0])))->tf_r12))))
#define SIG_R13(info, ctxt) \
	(*((uint64*)((((struct UserContext*)(ctxt))->type == 1) ? \
		&(((HwTrapframe*)(&((((UserContext*)(ctxt))->tf)[0])))->tf_r13) : \
		&(((SwTrapframe*)(&((((UserContext*)(ctxt))->tf)[0])))->tf_r13))))
#define SIG_R14(info, ctxt) \
	(*((uint64*)((((struct UserContext*)(ctxt))->type == 1) ? \
		&(((HwTrapframe*)(&((((UserContext*)(ctxt))->tf)[0])))->tf_r14) : \
		&(((SwTrapframe*)(&((((UserContext*)(ctxt))->tf)[0])))->tf_r14))))
#define SIG_R15(info, ctxt) \
	(*((uint64*)((((struct UserContext*)(ctxt))->type == 1) ? \
		&(((HwTrapframe*)(&((((UserContext*)(ctxt))->tf)[0])))->tf_r15) : \
		&(((SwTrapframe*)(&((((UserContext*)(ctxt))->tf)[0])))->tf_r15))))
#define SIG_RIP(info, ctxt) \
	(*((uint64*)((((struct UserContext*)(ctxt))->type == 1) ? \
		&(((HwTrapframe*)(&((((UserContext*)(ctxt))->tf)[0])))->tf_rip) : \
		&(((SwTrapframe*)(&((((UserContext*)(ctxt))->tf)[0])))->tf_rip))))
#define SIG_RSP(info, ctxt) \
	(*((uint64*)((((struct UserContext*)(ctxt))->type == 1) ? \
		&(((HwTrapframe*)(&((((UserContext*)(ctxt))->tf)[0])))->tf_rsp) : \
		&(((SwTrapframe*)(&((((UserContext*)(ctxt))->tf)[0])))->tf_rsp))))

// Registers only in the HW trapframe (0'd out if dealing with sw trapframe)
// Can't be used as an l-value
#define SIG_RAX(info, ctxt) \
	((((struct UserContext*)(ctxt))->type == 1) ? \
		((HwTrapframe*)(&((((UserContext*)(ctxt))->tf)[0])))->tf_rax : 0)
#define SIG_RCX(info, ctxt) \
	((((struct UserContext*)(ctxt))->type == 1) ? \
		((HwTrapframe*)(&((((UserContext*)(ctxt))->tf)[0])))->tf_rcx : 0)
#define SIG_RDX(info, ctxt) \
	((((struct UserContext*)(ctxt))->type == 1) ? \
		((HwTrapframe*)(&((((UserContext*)(ctxt))->tf)[0])))->tf_rdx : 0)
#define SIG_RDI(info, ctxt) \
	((((struct UserContext*)(ctxt))->type == 1) ? \
		((HwTrapframe*)(&((((UserContext*)(ctxt))->tf)[0])))->tf_rdi : 0)
#define SIG_RSI(info, ctxt) \
	((((struct UserContext*)(ctxt))->type == 1) ? \
		((HwTrapframe*)(&((((UserContext*)(ctxt))->tf)[0])))->tf_rsi : 0)
#define SIG_R8(info, ctxt) \
	((((struct UserContext*)(ctxt))->type == 1) ? \
		((HwTrapframe*)(&((((UserContext*)(ctxt))->tf)[0])))->tf_r8 : 0)
#define SIG_R9(info, ctxt) \
	((((struct UserContext*)(ctxt))->type == 1) ? \
		((HwTrapframe*)(&((((UserContext*)(ctxt))->tf)[0])))->tf_r9 : 0)
#define SIG_R10(info, ctxt) \
	((((struct UserContext*)(ctxt))->type == 1) ? \
		((HwTrapframe*)(&((((UserContext*)(ctxt))->tf)[0])))->tf_r10 : 0)
#define SIG_R11(info, ctxt) \
	((((struct UserContext*)(ctxt))->type == 1) ? \
		((HwTrapframe*)(&((((UserContext*)(ctxt))->tf)[0])))->tf_r11 : 0)
#define SIG_CS(info, ctxt) \
	((((struct UserContext*)(ctxt))->type == 1) ? \
		((HwTrapframe*)(&((((UserContext*)(ctxt))->tf)[0])))->tf_cs : 0)
#define SIG_RFLAGS(info, ctxt) \
	((((struct UserContext*)(ctxt))->type == 1) ? \
		((HwTrapframe*)(&((((UserContext*)(ctxt))->tf)[0])))->tf_rflags : 0)

