package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func Spawn(entry InitTabEntry) (*exec.Cmd, error) {
	cmdline := strings.Split(entry.Process, " ")
	cmd := exec.Command(cmdline[0], cmdline[1:]...)

	stdio := os.Stdout
	if entry.Device != "" {
		dev, err := GetDevice(entry.Device)
		if err != nil {
			return nil, fmt.Errorf("open device %s: %w", entry.Device, err)
		}
		stdio = dev
	}
	cmd.Stdin = stdio
	cmd.Stdout = stdio
	cmd.Stderr = stdio

	err := cmd.Start()
	if err != nil {
		return nil, err
	}

	return cmd, nil
}

func Exec(entry InitTabEntry) error {
	cmd, err := Spawn(entry)
	if err != nil {
		return fmt.Errorf("spawn: %w", err)
	}

	// skipping error handling due to wait4 in main
	cmd.Wait()

	return nil
}
