package main
package main

import (
	"syscall"
	"unsafe"
)

var (
	kernel32_DLL           = syscall.NewLazyDLL("kernel32.dll")
	procOpenProcess        = kernel32_DLL.NewProc("OpenProcess")
	procGetCurrentProcess  = kernel32_DLL.NewProc("GetCurrentProcess")
	procVirtualAllocEx     = kernel32_DLL.NewProc("VirtualAllocEx")
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

// Base: https://learn.microsoft.com/en-us/windows/win32/api/memoryapi/nf-memoryapi-virtualallocex
func VirtualAllocEx(hProcess uintptr, lpAddress uintptr, dwSize uint32,
	flAllocationType uint32, flProtect uint32) uintptr {

	ret, _, _ := procVirtualAllocEx.Call(
		hProcess,
		lpAddress,
		uintptr(dwSize),
		uintptr(flAllocationType),
		uintptr(flProtect))

	return ret // return value should be non-null
}

// Base: https://learn.microsoft.com/en-us/windows/win32/api/processthreadsapi/nf-processthreadsapi-getcurrentprocess
func GetCurrentProcess() uintptr {
	ret, _, _ := procGetCurrentProcess.Call()

	return ret // Return should be non-null (returns process handle)
}