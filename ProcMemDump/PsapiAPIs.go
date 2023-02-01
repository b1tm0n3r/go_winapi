package main

import (
	"syscall"
	"unsafe"
)

var (
	psapi_DLL                    = syscall.NewLazyDLL("psapi.dll")
	procEnumProcesses            = psapi_DLL.NewProc("EnumProcesses")
	procGetProcessImageFileNameA = psapi_DLL.NewProc("GetProcessImageFileNameA")
)

// Base: https://learn.microsoft.com/en-us/windows/win32/api/psapi/nf-psapi-enumprocesses
func EnumProcesses(out_lpidProcess_ptr []uint32, cb uint32, out_lpcbNeeded *uint32) uint {
	ret, _, _ := procEnumProcesses.Call(
		uintptr(unsafe.Pointer(&out_lpidProcess_ptr[0])),
		uintptr(cb),
		uintptr(unsafe.Pointer(out_lpcbNeeded)))

	return uint(ret) // Return should be non-zero
}

// Base: https://learn.microsoft.com/en-us/windows/win32/api/psapi/nf-psapi-enumprocessmodules
func GetProcessImageFileNameA(hProcess uintptr, out_lpImageFileName []byte, size uint32) uint {
	ret, _, _ := procGetProcessImageFileNameA.Call(
		uintptr(hProcess),
		uintptr(unsafe.Pointer(&out_lpImageFileName[0])),
		uintptr(size))

	return uint(ret) // Return specifies length of the string
}
