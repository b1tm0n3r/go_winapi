package main

import (
	"syscall"
	"unsafe"
)

var (
	user32DLL  = syscall.NewLazyDLL("user32.dll")
	procMsgBox = user32DLL.NewProc("MessageBoxW")
)

func MessageBoxW(hwnd uint32, lpText string, lpCaption string, uType uint32) uint {
	title_ptr, _ := syscall.UTF16PtrFromString(lpCaption)
	content_ptr, _ := syscall.UTF16PtrFromString(lpText)

	ret, _, _ := procMsgBox.Call(
		uintptr(hwnd),
		uintptr(unsafe.Pointer(title_ptr)),
		uintptr(unsafe.Pointer(content_ptr)),
		uintptr(uType))
	return uint(ret)
}
