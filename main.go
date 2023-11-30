package main

import (
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

func main() {
	defer func() {
		err := recover()
		if err != nil {
			fmt.Printf("====================\npanic caught in main: %s\n====================\n", err)
		}
	}()

	if os.Getpid() != 1 {
		fmt.Printf("%s should only be ran as pid1, exiting\n", os.Args[0])
		os.Exit(1)
	}

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

	inittab.Entries(SysInit).ExecAll()
	inittab.Entries(Wait).ExecAll()
	inittab.Entries(Once).SpawnAll()
	inittab.Entries(Respawn).RespawnAll()

	// TODO implement AskFirst handling

	sigs := make(chan os.Signal, 1)

	go func() {
		for {
			sig := <-sigs
			switch sig {
			case syscall.SIGUSR2:
				// shutdown
				inittab.Entries(Shutdown).ExecAll()
				syscall.Reboot(syscall.LINUX_REBOOT_CMD_POWER_OFF)

			case syscall.SIGTERM:
				// reboot
				inittab.Entries(Shutdown).ExecAll()
				syscall.Reboot(syscall.LINUX_REBOOT_CMD_RESTART)

			case syscall.SIGQUIT:
				inittab.Entries(Shutdown).ExecAll()
				restart := inittab.Entries(Restart)[0]
				cmdline := strings.Split(restart.Process, " ")
				syscall.Exec(cmdline[0], cmdline[1:], []string{})

			case syscall.SIGINT:
				inittab.Entries(CtrlAltDel).ExecAll()
			}
		}
	}()

	signal.Notify(sigs, syscall.SIGUSR2, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGINT)

	go func() {
		for {
			syscall.Wait4(-1, nil, 0, nil)
		}
	}()

	select {}
}

func (i InitTab) ExecAll() {
	for _, entry := range i {
		err := Exec(entry)
		if err != nil {
			fmt.Printf("error executing \"%s\": %s\n", entry.Process, err)
		}
	}
}

func (i InitTab) SpawnAll() {
	for _, entry := range i {
		_, err := Spawn(entry)
		if err != nil {
			fmt.Printf("error spawning \"%s\": %s\n", entry.Process, err)
		}
	}
}

func (i InitTab) RespawnAll() {
	for _, entry := range i {
		go func(entry InitTabEntry) {
			defer func() {
				err := recover()
				if err != nil {
					fmt.Printf("====================\npanic caught in RespawnAll: %s\n====================\n", err)
				}
			}()

			for {
				err := Exec(entry)
				if err != nil {
					fmt.Printf("error respawning \"%s\": %s\n", entry.Process, err)
					break
				}
			}
		}(entry)
	}
}
