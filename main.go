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
			fmt.Printf("error closing inittab: %s\n", err)
		}
	} else {
		fmt.Printf("error reading inittab: %s\n", err)
	}

	for _, entry := range inittab.Entries(SysInit) {
		err := Exec(entry)
		if err != nil {
			fmt.Printf("error running sysinit \"%s\": %s\n", entry.Process, err)
		}
	}

	for _, entry := range inittab.Entries(Wait) {
		err := Exec(entry)
		if err != nil {
			fmt.Printf("error running wait \"%s\": %s\n", entry.Process, err)
		}
	}

	for _, entry := range inittab.Entries(Once) {
		_, err := Spawn(entry)
		if err != nil {
			fmt.Printf("error running once \"%s\": %s\n", entry.Process, err)
		}
	}

	// TODO implement AskFirst handling

	for _, entry := range inittab.Entries(Respawn) {
		go func(entry InitTabEntry) {
			for {
				err := Exec(entry)
				if err != nil {
					fmt.Printf("error running respawn \"%s\": %s\n", entry.Process, err)
					break
				}
			}
		}(entry)
	}

	// TODO implement Shutdown, Restart and CtrlAltDel handling

	select {}
}
