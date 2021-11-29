package main

import (
	"fmt"
	"os"
)

func main() {
	inittab := DefaultInitTab
	if file, err := os.OpenFile("/etc/inittab", os.O_RDONLY, 0644); err == nil {
		inittab = ParseInitTab(file)
		err := file.Close()
		if err != nil {
			fmt.Println("close inittab:", err)
		}
	}

}
