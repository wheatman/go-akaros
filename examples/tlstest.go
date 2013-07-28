package main

/*
#include <stdio.h>
#include <stdint.h>

__thread void *tls_var1;
__thread void *tls_var2;

void print_tls_info() {
	void *tls_base;
#ifdef __x86_64__
	asm volatile("mov %%fs:16, %0" : "=r" (tls_base) ::);
#else
	asm volatile("mov %%gs:8, %0" : "=r" (tls_base) ::);
#endif
	printf("tls_base manual: %p\n", tls_base);
	printf("g_offset raw: %p\n", (void*)(((uintptr_t)tls_base - (uintptr_t)&tls_var1)));
	printf("m_offset raw: %p\n", (void*)(((uintptr_t)tls_base - (uintptr_t)&tls_var2)));
	printf("g_offset: %p\n", (void*)(0 - ((uintptr_t)tls_base - (uintptr_t)&tls_var1)));
	printf("m_offset: %p\n", (void*)(0 - ((uintptr_t)tls_base - (uintptr_t)&tls_var2)));
}
*/
import "C"

func main() {
    C.print_tls_info()
}

