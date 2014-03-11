// Copyright 2013 The Go Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

#define SIG_REGS(ctxt) (*((HwTrapframe*)(&((((struct UserContext*)ctxt)->tf)[0]))))

#define SIG_RAX(info, ctxt) (SIG_REGS(ctxt).tf_rax)
#define SIG_RBX(info, ctxt) (SIG_REGS(ctxt).tf_rbx)
#define SIG_RCX(info, ctxt) (SIG_REGS(ctxt).tf_rcx)
#define SIG_RDX(info, ctxt) (SIG_REGS(ctxt).tf_rdx)
#define SIG_RDI(info, ctxt) (SIG_REGS(ctxt).tf_rdi)
#define SIG_RSI(info, ctxt) (SIG_REGS(ctxt).tf_rsi)
#define SIG_RBP(info, ctxt) (SIG_REGS(ctxt).tf_rbp)
#define SIG_RSP(info, ctxt) (SIG_REGS(ctxt).tf_rsp)
#define SIG_R8(info, ctxt) (SIG_REGS(ctxt).tf_r8)
#define SIG_R9(info, ctxt) (SIG_REGS(ctxt).tf_r9)
#define SIG_R10(info, ctxt) (SIG_REGS(ctxt).tf_r10)
#define SIG_R11(info, ctxt) (SIG_REGS(ctxt).tf_r11)
#define SIG_R12(info, ctxt) (SIG_REGS(ctxt).tf_r12)
#define SIG_R13(info, ctxt) (SIG_REGS(ctxt).tf_r13)
#define SIG_R14(info, ctxt) (SIG_REGS(ctxt).tf_r14)
#define SIG_R15(info, ctxt) (SIG_REGS(ctxt).tf_r15)
#define SIG_RIP(info, ctxt) (SIG_REGS(ctxt).tf_rip)
#define SIG_RFLAGS(info, ctxt) ((uint64)SIG_REGS(ctxt).tf_rflags)

#define SIG_CS(info, ctxt) ((uint64)SIG_REGS(ctxt).tf_cs)
#define SIG_FS(info, ctxt) ((uint64)SIG_REGS(ctxt).tf_fsbase)
#define SIG_GS(info, ctxt) ((uint64)SIG_REGS(ctxt).tf_gsbase)

#define SIG_CODE0(info, ctxt) ((info)->si_code)
#define SIG_CODE1(info, ctxt) (*((uintptr*)(&((info)->si_addr))))

