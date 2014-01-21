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
	pi, _ := parlib.ProcinfoPackArgs(argv, envp)
	bs_cmd, _ := parlib.ByteSliceFromString(cmd)
	__cmd := uintptr(unsafe.Pointer(&bs_cmd[0]))
    __cmdlen := uintptr(parlib.Cstrlen(bs_cmd))
	__pi := uintptr(unsafe.Pointer(&pi))

    // Call proc create.
    r1, _, err1 = syscall.RawSyscall(syscall.SYS_PROC_CREATE, __cmd, __cmdlen, __pi)
    if err1 != nil {
        fmt.Printf("Error on SYS_PROC_CREATE: " + err1.Error())
    }

    // Proc create succeeded, now run it! 
    r1, _, err1 = syscall.RawSyscall(syscall.SYS_PROC_RUN, r1, 0, 0)
    if err1 != nil {
        fmt.Printf("Error on SYS_PROC_RUN: " + err1.Error())
    }
}

