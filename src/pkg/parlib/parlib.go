// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// +build akaros

package parlib

/*
#include <parlib.h>
#include <uthread.h>
#include <vcore.h>
#include <mcs.h>
*/
import "C"

// Dummy function called by the 'os' package to force a dependence on parlib.
// By forcing this dependance we are guaranteed to use the akaros
// cross-compiler linker to build the final executable.  In doing so, we use
// the standard libc _start() function rather than _rt0_GOARCH_GOOS() as the
// entry point to a go executable. Given our dependence on the c-parlib
// library, this is exactly what we want. In the future, we probably won't need
// this dummy init() call since we will likely access some of parlib's actualy
// functionality.  Additionally, we should make it a default dependance for the
// 'main' package when building for akaros.
func Init() {}

