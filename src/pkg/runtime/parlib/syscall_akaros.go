// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package parlib

/*
#include <benchutil/alarm.h>
#include <parlib/uthread.h>
#include <sys/syscall.h>

void abort_syscall_at(struct alarm_waiter *waiter, uint64_t deadline)
{
	init_awaiter(waiter, alarm_abort_sysc);
	waiter->data = current_uthread;
	set_awaiter_abs_unix(waiter, deadline);
	set_alarm(waiter);
}

void cancel_abort_syscall(struct alarm_waiter *waiter)
{
	unset_alarm(waiter);
}
*/
import "C"
import "unsafe"

type SyscallType C.struct_syscall
const (
	MAX_ERRSTR_LEN = C.MAX_ERRSTR_LEN
)

func Syscall(_num uint32, _a0, _a1, _a2, _a3, _a4, _a5 int) (ret int, err int32, errstr string) {
	var syscall SyscallType
	syscall.num = C.uint(_num)
	syscall.ev_q = (*C.struct_event_queue)(unsafe.Pointer(nil))
	syscall.arg0 = C.long(_a0)
	syscall.arg1 = C.long(_a1)
	syscall.arg2 = C.long(_a2)
	syscall.arg3 = C.long(_a3)
	syscall.arg4 = C.long(_a4)
	syscall.arg5 = C.long(_a5)
	C.ros_syscall_sync((*C.struct_syscall)(unsafe.Pointer(&syscall)))
	return int(syscall.retval), int32(syscall.err), C.GoString(&syscall.errstr[0])
}

func AbortSyscallAt(waiter *AlarmWaiter, deadline int64) {
	C.abort_syscall_at((*C.struct_alarm_waiter)(unsafe.Pointer(waiter)), C.uint64_t(deadline))
}

func CancelAbortSyscall(waiter *AlarmWaiter) {
	C.cancel_abort_syscall((*C.struct_alarm_waiter)(unsafe.Pointer(waiter)))
}
