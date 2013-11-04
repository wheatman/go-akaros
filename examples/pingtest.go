package main

import (
	"fmt"
	"net"
)

func main() {
	raddr := &net.IPAddr{IP: net.IPv4(127, 0, 0, 1).To4()}
	laddr := &net.IPAddr{IP: net.IPv4(127, 0, 0, 1)}
	fmt.Println("raddr: ", raddr)
	fmt.Println("laddr: ", laddr)

	ipconn, err := net.DialIP("ip:icmp", laddr, raddr)
	if err != nil {
		fmt.Println("Error: ", err)
	} else {
		fmt.Println(ipconn)
	}
}
