// Copyright 2013 The Go Authors.	All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"os"
	"net"
	"fmt"
)

const proto = "tcp"
const addr = "127.0.0.1:0"

func main() {

	fmt.Println("Starting....\n");
	// Open a connection for both the server and the client described below
	ln, err := net.Listen(proto, addr)
	if err != nil {
		fmt.Errorf("listen failed: %v", err)
		return
	}
	defer ln.Close()

	// Spawn two goroutines. This first reads one byte from stdin and writes to
	// a tcp conversation. The second reads from the byte from the tcp
	// conversation and writes two copies of it to stdout. 
	done := make(chan int)
	go server(ln, done)
	go client(ln, done)

	// Wait until one of them dies for some reason (maybe never)
	<-done
	fmt.Println("Exiting....\n");
}

func server(ln net.Listener, done chan<- int) {

	// Signal that we are done only once this goroutine exits
	defer func() { done <- 1 }()

	// Listen for a connection 
	c, err := ln.Accept()
	if err != nil {
		fmt.Errorf("Listener.Accept failed: %v", err)
		return
	}
	defer c.Close()

	// Do this forever....
	b := make([]byte, 1)
	for {
		// Read a single byte from the connection
		_, err := c.Read(b)
		if err != nil {
			fmt.Errorf("Conn.Read failed: %v", err)
			return
		}

		// And push two copies of it to stdout
		_, err = os.Stdout.Write(b)
		if err != nil {
			fmt.Errorf("os.Stdout.Write failed: %v", err)
			return
		}
		_, err = os.Stdout.Write(b)
		if err != nil {
			fmt.Errorf("os.Stdout.Write failed: %v", err)
			return
		}
	}
}

func client(ln net.Listener, done chan<- int) {

	// Signal that we are done only once this goroutine exits
	defer func() { done <- 1 }()

	// Dial the server
	c, err := net.Dial(proto, ln.Addr().String())
	if err != nil {
		fmt.Errorf("Dial failed: %v", err)
		return
	}
	defer c.Close()

	// Do this forever....
	var b []byte = make([]byte, 1)
	for {
		// Read a single byte from stdin
		_, err := os.Stdin.Read(b)
		if err != nil {
			fmt.Errorf("os.Stdin.Read failed: %v", err)
			return
		}

		// And push it to to the connection
		_, err = c.Write(b)
		if err != nil {
			fmt.Errorf("Conn.Write failed: %v", err)
			return
		}
	}
}

