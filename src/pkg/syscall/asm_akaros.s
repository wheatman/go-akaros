// Since we make all of our syscalls through the parlib interface fo for
// Akaros, I simply jump to the syscall wrappers provided in syscall_akaros.go
// which properly cast the arguments before sending them off to parlib.
//
// Additionally, Akaros sycalls NEVER block, so it's ok to implement regular
// syscalls the same as raw ones...

TEXT	·Syscall(SB),7,$0
	JMP	syscall·SyscallWrapper(SB)

TEXT	·Syscall6(SB),7,$0
	JMP	syscall·Syscall6Wrapper(SB)

TEXT	·RawSyscall(SB),7,$0
	JMP	syscall·SyscallWrapper(SB)

TEXT	·RawSyscall6(SB),7,$0
	JMP	syscall·Syscall6Wrapper(SB)

