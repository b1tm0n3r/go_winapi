package main

import (
	"fmt"
	"os"
	"strconv"
	"syscall"
)

func main() {

	fmt.Println("[i] Execute with elevated privileges!")
	fmt.Println("[i] Usage: appName.exe [<PID> <outFileName>]")

	if len(os.Args) == 1 { // no arg given - list processes, and try create handles for each one
		listAllProcesses()
	} else if len(os.Args) == 2 || len(os.Args) == 3 { // pid given to dump proc mem - dump memory of the process
		pid_proc, _ := strconv.ParseUint(os.Args[1], 10, 32)
		processPid := uint32(pid_proc)
		if len(os.Args) == 2 {
			dumpProcessMemoryToFile(processPid, "mem.dump")
		} else {
			dumpProcessMemoryToFile(processPid, os.Args[2])
		}
	}
}

func dumpProcessMemoryToFile(pid uint32, outFileName string) {

	fmt.Printf("[+] PID to dump memory from: %s", pid)

	var procAccess int32 = 0x1F0FFF
	var procHandle uintptr = OpenProcess(procAccess, false, pid)
	if procHandle == 0 {
		fmt.Printf("[!] Could not open handle to process with PID: %d\n", pid)
		return
	}
	fmt.Printf("[+] Obtained handle: %d for PID: %d\n", procHandle, pid)

	hFile := CreateFileA(outFileName,
		uint32(syscall.GENERIC_WRITE),
		uint32(syscall.FILE_SHARE_READ),
		uintptr(0),
		uint32(syscall.CREATE_ALWAYS),
		uint32(syscall.FILE_ATTRIBUTE_NORMAL),
		uintptr(0))
	if hFile == 0 {
		fmt.Printf("[!] Could not create a dump file!")
		return
	}
	fmt.Printf("[+] Created file and obtained handle.\n", hFile)
	fmt.Printf("[+] Dumping process memory to file: %s...\n", outFileName)
	MiniDumpWriteDump(procHandle, pid, hFile, 2, uintptr(0), uintptr(0), uintptr(0))
	CloseHandle(procHandle)
	CloseHandle(hFile)
	fmt.Printf("[+] Done!")
}

func listAllProcesses() {
	procArrSize := 1024
	procIdArr := make([]uint32, procArrSize)
	var lpcbNeeded uint32 = 0

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
