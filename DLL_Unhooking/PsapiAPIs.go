package main

import (
	"syscall"
	"unsafe"
)

var (
	psapi_DLL                = syscall.NewLazyDLL("psapi.dll")
	procGetModuleInformation = psapi_DLL.NewProc("GetModuleInformation")
)

// Base: https://learn.microsoft.com/en-us/windows/win32/api/psapi/nf-psapi-getmoduleinformation
func GetModuleInformation(hProcess uintptr, hModule uintptr, lpModInfo *MODULEINFO, cb uint32) bool {
	ret, _, _ := procGetModuleInformation.Call(
		uintptr(hProcess),
		uintptr(hModule),
		uintptr(unsafe.Pointer(lpModInfo)),
		uintptr(cb))

	return uint(ret) != 0 // non-zero on success
}
