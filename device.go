package main

import (
	"fmt"
	"os"
)

var devices = map[string]*os.File{}

func GetDevice(name string) (*os.File, error) {
	if dev, ok := devices[name]; ok {
		return dev, nil
	}

	dev, err := os.OpenFile("/dev/"+name, os.O_RDWR, 0644)
	if err != nil {
		return nil, fmt.Errorf("open: %w", err)
	}

	devices[name] = dev
	return dev, nil
}
