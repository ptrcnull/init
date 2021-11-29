package main

import (
	"fmt"
	"io"
	"os"
)

var devices map[string]io.ReadWriteCloser

func GetDevice(name string) (io.ReadWriteCloser, error) {
	if dev, ok := devices[name]; ok {
		return dev, nil
	}

	dev, err := os.OpenFile("/dev/" + name, os.O_RDWR, 0644)
	if err != nil {
		return nil, fmt.Errorf("open: %w", err)
	}

	devices[name] = dev
	return dev, nil
}
