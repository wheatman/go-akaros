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
	printf("%c\n", c);
}
void printString(char *s)
{
	printf("%s\n", s);
}
void printChars(char *s, int len)
{
	int i = 0;
	for (i=0; i<len; i++)
		printf("%c", *(s+i));
	printf("\n");
}
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
import "unsafe"

// Checkpoint function that I can sprinkle into the code to mark checkpoints
// during my porting effort
func Checkpoint(arg uint64, c int8) {
	C.checkpoint((C.ulonglong)(arg), (C.char)(c))
}

// Print an arbitraty character
func PrintChar(c byte) {
	C.printChar(C.char(c))
}

// Print an arbitraty string of bytes as characters
func PrintChars(s [256]int8) {
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

