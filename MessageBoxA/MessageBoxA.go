package main

import (
	"syscall"
	"unsafe"
)

var (
	user32DLL  = syscall.NewLazyDLL("user32.dll")
	procMsgBox = user32DLL.NewProc("MessageBoxA")
)

func MessageBoxA(hwnd uint, lpCaption string, lpText string, uType uint) uint {
	title_ptr, _ := syscall.BytePtrFromString(lpCaption)
	content_ptr, _ := syscall.BytePtrFromString(lpText)

	ret, _, _ := procMsgBox.Call(
		0,
		uintptr(unsafe.Pointer(title_ptr)),
		uintptr(unsafe.Pointer(content_ptr)),
		0)
	return uint(ret)
}
