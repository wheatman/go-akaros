// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ignore

/*
Input to cgo -cdefs

GOARCH=amd64 go tool cgo -cdefs defs_akaros.go defs1_akaros.go >defs_akaros_amd64.h
*/

package runtime

/*
// Akaros glibc and Akaros kernel define different and conflicting
// definitions for struct sigaction, struct timespec, etc.
// We want the kernel ones, which are in the asm/* headers.
// But then we'd get conflicts when we include the system
// headers for things like ucontext_t, so that happens in
// a separate file, defs1.go.

#include <asm/posix_types.h>
#define size_t __kernel_size_t
#include <asm/signal.h>
#include <asm/siginfo.h>
#include <asm/mman.h>
#include <asm-generic/errno.h>
#include <asm-generic/poll.h>
#include <linux/eventpoll.h>
*/
import "C"

const (
	EINTR  = C.EINTR
	EAGAIN = C.EAGAIN
	ENOMEM = C.ENOMEM

	PROT_NONE  = C.PROT_NONE
	PROT_READ  = C.PROT_READ
	PROT_WRITE = C.PROT_WRITE
	PROT_EXEC  = C.PROT_EXEC

	MAP_ANON    = C.MAP_ANONYMOUS
	MAP_PRIVATE = C.MAP_PRIVATE
	MAP_FIXED   = C.MAP_FIXED

	MADV_DONTNEED = C.MADV_DONTNEED

	FPE_INTDIV = C.FPE_INTDIV
	FPE_INTOVF = C.FPE_INTOVF
	FPE_FLTDIV = C.FPE_FLTDIV
	FPE_FLTOVF = C.FPE_FLTOVF
	FPE_FLTUND = C.FPE_FLTUND
	FPE_FLTRES = C.FPE_FLTRES
	FPE_FLTINV = C.FPE_FLTINV
	FPE_FLTSUB = C.FPE_FLTSUB

	BUS_ADRALN = C.BUS_ADRALN
	BUS_ADRERR = C.BUS_ADRERR
	BUS_OBJERR = C.BUS_OBJERR

	SEGV_MAPERR = C.SEGV_MAPERR
	SEGV_ACCERR = C.SEGV_ACCERR

	ITIMER_REAL    = C.ITIMER_REAL
	ITIMER_VIRTUAL = C.ITIMER_VIRTUAL
	ITIMER_PROF    = C.ITIMER_PROF

	EPOLLIN       = C.POLLIN
	EPOLLOUT      = C.POLLOUT
	EPOLLERR      = C.POLLERR
	EPOLLHUP      = C.POLLHUP
	EPOLLRDHUP    = C.POLLRDHUP
	EPOLLET       = C.EPOLLET
	EPOLL_CLOEXEC = C.EPOLL_CLOEXEC
	EPOLL_CTL_ADD = C.EPOLL_CTL_ADD
	EPOLL_CTL_DEL = C.EPOLL_CTL_DEL
	EPOLL_CTL_MOD = C.EPOLL_CTL_MOD
)

type Timespec C.struct_timespec
type Timeval C.struct_timeval
type Itimerval C.struct_itimerval
type EpollEvent C.struct_epoll_event
