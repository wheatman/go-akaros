// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// +build akaros

package parlib

/*
#include <stdint.h>
#include <stdio.h>

void checkpoint(long long unsigned arg, char c)
{
	if (c == 'x')
		printf("checkpoint: 0x%llx\n", arg);
	else
		printf("checkpoint: %llu\n", arg);
}
void syscallError(long num)
{
	printf("Error: Syscall %ld not yet supported!\n", num);
}
*/
import "C"

// Checkpoint function that I can sprinkle into the code to mark checkpoints
// during my porting effort
func Checkpoint(arg uint64, c int8) {
	C.checkpoint((C.ulonglong)(arg), (C.char)(c))
}

// Print out an error related to the syscall specified
func SyscallError(num uintptr) {
	C.syscallError(C.long(num))
}

