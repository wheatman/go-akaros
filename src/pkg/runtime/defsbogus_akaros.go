// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ignore

/*
Input to cgo -cdefs

go-akaros-386 tool cgo -cdefs defsbogus_akaros.go > defs_akaros_386.h
go-akaros-amd64 tool cgo -cdefs defsbogus_akaros.go > defs_akaros_amd64.h
*/

package parlib

import "C"

const (
	ITIMER_REAL    = 0
	ITIMER_VIRTUAL = 0
	ITIMER_PROF    = 0

	EPOLLIN       = 0
	EPOLLOUT      = 0
	EPOLLERR      = 0
	EPOLLHUP      = 0
	EPOLLRDHUP    = 0
	EPOLLET       = 0
	EPOLL_CLOEXEC = 0
	EPOLL_CTL_ADD = 0
	EPOLL_CTL_DEL = 0
	EPOLL_CTL_MOD = 0
)

type EpollEvent *void
