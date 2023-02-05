package main

import (
	"bytes"
	"fmt"
	"syscall"
)

func main() {

	fmt.Println("[i] WriteMemory ")

	fmt.Println("[+] Opening handle to current process...")

	hProcess := GetCurrentProcess()
	if hProcess == 0 {
		fmt.Println("[!] Could not obtain current process handle, aborting...")
		return
	}
	fmt.Printf("[+] Obtained process handle.\n")

	var allocMemSize uint32 = 1024
	memAddr := VirtualAllocEx(hProcess, 0, allocMemSize, (0x1000 | 0x2000), syscall.PAGE_EXECUTE_READWRITE)
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

	fmt.Println("[i] Validating")
	fmt.Println("[+] Read from address into buffer...")
	var bytesFromMemory []byte = []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
	var readSize uintptr = 0
	memReadRes := ReadProcessMemory(hProcess, memAddr, bytesFromMemory, uint32(len(valToWrite)), &readSize)
	if !memReadRes {
		fmt.Printf("[!] Error read bytes: <%s> from address: <0x%x> in current process. Read size: %d, Returned value: %d\n", bytesFromMemory, memAddr, readSize, memReadRes)
		return
	}
	fmt.Printf("[+] Read %d bytes under address: 0x%x\n", readSize, memAddr)

	//fmt.Scanln() // dbg test point
	if bytes.Equal(valToWrite, bytesFromMemory) {
		fmt.Printf("[+] Write/Read arrays store the same value: {%# 02x}\n", bytesFromMemory)
		fmt.Printf(" |-- This value was written under memory addr: 0x%x\n", memAddr)
	}

	fmt.Printf("[+] Calling VirtualFreeEx on allocated memory...\n")
	var MEM_RELEASE uint32 = 0x8000 // + memSize must be 0 when release + addr points to base addr returned on alloc
	freeRes := VirtualFreeEx(hProcess, memAddr, 0, MEM_RELEASE)
	if !freeRes {
		fmt.Printf("[!] Error when releasing the allocated memory")
		return
	}
	fmt.Printf("[+] Memory released")

	CloseHandle(hProcess)

}
