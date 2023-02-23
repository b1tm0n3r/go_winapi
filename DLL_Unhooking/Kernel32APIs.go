package main

import (
	"syscall"
	"unsafe"
)

var (
	kernel32_DLL          = syscall.NewLazyDLL("kernel32.dll")
	procGetCurrentProcess = kernel32_DLL.NewProc("GetCurrentProcess")
	procGetModuleHandleA  = kernel32_DLL.NewProc("GetModuleHandleA")
	procCloseHandle       = kernel32_DLL.NewProc("CloseHandle")
)

// Base: https://learn.microsoft.com/en-us/windows/win32/api/processthreadsapi/nf-processthreadsapi-getcurrentprocess
func GetCurrentProcess() uintptr {
	ret, _, _ := procGetCurrentProcess.Call()

	return ret // Return should be non-null (returns process handle)
}

// Base: https://learn.microsoft.com/en-us/windows/win32/api/libloaderapi/nf-libloaderapi-getmodulehandlea
func GetModuleHandleA(lpModuleHandle string) uintptr {
	moduleHandle_ptr, _ := syscall.BytePtrFromString(lpModuleHandle)

	ret, _, _ := procGetModuleHandleA.Call(uintptr(unsafe.Pointer(moduleHandle_ptr)))

	return ret
}

// Base: https://learn.microsoft.com/en-us/windows/win32/api/handleapi/nf-handleapi-closehandle
func CloseHandle(hObject uintptr) bool {
	ret, _, _ := procCloseHandle.Call(uintptr(unsafe.Pointer(&hObject)))

	return uint(ret) == 0
}
