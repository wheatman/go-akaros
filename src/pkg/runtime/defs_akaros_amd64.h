// Created by cgo -cdefs - DO NOT EDIT
// cgo -cdefs defs_linux.go defs1_linux.go

#include "parlib/goparlib.h"

enum {
	EINTR	= 0x4,
	EAGAIN	= 0xb,
	ENOMEM	= 0xc,

	PROT_NONE	= 0x0,
	PROT_READ	= 0x1,
	PROT_WRITE	= 0x2,
	PROT_EXEC	= 0x4,

	MAP_ANON	= 0x20,
	MAP_PRIVATE	= 0x2,
	MAP_FIXED	= 0x10,
	MAP_POPULATE = 0x8000,

	MADV_DONTNEED	= 0x4,

	FPE_INTDIV	= 0x1,
	FPE_INTOVF	= 0x2,
	FPE_FLTDIV	= 0x3,
	FPE_FLTOVF	= 0x4,
	FPE_FLTUND	= 0x5,
	FPE_FLTRES	= 0x6,
	FPE_FLTINV	= 0x7,
	FPE_FLTSUB	= 0x8,

	BUS_ADRALN	= 0x1,
	BUS_ADRERR	= 0x2,
	BUS_OBJERR	= 0x3,

	SEGV_MAPERR	= 0x1,
	SEGV_ACCERR	= 0x2,

	ITIMER_REAL	= 0x0,
	ITIMER_VIRTUAL	= 0x1,
	ITIMER_PROF	= 0x2,

	EPOLLIN		= 0x1,
	EPOLLOUT	= 0x4,
	EPOLLERR	= 0x8,
	EPOLLHUP	= 0x10,
	EPOLLRDHUP	= 0x2000,
	EPOLLET		= -0x80000000,
	EPOLL_CLOEXEC	= 0x80000,
	EPOLL_CTL_ADD	= 0x1,
	EPOLL_CTL_DEL	= 0x2,
	EPOLL_CTL_MOD	= 0x3,
};

typedef struct Itimerval Itimerval;
typedef struct EpollEvent EpollEvent;

#pragma pack on

struct Itimerval {
	Timeval	it_interval;
	Timeval	it_value;
};
struct EpollEvent {
	uint32	events;
	uint64	data;
};


#pragma pack off
// Created by cgo -cdefs - DO NOT EDIT
// cgo -cdefs defs_linux.go defs1_linux.go


enum {
	O_RDONLY	= 0x0,
	O_CLOEXEC	= 0x80000,
};

typedef struct Fpxreg Fpxreg;
typedef struct Xmmreg Xmmreg;
typedef struct Fpstate Fpstate;
typedef struct Fpxreg1 Fpxreg1;
typedef struct Xmmreg1 Xmmreg1;
typedef struct Fpstate1 Fpstate1;
typedef struct Fpreg1 Fpreg1;
typedef struct Mcontext Mcontext;

#pragma pack on

struct Fpxreg {
	uint16	significand[4];
	uint16	exponent;
	uint16	padding[3];
};
struct Xmmreg {
	uint32	element[4];
};
struct Fpstate {
	uint16	cwd;
	uint16	swd;
	uint16	ftw;
	uint16	fop;
	uint64	rip;
	uint64	rdp;
	uint32	mxcsr;
	uint32	mxcr_mask;
	Fpxreg	_st[8];
	Xmmreg	_xmm[16];
	uint32	padding[24];
};
struct Fpxreg1 {
	uint16	significand[4];
	uint16	exponent;
	uint16	padding[3];
};
struct Xmmreg1 {
	uint32	element[4];
};
struct Fpstate1 {
	uint16	cwd;
	uint16	swd;
	uint16	ftw;
	uint16	fop;
	uint64	rip;
	uint64	rdp;
	uint32	mxcsr;
	uint32	mxcr_mask;
	Fpxreg1	_st[8];
	Xmmreg1	_xmm[16];
	uint32	padding[24];
};
struct Fpreg1 {
	uint16	significand[4];
	uint16	exponent;
};
struct Mcontext {
	int64	gregs[23];
	Fpstate	*fpregs;
	uint64	__reserved1[8];
};


#pragma pack off
