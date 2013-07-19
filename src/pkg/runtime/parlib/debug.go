// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// +build akaros

package parlib

/*
#include <stdint.h>
#include <stdio.h>

void checkpoint(uint32_t arg)
{
	printf("checkpoint: %d\n", arg);
}
*/
import "C"

// Checkpoint function that I can sprinkle into the code to mark checkpoints
// during my porting effort
func Checkpoint( arg uint32 ) {
	C.checkpoint((C.uint32_t)(arg))
}

