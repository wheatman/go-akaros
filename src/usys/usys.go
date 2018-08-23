package usys

import "unsafe"

const (
	USYS_ABORT_SYSCALL_AT_ABS_UNIX = 0
	USYS_UNSET_ALARM               = 1
	USYS_GO_SYSCALL                = 2
	USYS_FUNC                      = 3
	USYS_FUTEX                     = 4
	USYS_SERIALIZE_ARGV_ENVP       = 5
	USYS_FREE                      = 6
	num_functions                  = 7
)

var functions [num_functions + 1]uintptr

func call(f uintptr, opt ...uintptr) int64
func call1(f, arg0 uintptr) int64

func init() {
	functions[num_functions] = 0xDEADBEEF
	call(0xDEAD000000000000 | uintptr(unsafe.Pointer(&functions[0])))
}

//go:nosplit
func Call(f int, opt ...uintptr) int64 {
	return call(functions[f], opt...)
}

//go:nosplit
func Call1(f int, arg0 uintptr) int64 {
	return call1(functions[f], arg0)
}
