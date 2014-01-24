// Created by cgo -cdefs - DO NOT EDIT
// cgo -cdefs defs2_linux.go

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

	O_RDONLY	= 0x0,
	O_CLOEXEC	= 0x80000,

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

typedef struct Fpreg Fpreg;
typedef struct Fpxreg Fpxreg;
typedef struct Xmmreg Xmmreg;
typedef struct Fpstate Fpstate;
typedef struct Itimerval Itimerval;
typedef struct EpollEvent EpollEvent;

#pragma pack on

struct Fpreg {
	uint16	significand[4];
	uint16	exponent;
};
struct Fpxreg {
	uint16	significand[4];
	uint16	exponent;
	uint16	padding[3];
};
struct Xmmreg {
	uint32	element[4];
};
struct Fpstate {
	uint32	cw;
	uint32	sw;
	uint32	tag;
	uint32	ipoff;
	uint32	cssel;
	uint32	dataoff;
	uint32	datasel;
	Fpreg	_st[8];
	uint16	status;
	uint16	magic;
	uint32	_fxsr_env[6];
	uint32	mxcsr;
	uint32	reserved;
	Fpxreg	_fxsr_st[8];
	Xmmreg	_xmm[8];
	uint32	padding1[44];
	byte	anon0[48];
};
struct Itimerval {
	Timeval	it_interval;
	Timeval	it_value;
};
struct EpollEvent {
	uint32	events;
	uint64	data;
};


#pragma pack off
