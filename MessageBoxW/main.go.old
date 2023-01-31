package main

import (
	"fmt"
	"unsafe"

	"golang.org/x/sys/windows"
)

var (
	user32DLL = windows.NewLazyDLL("user32.dll")
	proc      = user32DLL.NewProc("MessageBoxW")
)

func main() {
	title, _ := windows.UTF16PtrFromString("title")
	text, _ := windows.UTF16PtrFromString("text")

	title_ptr := uintptr(unsafe.Pointer(title))
	text_ptr := uintptr(unsafe.Pointer(text))

	fmt.Println("[+] Message Box should pop up.")
	proc.Call(0, title_ptr, text_ptr, 0)
}
