package main

import (
	"syscall"
	"unsafe"
)

var (
	kernel32_DLL    = syscall.NewLazyDLL("kernel32.dll")
	procOpenProcess = kernel32_DLL.NewProc("OpenProcess")
	procCloseHandle = kernel32_DLL.NewProc("CloseHandle")
)

// Base: https://learn.microsoft.com/en-us/windows/win32/api/processthreadsapi/nf-processthreadsapi-openprocess
func OpenProcess(dwDesiredAccess int32, bInheritHandle bool, dwProcessId uint32) uintptr {
	inherit_boolToUint := 0
	if bInheritHandle == true {
		inherit_boolToUint = 1
	}

	ret, _, _ := procOpenProcess.Call(
		uintptr(dwDesiredAccess),
		uintptr(inherit_boolToUint),
		uintptr(dwProcessId))

	return ret // Return should be non-null (returns process handle)
}

// Base: https://learn.microsoft.com/en-us/windows/win32/api/handleapi/nf-handleapi-closehandle
func CloseHandle(hObject uintptr) bool {
	ret, _, _ := procCloseHandle.Call(uintptr(unsafe.Pointer(&hObject)))

	return uint(ret) == 0 // Return should be non-null (returns process handle)
}
