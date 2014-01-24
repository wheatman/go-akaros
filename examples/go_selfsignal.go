package main

import (
	"os"
	"fmt"
	"time"
	"syscall"
	"runtime"
	"math/rand"
	"runtime/parlib"
)

var done chan bool = make(chan bool)

func handle_signal(sig int) {
	fmt.Printf("Got posix signal %s\n", syscall.Signal(sig))
	if (syscall.Signal(sig) == syscall.SIGTERM) { done <- true; for { runtime.Gosched() } }
}

func main() {
	fmt.Printf("Hello world from program %s!!\n", os.Args[0])
	rand.Seed( time.Now().UTC().UnixNano())
	for sig := 1; sig < parlib.SIGRTMIN; sig++ {
		parlib.Signal(sig, handle_signal)
		if (sig != int(syscall.SIGTERM)) && (sig != int(syscall.SIGKILL)) {
			d := time.Duration(rand.Intn(1000))
			go func(sig int, d time.Duration) {
				time.Sleep(d * time.Microsecond)
				syscall.Kill(int(parlib.Procinfo.Pid), syscall.Signal(sig))
			}(sig, d)
		}
	}
	time.Sleep(1000 * time.Millisecond)
	syscall.Kill(int(parlib.Procinfo.Pid), syscall.SIGTERM)
	<-done
	fmt.Printf("Exiting...\n")
}

