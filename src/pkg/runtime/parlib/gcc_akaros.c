// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

#include <futex.h>
#include <pthread.h>
#include <sys/syscall.h>
#include "gcc_akaros.h"

// Akaros syscalls
typedef void (*gcc_call_t)(void *arg);
const gcc_call_t gcc_syscall = (gcc_call_t)ros_syscall_sync;

// Akaros style futexes
static void __gcc_futex(void *__arg)
{
	// For now, akaros futexes don't support uaddr2 or val3, so we
	// just 0 them out.
	gcc_futex_arg_t *a = (gcc_futex_arg_t*)__arg;
	a->uaddr2 = NULL;
	a->val3 = 0;
	// Also, the minimum timout is 1us, so up it to that if it's too small
	if(a->timeout != NULL) {
		if(a->timeout->tv_sec == 0)
			if(a->timeout->tv_nsec < 1000L)
				a->timeout->tv_nsec = 1000L;
    }
	a->retval = futex(a->uaddr, a->op, a->val, a->timeout, a->uaddr2, a->val3);
}
const gcc_call_t gcc_futex = __gcc_futex;

// Akaros style pthread yields
static void __gcc_myield(void *__arg)
{
	// We should never pass an argument here
	assert(__arg == NULL);
	pthread_yield();
}
const gcc_call_t gcc_myield = __gcc_myield;

// Akaros style sigactions
static void __gcc_sigaction(void *__arg)
{
	gcc_sigaction_arg_t *a = (gcc_sigaction_arg_t*)__arg;
	a->ret = sigaction(a->sig, (struct sigaction*)a->act, (struct sigaction*)a->oact);
}
const gcc_call_t gcc_sigaction = __gcc_sigaction;

// Akaros sigprocmask
static void __gcc_sigprocmask(void *__arg)
{
	gcc_sigprocmask_arg_t *a = (gcc_sigprocmask_arg_t*)__arg;
	a->retval = pthread_sigmask(a->how, a->set, a->oset);
}
const gcc_call_t gcc_sigprocmask = __gcc_sigprocmask;

// enable_profalarm() and disable_profalarm()
static void __gcc_enable_profalarm(void *__arg)
{
	enable_profalarm(*((uint64_t*)__arg));
}
const gcc_call_t gcc_enable_profalarm = __gcc_enable_profalarm;

static void __gcc_disable_profalarm(void *__arg)
{
	disable_profalarm();
}
const gcc_call_t gcc_disable_profalarm = __gcc_disable_profalarm;
