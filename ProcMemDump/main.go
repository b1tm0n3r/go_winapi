package main

import (
	"fmt"
)

func main() {

	procArrSize := 1024
	procIdArr := make([]uint32, procArrSize)
	var lpcbNeeded uint32 = 0

	fmt.Println("[i] Execute with elevated privileges!")
	if EnumProcesses(procIdArr, 1024, &lpcbNeeded) == 0 {
		fmt.Println("[!] Failed to read the processes. Leaving...")
		return
	}

	fmt.Println("[+] Listing processes ids...")
	for _, p := range procIdArr[:lpcbNeeded/4] {
		fmt.Println(p)
	}

	fmt.Println("[+] Looking for specific process...")

	var procAccess int32 = 0x1F0FFF
	for _, p := range procIdArr[:lpcbNeeded/4] {
		if p == 0 {
			continue
		}
		fmt.Println("[+] Opening process handle...")
		var procHandle uintptr = OpenProcess(procAccess, false, p)
		if procHandle == 0 {
			fmt.Printf("[!] Could not open handle to process with PID: %d\n", p)
			continue
		}
		fmt.Printf("[+] Obtained handle: %d for PID: %d\n", procHandle, p)

		fileNameMaxSize := 256
		fileNameBuff := make([]byte, fileNameMaxSize)

		res := GetProcessImageFileNameA(procHandle, fileNameBuff, 256)
		var resolvedProcessFileName string
		if res != 0 {
			resolvedProcessFileName = string(fileNameBuff[:])
			fmt.Printf("[+] Found process fileName: %s\n", resolvedProcessFileName)
		}

		fmt.Printf("[+] Closing handle: %d\n", procHandle)
		if !CloseHandle(procHandle) {
			fmt.Printf("[+] Could not close process handle: %d\n", procHandle)
		} else {
			fmt.Printf("[+] Successfully closed process handle.\n")
		}
	}

}
