// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// +build akaros

package parlib

/*
#include <stdlib.h>
#include <stdint.h>
#include <futex.h>
#include <limits.h>
#include <signal.h>

uint64_t __sigmap = 0;
int __sigpending = 0;
void sig_hand(int signr) {
    __sigmap |= ((uint64_t)(1)) << (signr-1);
    __sigpending = 1;
    futex(&__sigpending, FUTEX_WAKE, INT_MAX, NULL, NULL, 0);
}
*/
import "C"
import (
	"unsafe"
)

type SignalHandler func(sig int)

var sigact = Sigaction{(C.__sighandler_t)(C.sig_hand), 0, nil, 0};
var sighandlers = map[int]SignalHandler{}

func init() {
	go process_signals()
}

func process_signals() {
    for {
        Futex((*int32)(&C.__sigpending), FUTEX_WAIT, 0, nil, nil, 0)
        C.__sigpending = 0
        sigmap := C.__sigmap
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
            sighandlers[signr](signr)
            C.__sigmap &^= ((_Ctype_uint64_t)(1) << (uint(signr)-1))
        }
    }
}

func Signal(sig int, newh SignalHandler) (SignalHandler, int) {
	ret := int(C.sigaction(C.int(sig), (*C.struct_sigaction)(unsafe.Pointer(&sigact)), nil))
	if ret != 0 {
		return nil, ret
	}
	oldh := sighandlers[sig]
	sighandlers[sig] = newh
	return oldh, ret
}

