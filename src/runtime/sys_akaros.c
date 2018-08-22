// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//
// System calls and other sys stuff for Akaros
//

#include "runtime.h"
#include "cgocall.h"
#include "defs_GOOS_GOARCH.h"
#include "os_GOOS.h"
#include "zsyscall_akaros.h"
#include "../cmd/ld/textflag.h"

// We extern these libc functions here, so we don't have to reimplement the logic
// they entails just for compilation in the kenc world.
// We use asmcgocall() to call them.
#pragma cgo_import_static gcc_syscall
#pragma cgo_import_static gcc_futex
#pragma cgo_import_static gcc_myield
#pragma cgo_import_static gcc_sigprocmask
#pragma cgo_import_static gcc_enable_profalarm
#pragma cgo_import_static gcc_disable_profalarm
typedef void (*gcc_call_t)(void *arg);
extern gcc_call_t gcc_syscall;
extern gcc_call_t gcc_futex;
extern gcc_call_t gcc_myield;
extern gcc_call_t gcc_sigprocmask;
extern gcc_call_t gcc_enable_profalarm;
extern gcc_call_t gcc_disable_profalarm;

#pragma textflag NOSPLIT
static intgo strlen(int8 *string)
{
	int8 *temp = string;
	while(*temp != '\0')
		temp++;
	return temp-string;
}

#pragma textflag NOSPLIT
static inline bool mult_will_overflow_int64(int64 a, int64 b)
{
    if (!a)
        return false;
    return ((uint64)(-1) >> 1) / a < b;
}

#pragma textflag NOSPLIT
static inline int64 tsc2sec(int64 tsc_time)
{
    return tsc_time / __procinfo.tsc_freq;
}

#pragma textflag NOSPLIT
static inline int64 tsc2msec(int64 tsc_time)
{
    if (mult_will_overflow_int64(tsc_time, 1000LL))
        return tsc2sec(tsc_time) * 1000LL;
    else
        return (tsc_time * 1000LL) / __procinfo.tsc_freq;
}

#pragma textflag NOSPLIT
static inline int64 tsc2usec(int64 tsc_time)
{
    if (mult_will_overflow_int64(tsc_time, 1000000LL))
        return tsc2msec(tsc_time) * 1000LL;
    else
        return (tsc_time * 1000000LL) / __procinfo.tsc_freq;
}

#pragma textflag NOSPLIT
static inline int64 tsc2nsec(int64 tsc_time)
{
    if (mult_will_overflow_int64(tsc_time, 1000000000LL))
        return tsc2usec(tsc_time) * 1000LL;
    else
        return (tsc_time * 1000000000LL) / __procinfo.tsc_freq;
}

#pragma textflag NOSPLIT
static inline int64 sec2tsc(int64 sec)
{
    if (mult_will_overflow_int64(sec, __procinfo.tsc_freq))
        return (int64)(-1);
    else
        return sec * __procinfo.tsc_freq;
}

#pragma textflag NOSPLIT
static inline int64 msec2tsc(int64 msec)
{
    if (mult_will_overflow_int64(msec, __procinfo.tsc_freq))
        return sec2tsc(msec / 1000LL);
    else
        return (msec * __procinfo.tsc_freq) / 1000LL;
}

#pragma textflag NOSPLIT
static inline int64 usec2tsc(int64 usec)
{
    if (mult_will_overflow_int64(usec, __procinfo.tsc_freq))
        return msec2tsc(usec / 1000LL);
    else
        return (usec * __procinfo.tsc_freq) / 1000000LL;
}

#pragma textflag NOSPLIT
static inline int64 nsec2tsc(int64 nsec)
{
    if (mult_will_overflow_int64(nsec, __procinfo.tsc_freq))
        return usec2tsc(nsec / 1000LL);
    else
        return (nsec * __procinfo.tsc_freq) / 1000000000LL;
}

// Wrapper for making an akaros syscall through gcc
#define __akaros_syscall(sysc, n, a0, a1, a2, a3, a4, a5, perrno) \
do { \
	sysc->num = n;       \
	sysc->err = 0;       \
	sysc->retval = 0;    \
	sysc->flags = 0;     \
	sysc->ev_q = 0;      \
	sysc->u_data = 0;    \
	sysc->arg0 = a0;     \
	sysc->arg1 = a1;     \
	sysc->arg2 = a2;     \
	sysc->arg3 = a3;     \
	sysc->arg4 = a4;     \
	sysc->arg5 = a5;     \
	sysc->errstr[0] = 0; \
	runtime·asmcgocall(gcc_syscall, sysc); \
	if(perrno != nil) \
    	*perrno = sysc->err; \
} while(0);

#define akaros_syscall(sysc, n, a0, a1, a2, a3, a4, a5, perrno) \
	__akaros_syscall(((sysc)), ((uint32)(n)),                   \
	                 ((intgo)(a0)), ((intgo)(a1)),              \
	                 ((intgo)(a2)), ((intgo)(a3)),              \
	                 ((intgo)(a4)), ((intgo)(a5)),              \
	                 ((int32*)(perrno)))

#pragma textflag NOSPLIT
int32 runtime·getpid(void)
{
	return 	__procinfo.pid;
}

#pragma textflag NOSPLIT
int32 runtime·read(int32 fd, void* buf, int32 count)
{
	int32 errno;
	SyscallArg *sysc = (SyscallArg *)(g->sysc);
	akaros_syscall(sysc, SYS_read, fd, buf, count, 0, 0, 0, &errno);
	return sysc->retval;
}

#pragma textflag NOSPLIT
int32 runtime·write(uintptr fd, void* buf, int32 count)
{
	int32 errno;
	SyscallArg *sysc = (SyscallArg *)(g->sysc);
	akaros_syscall(sysc, SYS_write, fd, buf, count, 0, 0, 0, &errno);
	return sysc->retval;
}

#pragma textflag NOSPLIT
int32 runtime·open(int8* pathname, int32 flags, int32 mode)
{
	int32 errno;
	intgo len = strlen(pathname);
	SyscallArg *sysc = (SyscallArg *)(g->sysc);
	akaros_syscall(sysc, SYS_openat, AT_FDCWD, pathname, len, flags, mode, 0, &errno);
	return sysc->retval;
}

#pragma textflag NOSPLIT
int32 runtime·close(int32 fd)
{
	int32 errno;
	SyscallArg *sysc = (SyscallArg *)(g->sysc);
	akaros_syscall(sysc, SYS_close, fd, 0, 0, 0, 0, 0, &errno);
	return sysc->retval;
}

#pragma textflag NOSPLIT
uint8* runtime·mmap(byte* addr, uintptr len, int32 prot,
                    int32 flags, int32 fd, uint32 offset)
{
	int32 errno;
	SyscallArg *sysc = (SyscallArg *)(g->sysc);
	akaros_syscall(sysc, SYS_mmap, addr, len, prot, flags, fd, offset, &errno);
	return (uint8*)sysc->retval;
}

#pragma textflag NOSPLIT
void runtime·munmap(byte* addr, uintptr len)
{
	SyscallArg *sysc = (SyscallArg *)(g->sysc);
	akaros_syscall(sysc, SYS_munmap, addr, len, 0, 0, 0, 0, nil);
}

#pragma textflag NOSPLIT
void runtime·osyield(void)
{
	runtime·asmcgocall(gcc_myield, nil);
}

#pragma textflag NOSPLIT
void runtime·usleep(uint32 usec)
{
	SyscallArg *sysc = (SyscallArg *)(g->sysc);
	akaros_syscall(sysc, SYS_block, usec, 0, 0, 0, 0, 0, nil);
}

#pragma textflag NOSPLIT
int64 runtime·nanotime(void)
{
	return tsc2nsec(runtime·cputicks());
}

#pragma textflag NOSPLIT
void time·now(int64 sec, int32 nsec)
{
	// NOTE: we should add real time base to the nsec since start,
	// but gettimeofday syscall is missing on akaros.
	int64 time = tsc2nsec(runtime·cputicks() - __proc_global_info.tsc_cycles_last) + __proc_global_info.walltime_ns_last;
	sec = time / 1000000000LL;
	nsec = time - sec * 1000000000LL;
	FLUSH(&sec);
	FLUSH(&nsec);
}

#pragma textflag NOSPLIT
void runtime·exit(int32 status)
{
	intgo pid = runtime·getpid();
	SyscallArg *sysc = (SyscallArg *)(g->sysc);
	akaros_syscall(sysc, SYS_proc_destroy, pid, status, 0, 0, 0, 0, nil);
	runtime·throw("Exit Returned: We should never get here!");
}

#pragma textflag NOSPLIT
int32 runtime·futex(uint32 *uaddr, int32 op, uint32 val,
                    Timespec *timeout, uint32 *uaddr2, uint32 val3)
{
	FutexArg a;
	a.uaddr = (int32*)uaddr;
	a.op = op;
	a.val = val;
	a.timeout = timeout;
	a.uaddr2 = (int32*)uaddr2;
	a.val3 = val3;
	runtime·asmcgocall(gcc_futex, &a);
	return a.retval;
}

#pragma textflag NOSPLIT
void runtime·raise(int32 sig)
{
	USED(sig);
}

#pragma textflag NOSPLIT
int32 runtime·sigprocmask(int32 how, Sigset *set, Sigset *oldset)
{
	SigprocmaskArg a;
	a.how = how;
	a.set = set;
	a.oset = oldset;
	runtime·asmcgocall(gcc_sigprocmask, &a);
	return a.retval;
}

#pragma textflag NOSPLIT
void runtime·enable_profalarm(uint64 usecs)
{
	runtime·asmcgocall(gcc_enable_profalarm, &usecs);
}

#pragma textflag NOSPLIT
void runtime·disable_profalarm()
{
	runtime·asmcgocall(gcc_disable_profalarm, nil);
}

int32 runtime·epollcreate(int32 size)
{ USED(size); return -1; }

int32 runtime·epollcreate1(int32 flags)
{ USED(flags); return -1; }

int32 runtime·epollctl(int32 epfd, int32 op, int32 fd, EpollEvent *ev)
{ USED(epfd); USED(op); USED(fd); USED(ev); return -1; }

int32 runtime·epollwait(int32 epfd, EpollEvent *ev, int32 nev, int32 timeout)
{ USED(epfd); USED(ev); USED(nev); USED(timeout); return -1; }

void runtime·closeonexec(int32 fd)
{ USED(fd); }
