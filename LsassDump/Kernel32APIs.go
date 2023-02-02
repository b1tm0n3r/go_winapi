package main

import (
	"syscall"
	"unsafe"
)

var (
	kernel32_DLL                 = syscall.NewLazyDLL("kernel32.dll")
	procOpenProcess              = kernel32_DLL.NewProc("OpenProcess")
	procCloseHandle              = kernel32_DLL.NewProc("CloseHandle")
	procCreateFileA              = kernel32_DLL.NewProc("CreateFileA")
	procopenProcess              = kernel32_DLL.NewProc("OpenProcess")
	procGetCurrentProcess        = kernel32_DLL.NewProc("GetCurrentProcess")
	procCreateToolhelp32Snapshot = kernel32_DLL.NewProc("CreateToolhelp32Snapshot")
)

// Base: https://learn.microsoft.com/en-us/windows/win32/api/processthreadsapi/nf-processthreadsapi-openprocess
func OpenProcess(dwDesiredAccess int32, bInheritHandle bool, dwProcessId uint32) uintptr {
	inherit_boolToUint := 0
	if bInheritHandle {
		inherit_boolToUint = 1
	}

	ret, _, _ := procOpenProcess.Call(
		uintptr(dwDesiredAccess),
		uintptr(inherit_boolToUint),
		uintptr(dwProcessId))

	return ret // Return should be non-null (returns process handle)
}

// Base: https://learn.microsoft.com/en-us/windows/win32/api/handleapi/nf-handleapi-closehandle
func CloseHandle(hObject uintptr) bool {
	ret, _, _ := procCloseHandle.Call(uintptr(unsafe.Pointer(&hObject)))

	return uint(ret) == 0
}

// Base: https://learn.microsoft.com/en-us/windows/win32/api/fileapi/nf-fileapi-createfilea
func CreateFileA(lpFileName string, dwDesiredAccess uint32, dwShareMode uint32, lpSecurityAttributes uintptr,
	dwCreationDisposition uint32, dwFlagsAndAttributes uint32, hTemplateFile uintptr) uintptr {
	fileName_ptr, _ := syscall.BytePtrFromString(lpFileName)

	ret, _, _ := procCreateFileA.Call(
		uintptr(unsafe.Pointer(fileName_ptr)),
		uintptr(dwDesiredAccess),
		uintptr(dwShareMode),
		uintptr(unsafe.Pointer(&lpSecurityAttributes)),
		uintptr(dwCreationDisposition),
		uintptr(dwFlagsAndAttributes),
		uintptr(unsafe.Pointer(&hTemplateFile)))

	return ret // Return should be non-null (returns file handle)
}

// Base: https://learn.microsoft.com/en-us/windows/win32/api/processthreadsapi/nf-processthreadsapi-getcurrentprocessid
func GetCurrentProcess() uintptr {
	ret, _, _ := procGetCurrentProcess.Call()

	return ret // Return should be non-null (returns process handle)
}

// Base: https://learn.microsoft.com/en-us/windows/win32/api/tlhelp32/nf-tlhelp32-createtoolhelp32snapshot
func CreateToolhelp32Snapshot(dwFlags uint32, th32ProcessID uint32) uintptr {
	ret, _, _ := procCreateToolhelp32Snapshot.Call(uintptr(dwFlags), uintptr(th32ProcessID))

	return ret // Return should be non-null (returns snapshot handle)
}
