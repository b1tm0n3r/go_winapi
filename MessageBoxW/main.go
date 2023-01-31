package main

import (
	"fmt"
)

func main() {
	fmt.Println("[+] Message Box should pop up.")
	MessageBoxW(0, "test", "text", 0)
}
