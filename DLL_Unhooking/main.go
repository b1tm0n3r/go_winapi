package main

import (
	"fmt"
	"strings"
	"syscall"
	"unsafe"
)

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

	fmt.Printf("[+] Base dll location: %x\n", modinfo.BaseOfDll)
	dosHeader := (*PIMAGE_DOS_HEADER)(unsafe.Pointer(modinfo.BaseOfDll))
	fmt.Printf("[+] DOS header e_lfanew val: %x\n", dosHeader.e_lfanew)
	fmt.Printf("[+] Base NtDll header location: %x\n", (modinfo.BaseOfDll + uintptr(dosHeader.e_lfanew)))
	ntdllHeader := (*PIMAGE_NT_HEADERS)(unsafe.Pointer(modinfo.BaseOfDll + uintptr(dosHeader.e_lfanew)))
	fmt.Printf("[+] NtDll header (optional) Magic: %x\n", ntdllHeader.OptionalHeader.Magic)

	fmt.Printf("[+] Iterating over number of sections: %d\n", ntdllHeader.FileHeader.NumberOfSections)
	for i := uint32(0); i < uint32(ntdllHeader.FileHeader.NumberOfSections); i++ {
		sectionHeader := (*PIMAGE_SECTION_HEADER)(unsafe.Pointer(uintptr(unsafe.Pointer(IMAGE_FIRST_SECTION(ntdllHeader))) + uintptr(IMAGE_SIZEOF_SECTION_HEADER*i) + uintptr(4))) // add 4 - probably bug somewhere earlier, it sets padding correctly
		fmt.Printf("[+] Processed section address: %x\n", unsafe.Pointer(sectionHeader))

		fmt.Printf("[+] Processed section name: %s\n", string(sectionHeader.Name[:]))
		fmt.Printf("[+] Processed section name: %x\n", sectionHeader.Name)

		if strings.HasPrefix(string(sectionHeader.Name[:]), ".text") {
			fmt.Printf("[+] Found .text section at address: %x\n", unsafe.Pointer(sectionHeader))

		}
	}

	CloseHandle(hNtdllFile)
	CloseHandle(hNtdllFileMapping)
	CloseHandle(hProcess)
}

func IMAGE_FIRST_SECTION(ntheader *PIMAGE_NT_HEADERS) *PIMAGE_SECTION_HEADER {
	// Calculate the address of the first section header
	fileHeaderSize := unsafe.Sizeof(ntheader.FileHeader)
	optionalHeaderSize := uintptr(ntheader.FileHeader.SizeOfOptionalHeader)
	firstSectionOffset := uintptr(unsafe.Pointer(ntheader)) + fileHeaderSize + optionalHeaderSize

	// Iterate over the section headers and find the first one
	sectionHeaderSize := unsafe.Sizeof(PIMAGE_SECTION_HEADER{})
	sectionHeaderPtr := (*PIMAGE_SECTION_HEADER)(unsafe.Pointer(firstSectionOffset))
	for i := 0; uint16(i) < ntheader.FileHeader.NumberOfSections; i++ {
		if sectionHeaderPtr.SizeOfRawData != 0 {
			return sectionHeaderPtr
		}
		sectionHeaderPtr = (*PIMAGE_SECTION_HEADER)(unsafe.Pointer(uintptr(unsafe.Pointer(sectionHeaderPtr)) + sectionHeaderSize))
	}

	return nil
}
