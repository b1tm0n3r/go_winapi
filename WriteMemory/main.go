package main

import (
	"bytes"
	"fmt"
	"syscall"
)

func main() {

	fmt.Println("[i] WriteMemory")

	fmt.Println("[+] Opening handle to current process...")

	hProcess := GetCurrentProcess()
	if hProcess == 0 {
		fmt.Println("[!] Could not obtain current process handle, aborting...")
		return
	}
	fmt.Printf("[+] Obtained process handle.\n")

	memAddr := VirtualAllocEx(hProcess, 0, 1024, (0x1000 | 0x2000), syscall.PAGE_EXECUTE_READWRITE)
	if memAddr == 0 {
		fmt.Printf("[!] Error allocating memory - returned value: %d\n", memAddr)
		return
	}
	fmt.Printf("[+] Successfully allocated memory at addr: 0x%x\n", memAddr)

	//fmt.Println("[dbg] addr ptr loc: ", memAddr)
	var valToWrite []byte = []byte{0x90, 0x90, 0x90, 0x90, 0x90, 0x90, 0x90, 0x90, 0x90, 0x90}
	var writtenSize uintptr = 0
	memWriteRes := WriteProcessMemory(hProcess, memAddr, valToWrite, uint32(len(valToWrite)), &writtenSize)
	if memWriteRes == 0 {
		fmt.Printf("[!] Error writing bytes: <%x> to address: <0x%x> in current process. Written size: %d, Returned value: %d\n", valToWrite, memAddr, writtenSize, memWriteRes)
		return
	}
	fmt.Printf("[+] Written %d bytes under address: 0x%x\n", writtenSize, memAddr)

	fmt.Printf("[+] Validating - Read from address into buffer...\n")
	var bytesFromMemory []byte = []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
	var readSize uintptr = 0
	memReadRes := ReadProcessMemory(hProcess, memAddr, bytesFromMemory, uint32(len(valToWrite)), &readSize)
	if !memReadRes {
		fmt.Printf("[!] Error read bytes: <%s> from address: <0x%x> in current process. Read size: %d, Returned value: %d\n", bytesFromMemory, memAddr, readSize, memReadRes)
		return
	}
	fmt.Printf("[+] Read %d bytes under address: 0x%x\n", readSize, memAddr)

	if bytes.Equal(valToWrite, bytesFromMemory) {
		fmt.Println("[+] Write/Read arrays store the same value\n")
	}

	CloseHandle(hProcess)

}
