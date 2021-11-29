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

	// TODO: add stdio handling
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

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
