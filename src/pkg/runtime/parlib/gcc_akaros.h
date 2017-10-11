// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

#include <signal.h>
#include <time.h>
#include <parlib/alarm.h>
#include <sys/syscall.h>

typedef struct syscall gcc_syscall_arg_t;

typedef struct gcc_futex_arg {
	int *uaddr;
	int op;
	int val;
	struct timespec *timeout;
	int *uaddr2;
	int val3;
	int retval;	
} gcc_futex_arg_t;

#undef sa_handler
struct parlib_sigaction {
    __sighandler_t sa_handler;
    sigset_t sa_mask;
    unsigned int sa_flags;
	unsigned int padding;
    void (*sa_restorer) (void);
};

typedef struct gcc_sigaction_arg {
	int sig;
	struct parlib_sigaction *act;
	struct parlib_sigaction *oact;
	int ret;
} gcc_sigaction_arg_t;

typedef struct gcc_sigprocmask_arg {
	int how;
	sigset_t *set;
	sigset_t *oset;
	int retval;
} gcc_sigprocmask_arg_t;

typedef TAILQ_ENTRY(parlib_alarm_waiter) parlib_alarm_waiter_tailq_entry_t;
struct parlib_alarm_waiter {
    uint64_t                          wake_up_time;   /* tsc time */
    void (*func) (struct parlib_alarm_waiter *waiter);
    void                              *data;
    parlib_alarm_waiter_tailq_entry_t next;
    bool                              on_tchain;
};

