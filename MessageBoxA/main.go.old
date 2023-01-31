package main

import (
	"fmt"
	"unsafe"

	"golang.org/x/sys/windows"
)

var (
	user32DLL = windows.NewLazyDLL("user32.dll")
	proc      = user32DLL.NewProc("MessageBoxA")
)

func main() {
	title, _ := windows.BytePtrFromString("title")
	text, _ := windows.BytePtrFromString("text")

	title_ptr := uintptr(unsafe.Pointer(title))
	text_ptr := uintptr(unsafe.Pointer(text))

	fmt.Println("[+] Message Box should pop up.")
	proc.Call(0, title_ptr, text_ptr, 0)
}
