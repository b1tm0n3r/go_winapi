package main

import (
	"fmt"
)

func main() {
	fmt.Println("[+] Message Box should pop up.")
	MessageBoxA(0, "test", "text", 0)
}
