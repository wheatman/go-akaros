// Copyright 2014 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package signal

import (
	"os"
	"syscall"
	"runtime/parlib"
)

const (
	numSig = parlib.NSIG
)

func signum(sig os.Signal) int {
	switch sig := sig.(type) {
	case syscall.Signal:
		i := int(sig)
		if i < 1 || i >= numSig {
			return -1
		}
		return i
	default:
		return -1
	}
}

func process_wrapper(sig int) {
	process(syscall.Signal(sig))
}

func enableSignal(sig int) {
	parlib.Signal(sig, process_wrapper)
}

func disableSignal(sig int) {
	parlib.Signal(sig, nil)
}
