package main

/*
#include <stdlib.h>
#include <stdint.h>
#include <stdio.h>
#include <futex.h>
#include <limits.h>

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
	"os"
	"fmt"
	"time"
	"syscall"
	"runtime"
	"math/rand"
	"runtime/parlib"
)

const NSIG = 32
var done chan bool = make(chan bool)

func handle_signal(sig syscall.Signal) {
	fmt.Printf("Got posix signal %s\n", sig)
	if (sig == syscall.SIGTERM) { done <- true; for { runtime.Gosched() } }
}

func process_signals() {
	for {
		parlib.Futex((*int32)(&C.__sigpending), parlib.FUTEX_WAIT, 0, nil, nil, 0)
		C.__sigpending = 0
		sigmap := C.__sigmap
		signr := uint(0)
		for sigmap != 0 {
			for {
				bit := sigmap & 1
				sigmap >>= 1
				signr++
				if bit == 1 {
					break
				}
			}
			handle_signal(syscall.Signal(signr))
			C.__sigmap &^= ((_Ctype_uint64_t)(1) << (signr-1))
		}
	}
}

func main() {
	fmt.Printf("Hello world from program %s!!\n", os.Args[0])
	rand.Seed( time.Now().UTC().UnixNano())
	go process_signals()
	for sig := 1; sig < NSIG; sig++ {
		sigact := parlib.SigactionType{(*[0]byte)(C.sig_hand), 0, nil, 0};
		parlib.Sigaction(sig, &sigact, nil)
		if (sig != int(syscall.SIGTERM)) && (sig != int(syscall.SIGKILL)) {
			d := time.Duration(rand.Intn(1000))
			go func(sig int, d time.Duration) {
				time.Sleep(d * time.Microsecond)
				syscall.Kill(parlib.Procinfo.Pid(), syscall.Signal(sig))
			}(sig, d)
		}
	}
	time.Sleep(1000 * time.Millisecond)
	syscall.Kill(parlib.Procinfo.Pid(), syscall.SIGTERM)
	<-done
	fmt.Printf("Exiting...\n")
}

