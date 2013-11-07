// Copyright 2013 The Go Authors.	All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"os"
	"net"
	"fmt"
	"bufio"
)

const proto = "tcp"
var host string = "128.32.37.180"
var page string = "files/test.html"

func main() {

	// Parse the args
    args := os.Args
	if len(args) != 3 {
		fmt.Printf("Usage: %s HOST PAGE\n", args[0])
    } else {
        host = args[1]
        page = args[2]
    }
    fmt.Printf("Trying to access http://%s/%s\n", host, page)

    // Build the address for use in Dial
	addr := net.JoinHostPort(host, "80")

	// Dial the addr
	conn, err := net.Dial(proto, addr)
	if err != nil {
		fmt.Errorf("Dial failed: %v", err)
		return
	}
	defer conn.Close()

	// Build the GET string
	req := fmt.Sprintf("GET /%s\r\n\r\n", page)

	// Write it to the connection
	_, err = fmt.Fprintf(conn, req)
	if err != nil {
		fmt.Errorf("Fprintf failed: %v", err)
		return
	}

	// Finally, echo all of the data returned on the connection to stdout
	reader := bufio.NewReader(conn)
	for {
		line, err := reader.ReadString('\n')
		fmt.Printf(line)
		if err != nil {
			break;
		}
	}
}

