package main

import (
	"syscall"
	"unsafe"
)

var (
	kernel32_DLL           = syscall.NewLazyDLL("kernel32.dll")
	procCreateFileA        = kernel32_DLL.NewProc("CreateFileA")
	procCreateFileMappingA = kernel32_DLL.NewProc("CreateFileMappingA")
	procGetCurrentProcess  = kernel32_DLL.NewProc("GetCurrentProcess")
	procGetModuleHandleA   = kernel32_DLL.NewProc("GetModuleHandleA")
	procMapViewOfFile      = kernel32_DLL.NewProc("MapViewOfFile")
	procVirtualProtect     = kernel32_DLL.NewProc("VirtualProtect")
	procCloseHandle        = kernel32_DLL.NewProc("CloseHandle")
	procCopyMemory         = kernel32_DLL.NewProc("RtlCopyMemory")
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

// Base: https://learn.microsoft.com/en-us/windows/win32/api/winbase/nf-winbase-createfilemappinga
func CreateFileMappingA(hFile uintptr, lpFileMappingAttributes uintptr, flProtect uint32,
	dwMaximumSizeHigh uint32, dwMaximumSizeLow uint32, lpName uintptr) uintptr {

	ret, _, _ := procCreateFileMappingA.Call(
		hFile,
		lpFileMappingAttributes,
		uintptr(flProtect),
		uintptr(dwMaximumSizeHigh),
		uintptr(dwMaximumSizeLow),
		lpName)

	return ret // Return should be non-null (returns file mapping handle)
}

// Base: https://learn.microsoft.com/en-us/windows/win32/api/memoryapi/nf-memoryapi-mapviewoffile
func MapViewOfFile(hFileMappingObject uintptr, dwDesiredAccess uint32, dwFileOffsetHigh uint32,
	dwFileOffsetLow uint32, dwNumberOfBytesToMap uint32) uintptr {

	ret, _, _ := procMapViewOfFile.Call(
		hFileMappingObject,
		uintptr(dwDesiredAccess),
		uintptr(dwFileOffsetHigh),
		uintptr(dwFileOffsetLow),
		uintptr(dwNumberOfBytesToMap))

	return ret // retrun should be non-null (starting address of mapped view)
}

// Base: https://learn.microsoft.com/en-us/windows/win32/api/memoryapi/nf-memoryapi-virtualprotect
func VirtualProtect(lpAddress uintptr, dwSize uint32, flNewProtect uint32, lpflOldProtect *uint32) bool {

	ret, _, _ := procVirtualProtect.Call(
		lpAddress,
		uintptr(dwSize),
		uintptr(flNewProtect),
		uintptr(unsafe.Pointer(lpflOldProtect)))

	return uint(ret) != 0 // return value is nonzero on success.
}

// Base: https://learn.microsoft.com/en-us/windows/win32/api/handleapi/nf-handleapi-closehandle
func CloseHandle(hObject uintptr) bool {
	ret, _, _ := procCloseHandle.Call(uintptr(unsafe.Pointer(&hObject)))

	return uint(ret) != 0 // return value is nonzero on success.
}

// Base: https://learn.microsoft.com/en-us/previous-versions/windows/desktop/legacy/aa366535(v=vs.85)
func CopyMemory(destination uintptr, source uintptr, length uint32) {

	procCopyMemory.Call(
		destination,
		source,
		uintptr(length))
}
