package main

import (
	"syscall"
	"unsafe"
)

var (
	dbgHelp_DLL           = syscall.NewLazyDLL("Dbghelp.dll")
	procMiniDumpWriteDump = dbgHelp_DLL.NewProc("MiniDumpWriteDump")
)

// Base: https://learn.microsoft.com/en-us/windows/win32/api/minidumpapiset/nf-minidumpapiset-minidumpwritedump
func ProcMiniDumpWriteDump(hProcess uintptr, processId uint32, hFile uintptr, dumpType int,
	exceptionParam uintptr, userStreamParam uintptr, callbackParam uintptr) bool {
	ret, _, _ := procMiniDumpWriteDump.Call(
		uintptr(unsafe.Pointer(&hProcess)),
		uintptr(processId),
		uintptr(unsafe.Pointer(&hFile)),
		uintptr(dumpType),
		uintptr(unsafe.Pointer(&exceptionParam)),
		uintptr(unsafe.Pointer(&userStreamParam)),
		uintptr(unsafe.Pointer(&callbackParam)))

	return uint(ret) == 0 // Return should be non-zero
}
