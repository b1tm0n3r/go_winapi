package main

import (
	"fmt"
)

func main() {

	procArrSize := 1024
	procIdArr := make([]uint32, procArrSize)
	var lpcbNeeded uint32 = 0

	if EnumProcesses(procIdArr, 255, &lpcbNeeded) == 0 {
		fmt.Println("[!] Failed to read the processes. Leaving...")
		return
	}

	fmt.Println("[+] Listing processes ids...")
	for _, p := range procIdArr[:lpcbNeeded/4] {
		fmt.Println(p)
	}

	fmt.Println("[+] Looking for specific process...")
	// TODO: implement -> process lookup by name + obtaining proc handle with OpenProcess from kernel32 + using MiniDumpWriteDump from Dbghelp

}
