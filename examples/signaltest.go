// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"io/ioutil"
	"os"
	"os/signal"
	"os/exec"
	"runtime"
	"strconv"
	"syscall"
	"time"
	"fmt"
)

func waitSig(c <-chan os.Signal, sig os.Signal) {
	select {
	case s := <-c:
		if s != sig {
			fmt.Printf("signal was %v, want %v", s, sig)
		}
	case <-time.After(1 * time.Second):
		fmt.Printf("timeout waiting for %v", sig)
	}
}

func main() {
	TestSignal()
	TestStress(false)
	TestStop()
}

func TestSignal() {
	// Ask for SIGHUP
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP)
	defer signal.Stop(c)

	// Send this process a SIGHUP
	syscall.Kill(syscall.Getpid(), syscall.SIGHUP)
	waitSig(c, syscall.SIGHUP)

	// Ask for everything we can get.
	c1 := make(chan os.Signal, 1)
	signal.Notify(c1)

	// Send this process a SIGWINCH
	syscall.Kill(syscall.Getpid(), syscall.SIGWINCH)
	waitSig(c1, syscall.SIGWINCH)

	// Send two more SIGHUPs, to make sure that
	// they get delivered on c1 and that not reading
	// from c does not block everything.
	syscall.Kill(syscall.Getpid(), syscall.SIGHUP)
	waitSig(c1, syscall.SIGHUP)
	syscall.Kill(syscall.Getpid(), syscall.SIGHUP)
	waitSig(c1, syscall.SIGHUP)

	// The first SIGHUP should be waiting for us on c.
	waitSig(c, syscall.SIGHUP)
}

func TestStress(short bool) {
	dur := 3 * time.Second
	if short {
		dur = 100 * time.Millisecond
	}
	defer runtime.GOMAXPROCS(runtime.GOMAXPROCS(4))
	done := make(chan bool)
	finished := make(chan bool)
	go func() {
		sig := make(chan os.Signal, 1)
		signal.Notify(sig, syscall.SIGUSR1)
		defer signal.Stop(sig)
	Loop:
		for {
			select {
			case <-sig:
			case <-done:
				break Loop
			}
		}
		finished <- true
		//for {runtime.Gosched()}
	}()
	go func() {
	Loop:
		for {
			select {
			case <-done:
				break Loop
			default:
				syscall.Kill(syscall.Getpid(), syscall.SIGUSR1)
				runtime.Gosched()
			}
		}
		finished <- true
		//for {runtime.Gosched()}
	}()
	time.Sleep(dur)
	close(done)
	<-finished
	<-finished
	// When run with 'go test -cpu=1,2,4' SIGUSR1 from this test can slip
	// into subsequent TestSignal() causing failure.
	// Sleep for a while to reduce the possibility of the failure.
	time.Sleep(10 * time.Millisecond)
}

var sendUncaughtSighup = flag.Int("send_uncaught_sighup", 0, "send uncaught SIGHUP during TestStop")

// Test that Stop cancels the channel's registrations.
func TestStop() {
	sigs := []syscall.Signal{
		syscall.SIGWINCH,
		syscall.SIGHUP,
	}

	for _, sig := range sigs {
		// Send the signal.
		// If it's SIGWINCH, we should not see it.
		// If it's SIGHUP, maybe we'll die. Let the flag tell us what to do.
		if sig != syscall.SIGHUP || *sendUncaughtSighup == 1 {
			syscall.Kill(syscall.Getpid(), sig)
		}
		time.Sleep(10 * time.Millisecond)

		// Ask for signal
		c := make(chan os.Signal, 1)
		signal.Notify(c, sig)
		defer signal.Stop(c)

		// Send this process that signal
		syscall.Kill(syscall.Getpid(), sig)
		waitSig(c, sig)

		signal.Stop(c)
		select {
		case s := <-c:
			fmt.Printf("unexpected signal %v\n", s)
		case <-time.After(10 * time.Millisecond):
			// nothing to read - good
		}

		// Send the signal.
		// If it's SIGWINCH, we should not see it.
		// If it's SIGHUP, maybe we'll die. Let the flag tell us what to do.
		if sig != syscall.SIGHUP || *sendUncaughtSighup == 2 {
			syscall.Kill(syscall.Getpid(), sig)
		}

		select {
		case s := <-c:
			fmt.Printf("unexpected signal %v\n", s)
		case <-time.After(10 * time.Millisecond):
			// nothing to read - good
		}
	}
}

// Test that when run under nohup, an uncaught SIGHUP does not kill the program,
// but a
func TestNohup() {
	// Ugly: ask for SIGHUP so that child will not have no-hup set
	// even if test is running under nohup environment.
	// We have no intention of reading from c.
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP)

	// When run without nohup, the test should crash on an uncaught SIGHUP.
	// When run under nohup, the test should ignore uncaught SIGHUPs,
	// because the runtime is not supposed to be listening for them.
	// Either way, TestStop should still be able to catch them when it wants them
	// and then when it stops wanting them, the original behavior should resume.
	//
	// send_uncaught_sighup=1 sends the SIGHUP before starting to listen for SIGHUPs.
	// send_uncaught_sighup=2 sends the SIGHUP after no longer listening for SIGHUPs.
	//
	// Both should fail without nohup and succeed with nohup.

	for i := 1; i <= 2; i++ {
		out, err := exec.Command(os.Args[0], "-test.run=TestStop", "-send_uncaught_sighup="+strconv.Itoa(i)).CombinedOutput()
		if err == nil {
			fmt.Printf("ran test with -send_uncaught_sighup=%d and it succeeded: expected failure.\nOutput:\n%s\n", i, out)
		}
	}

	signal.Stop(c)

	// Again, this time with nohup, assuming we can find it.
	_, err := os.Stat("/usr/bin/nohup")
	if err != nil {
		fmt.Printf("cannot find nohup; skipping second half of test")
		return;
	}

	for i := 1; i <= 2; i++ {
		os.Remove("nohup.out")
		out, err := exec.Command("/usr/bin/nohup", os.Args[0], "-test.run=TestStop", "-send_uncaught_sighup="+strconv.Itoa(i)).CombinedOutput()

		data, _ := ioutil.ReadFile("nohup.out")
		os.Remove("nohup.out")
		if err != nil {
			fmt.Printf("ran test with -send_uncaught_sighup=%d under nohup and it failed: expected success.\nError: %v\nOutput:\n%s%s\n", i, err, out, data)
		}
	}
}
