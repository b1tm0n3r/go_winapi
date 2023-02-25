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
}
