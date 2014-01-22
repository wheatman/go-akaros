package main

import (
	"runtime/parlib"
	"syscall"
	"unsafe"
	"fmt"
)

func main() {
	cmd := "/bin/hello"
	argv := []string{cmd}
	envp := []string{"LD_LIBRARY_PATH=/lib"}

	// Declare all variables at top in case any
	// declarations require heap allocation (e.g., err1).
	var (
	    r1     uintptr
	    err1   error
	)

	// Set up arguments for proc_create
	bs_cmd, _ := parlib.ByteSliceFromString(cmd)
	bs_argv, _ := parlib.SlicePtrFromStrings(argv)
	bs_envp, _ := parlib.SlicePtrFromStrings(envp)
	pi, _ := parlib.ProcinfoPackArgs(bs_argv, bs_envp)
	__cmd := uintptr(unsafe.Pointer(&bs_cmd[0]))
	__cmdlen := uintptr(parlib.Cstrlen(bs_cmd))
	__pi := uintptr(unsafe.Pointer(&pi))

	// Call proc create.
	r1, _, err1 = syscall.RawSyscall(syscall.SYS_PROC_CREATE, __cmd, __cmdlen, __pi)
	if err1 != nil {
	    fmt.Printf("Error on SYS_PROC_CREATE: " + err1.Error())
	}
	child := int(r1)

	// Proc create succeeded, now run it! 
	r1, _, err1 = syscall.RawSyscall(syscall.SYS_PROC_RUN, r1, 0, 0)
	if err1 != nil {
	    fmt.Printf("Error on SYS_PROC_RUN: " + err1.Error())
	}

	// And wait for it...
	var status syscall.WaitStatus
	_, err1 = syscall.Waitpid(child, &status, 0)
	if err1 != nil {
	    fmt.Printf("Error on SYS_WAITPID" + err1.Error())
	}
}

