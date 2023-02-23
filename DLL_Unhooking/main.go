package main

import (
	"fmt"
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

}
