package skymp_wrapper

import (
	"syscall"
)

var (
	scampLib, _         = syscall.LoadLibrary("scamp_lib")
	createServerProc, _ = syscall.GetProcAddress(scampLib, "CreateServer")
)

func Free() {
	syscall.FreeLibrary(scampLib)
}

func CreateServer() {
	syscall.SyscallN(uintptr(createServerProc), 3000, 10)
}
