// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// +build akaros

package parlib

/*
#include <stdint.h>
#include <stdio.h>

void printChar(char c)
{
	fprintf(stderr, "%c\n", c);
}
void printInt(int d)
{
	fprintf(stderr, "%d", d);
}
void printString(char *s)
{
	fprintf(stderr, "%s", s);
}
void printChars(char *s, int len)
{
	int i = 0;
	for (i=0; i<len; i++)
		fprintf(stderr, "%c", *(s+i));
}
void checkpoint(long long unsigned arg, char c)
{
	if (c == 'x')
		fprintf(stderr, "checkpoint: 0x%llx\n", arg);
	else
		fprintf(stderr, "checkpoint: %llu\n", arg);
}
void syscallError(long num)
{
	fprintf(stderr, "Error: Syscall %ld not yet supported!\n", num);
}
*/
import "C"
import "unsafe"

// Checkpoint function that I can sprinkle into the code to mark checkpoints
// during my porting effort
func Checkpoint(arg uint64, c int8) {
	C.checkpoint((C.ulonglong)(arg), (C.char)(c))
}

// Print an arbitraty character
func PrintInt(d int) {
	C.printInt(C.int(d))
}

// Print an arbitraty character
func PrintChar(c byte) {
	C.printChar(C.char(c))
}

// Print an arbitraty string of bytes as characters
func PrintChars(s []byte) {
	C.printChars((*C.char)(unsafe.Pointer(&s[0])), C.int(len(s)))
}

// Print a string 
func PrintString(s string) {
	C.printString(C.CString(s))
}

// Print out an error related to the syscall specified
func SyscallError(num uintptr) {
	C.syscallError(C.long(num))
}

