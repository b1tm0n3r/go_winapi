package main

import (
	"syscall"
	"unsafe"
)

var (
	kernel32_DLL            = syscall.NewLazyDLL("kernel32.dll")
	procVirtualAllocEx      = kernel32_DLL.NewProc("VirtualAllocEx")
	procCreateThread        = kernel32_DLL.NewProc("CreateThread")
	procGetCurrentProcess   = kernel32_DLL.NewProc("GetCurrentProcess")
	procWriteProcessMemory  = kernel32_DLL.NewProc("WriteProcessMemory")
	procWaitForSingleObject = kernel32_DLL.NewProc("WaitForSingleObject")
)

// Base: https://learn.microsoft.com/en-us/windows/win32/api/processthreadsapi/nf-processthreadsapi-createthread
func CreateThread(lpThreadAttributes uintptr, dwStackSize uint32, lpStartAddress uintptr,
	lpParameter uintptr, dwCreationFlags uint32, lpThreadId uintptr) uintptr {

	ret, _, _ := procCreateThread.Call(
		lpThreadAttributes,
		uintptr(dwStackSize),
		lpStartAddress,
		lpParameter,
		uintptr(dwCreationFlags),
		lpThreadId)

	return ret
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

// Base: https://learn.microsoft.com/en-us/windows/win32/api/memoryapi/nf-memoryapi-writeprocessmemory
func WriteProcessMemory(hProcess uintptr, lpBaseAddress uintptr, lpBuffer []byte,
	nSize uint32, lpNumberOfBytesWritten *uintptr) uint32 {

	ret, _, _ := procWriteProcessMemory.Call(
		hProcess,
		lpBaseAddress,
		uintptr(unsafe.Pointer(&lpBuffer[0])),
		uintptr(nSize),
		uintptr(unsafe.Pointer(lpNumberOfBytesWritten)))

	return uint32(ret) // return value should be non-zero
}

// Base: https://learn.microsoft.com/en-us/windows/win32/api/synchapi/nf-synchapi-waitforsingleobject
func WaitForSingleObject(hHandle uintptr, dwMilliseconds uint32) uint32 {
	ret, _, _ := procWaitForSingleObject.Call(
		hHandle,
		uintptr(dwMilliseconds))

	return uint32(ret)
}

// Base: https://learn.microsoft.com/en-us/windows/win32/api/processthreadsapi/nf-processthreadsapi-getcurrentprocess
func GetCurrentProcess() uintptr {
	ret, _, _ := procGetCurrentProcess.Call()

	return ret // Return should be non-null (returns process handle)
}
