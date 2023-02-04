package main

import (
	"syscall"
)

var (
	kernel32_DLL = syscall.NewLazyDLL("kernel32.dll")
)
