// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// +build akaros

package parlib

/*
#include <stdlib.h>
#include <stdint.h>
#include <stdio.h>
#include <futex.h>
#include <limits.h>
#include <signal.h>
#include <pthread.h>

uint64_t __sigmap = 0;
int __sigpending = 0;
//void (*wtf)(int);
uint64_t wtf = 0;

// must match the one in gcc_akaros.h
typedef void (*__sigaction_t) (int, siginfo_t *, void *);


void sig_hand(int signr, void *info, void *ctxt) {
        if (in_vcore_context()) {
          __sigmap |= ((uint64_t)(1)) << (signr-1);
          __sigpending = 1;
          futex(&__sigpending, FUTEX_WAKE, INT_MAX, NULL, NULL, 0);
        } else {
          ((__sigaction_t)wtf)(signr, info, ctxt);

        }

}

void pthread_wake(int signr) {
        pthread_kill(pthread_self(), signr);
        pthread_yield();
}
*/
import "C"
import (
	"unsafe"
)

var (
	__SIG_ERR = -1
	__SIG_IGN = 1
	__SIG_DFL = 0
	SIG_ERR = *((*C.__sigaction_t)(unsafe.Pointer(&__SIG_ERR)))
	SIG_IGN = *((*C.__sigaction_t)(unsafe.Pointer(&__SIG_IGN)))
	SIG_DFL = *((*C.__sigaction_t)(unsafe.Pointer(&__SIG_DFL)))
)
const (
	NSIG = C._NSIG
)
type SignalHandler func(sig int)

var sighandlers [NSIG-1]SignalHandler
var sigact = SigactionT{Sigact: (C.__sigaction_t)(C.sig_hand), Flags: C.SA_SIGINFO};

// Implemented in runtime/sys_{GOOS}_{GOARCH}.s 
func defaultSighandler(sig int)
func defaultSighandler_for_c(sig int)

const ptrSize = 4 << (^uintptr(0) >> 63)
//go:nosplit
func add(p unsafe.Pointer, x uintptr) unsafe.Pointer {
        return unsafe.Pointer(uintptr(p) + x)
}
func get_value(f interface{}) (uint64) {
	return uint64(**(**uintptr)(unsafe.Pointer((add(unsafe.Pointer(&f), ptrSize)))))
}


func init() {
	for i := 0; i < (NSIG-1); i++ {
		Signal(i, defaultSighandler)
	}
	C.wtf = C.uint64_t(get_value(defaultSighandler_for_c))
	go process_signals()
}

func process_signals() {
	for {
		Futex((*int32)(&C.__sigpending), FUTEX_WAIT, 0, nil, nil, 0)
		C.__sigpending = 0
		sigmap := C.__sigmap
		C.__sigmap &^= sigmap

		signr := 0
		for sigmap != 0 {
			for {
				bit := sigmap & 1
				sigmap >>= 1
				signr++
				if bit == 1 {
					break
				}
			}
			// if somebody has updated the signal handler call the new one
			// else convert to internal signal and loop back around
			if (get_value(sighandlers[signr-1]) == get_value(defaultSighandler)) {
				C.pthread_wake(C.int(signr))
			} else {
				sighandlers[signr-1](signr)
			}
		}
	}
}

func Signal(signr int, newh SignalHandler) (SignalHandler, int) {
	if signr < 1 || signr >= NSIG {
		return nil, -1
	}

	oldh := sighandlers[signr-1]
	sighandlers[signr-1] = newh

	__sigact := sigact
	if newh == nil {
		__sigact.Sigact = SIG_DFL
	}
	ret := int(C.sigaction(C.int(signr), (*C.struct_sigaction)(unsafe.Pointer(&__sigact)), nil))
	if ret != 0 {
		return nil, ret
	}
	return oldh, ret
}

