package main

import (
	"fmt"
	"os"
)

func main() {

	if len(os.Args) != 2 {
		fmt.Println("[i] Usage: appName.exe \"<cmd_to_run>\"")
	}

	fmt.Printf("[+] Executing command: %s\n", os.Args[1])
	err := ShellExecuteA(0, "open", os.Args[1], "", "", 0)
	if err != nil {
		fmt.Println("[!] Error:", err)
	}
}
