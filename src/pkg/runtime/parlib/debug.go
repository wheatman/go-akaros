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
*/
import "C"

// Checkpoint function that I can sprinkle into the code to mark checkpoints
// during my porting effort
func Checkpoint( arg uint64, c int8 ) {
	C.checkpoint((C.ulonglong)(arg), (C.char)(c))
}

