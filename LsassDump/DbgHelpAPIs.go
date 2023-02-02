package main

import (
	"syscall"
)

var (
	dbgHelp_DLL           = syscall.NewLazyDLL("Dbghelp.dll")
	procMiniDumpWriteDump = dbgHelp_DLL.NewProc("MiniDumpWriteDump")
)

// Base: https://learn.microsoft.com/en-us/windows/win32/api/minidumpapiset/nf-minidumpapiset-minidumpwritedump
func MiniDumpWriteDump(hProcess uintptr, processId uint32, hFile uintptr, dumpType uint32,
	exceptionParam uintptr, userStreamParam uintptr, callbackParam uintptr) bool {
	ret, _, _ := procMiniDumpWriteDump.Call(
		uintptr(hProcess),
		uintptr(processId),
		uintptr(hFile),
		uintptr(dumpType),
		exceptionParam,
		userStreamParam,
		callbackParam)

	return uint(ret) == 0 // Return should be non-zero
}
