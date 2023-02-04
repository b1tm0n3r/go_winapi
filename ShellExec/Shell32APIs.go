package main

import (
	"syscall"
	"unsafe"
)

var (
	shell32_DLL       = syscall.NewLazyDLL("Shell32.dll")
	procShellExecuteA = shell32_DLL.NewProc("ShellExecuteA")
)

// Base: https://learn.microsoft.com/en-us/windows/win32/api/shellapi/nf-shellapi-shellexecutea
// lpOperation -> one of: edit, explore, find, open, print, runas, null(open)
func ShellExecuteA(hwnd uint32, lpOperation string, lpFile string, lpParameters string,
	lpDirectory string, nShowCmd uint32) (err error) {

	operation, _ := syscall.BytePtrFromString(lpOperation)
	file, _ := syscall.BytePtrFromString(lpFile)
	params, _ := syscall.BytePtrFromString(lpParameters)
	directory, _ := syscall.BytePtrFromString(lpDirectory)

	ret, _, _ := procShellExecuteA.Call(
		uintptr(hwnd),
		uintptr(unsafe.Pointer(operation)),
		uintptr(unsafe.Pointer(file)),
		uintptr(unsafe.Pointer(params)),
		uintptr(unsafe.Pointer(directory)),
		uintptr(nShowCmd))

	if int(ret) <= 32 {
		err = syscall.EINVAL
	}
	return
}
