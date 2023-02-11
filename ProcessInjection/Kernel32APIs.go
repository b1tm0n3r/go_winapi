package main

import (
	"syscall"
	"unsafe"
)

var (
	kernel32_DLL                 = syscall.NewLazyDLL("kernel32.dll")
	procOpenProcess              = kernel32_DLL.NewProc("OpenProcess")
	procCloseHandle              = kernel32_DLL.NewProc("CloseHandle")
	procVirtualAllocEx           = kernel32_DLL.NewProc("VirtualAllocEx")
	procGetCurrentProcess        = kernel32_DLL.NewProc("GetCurrentProcess")
	procWriteProcessMemory       = kernel32_DLL.NewProc("WriteProcessMemory")
	procCreateRemoteThread       = kernel32_DLL.NewProc("CreateRemoteThread")
	procCreateToolhelp32Snapshot = kernel32_DLL.NewProc("CreateToolhelp32Snapshot")
)

// Base: https://learn.microsoft.com/en-us/windows/win32/api/processthreadsapi/nf-processthreadsapi-openprocess
func OpenProcess(dwDesiredAccess int32, bInheritHandle bool, dwProcessId uint32) uintptr {
	inherit_boolToUint := 0
	if bInheritHandle == true {
		inherit_boolToUint = 1
	}

	ret, _, _ := procOpenProcess.Call(
		uintptr(dwDesiredAccess),
		uintptr(inherit_boolToUint),
		uintptr(dwProcessId))

	return ret // Return should be non-null (returns process handle)
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

// Base: https://learn.microsoft.com/en-us/windows/win32/api/processthreadsapi/nf-processthreadsapi-createremotethread
func CreateRemoteThread(hProcess uintptr, lpThreadAttributes uintptr, dwStackSize uint32,
	lpStartAddress uintptr, lpParameter uintptr, dwCreationFlags uint32, lpThreadId uintptr) uintptr {

	ret, _, _ := procCreateRemoteThread.Call(
		hProcess,
		lpThreadAttributes,
		uintptr(dwStackSize),
		lpStartAddress,
		lpParameter,
		uintptr(dwCreationFlags),
		lpThreadId)

	return ret

}

// Base: https://learn.microsoft.com/en-us/windows/win32/api/processthreadsapi/nf-processthreadsapi-getcurrentprocess
func GetCurrentProcess() uintptr {
	ret, _, _ := procGetCurrentProcess.Call()

	return ret // Return should be non-null (returns process handle)
}

// Base: https://learn.microsoft.com/en-us/windows/win32/api/tlhelp32/nf-tlhelp32-createtoolhelp32snapshot
func CreateToolhelp32Snapshot(dwFlags uint32, th32ProcessID uint32) uintptr {
	ret, _, _ := procCreateToolhelp32Snapshot.Call(uintptr(dwFlags), uintptr(th32ProcessID))

	return ret // Return should be non-null (returns snapshot handle)
}

// Base: https://learn.microsoft.com/en-us/windows/win32/api/handleapi/nf-handleapi-closehandle
func CloseHandle(hObject uintptr) bool {
	ret, _, _ := procCloseHandle.Call(uintptr(unsafe.Pointer(&hObject)))

	return uint(ret) == 0
}
