// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package time

import (
	"errors"
	"runtime"
)

func init() {
	// force US/Pacific for time zone tests
	ForceUSPacificForTesting()
}

var Interrupt = interrupt
var DaysIn = daysIn

func empty(now int64, arg interface{}) {}

// Test that a runtimeTimer with a duration so large it overflows
// does not cause other timers to hang.
//
// This test has to be in internal_test.go since it fiddles with
// unexported data structures.
func CheckRuntimeTimerOverflow() error {
	// We manually create a runtimeTimer to bypass the overflow
	// detection logic in NewTimer: we're testing the underlying
	// runtime.addtimer function.
	r := &runtimeTimer{
		when: runtimeNano() + (1<<63 - 1),
		f:    empty,
		arg:  nil,
	}
	startTimer(r)

	timeout := 100 * Millisecond
	switch runtime.GOOS {
	// Allow more time for gobuilder to succeed.
	case "windows":
		timeout = Second
	case "plan9":
		// TODO(0intro): We don't know why it is needed.
		timeout = 3 * Second
	case "akaros":
		// TODO(wheatman): We also don't know why it is needed.
		timeout = Second
	}

	// Start a goroutine that should send on t.C before the timeout.
	t := NewTimer(1)

	defer func() {
		// Subsequent tests won't work correctly if we don't stop the
		// overflow timer and kick the timer proc back into service.
		//
		// The timer proc is now sleeping and can only be awoken by
		// adding a timer to the *beginning* of the heap. We can't
		// wake it up by calling NewTimer since other tests may have
		// left timers running that should have expired before ours.
		// Instead we zero the overflow timer duration and start it
		// once more.
		stopTimer(r)
		t.Stop()
		r.when = 0
		startTimer(r)
	}()

	// Try to receive from t.C before the timeout. It will succeed
	// iff the previous sleep was able to finish. We're forced to
	// spin and yield after trying to receive since we can't start
	// any more timers (they might hang due to the same bug we're
	// now testing).
	stop := Now().Add(timeout)
	for {
		select {
		case <-t.C:
			return nil // It worked!
		default:
			if Now().After(stop) {
				return errors.New("runtime timer stuck: overflow in addtimer")
			}
			// Issue 6874. This test previously called runtime.Gosched to try to yield
			// to the goroutine servicing t, however the scheduler has a bias towards the
			// previously running goroutine in an idle system. Combined with high load due
			// to all CPUs busy running tests t's goroutine could be delayed beyond the
			// timeout window.
			//
			// Calling runtime.GC() reduces the worst case lantency for scheduling t by 20x
			// under the current Go 1.3 scheduler.
			runtime.GC()
		}
	}
}
