package main

import (
	"fmt"
	"os"
	"sync"
)

type DeviceMap struct {
	sync.Mutex
	m map[string]*os.File
}

var devices = DeviceMap{
	m: map[string]*os.File{},
}

func GetDevice(name string) (*os.File, error) {
	if dev, ok := devices.m[name]; ok {
		return dev, nil
	}

	dev, err := os.OpenFile("/dev/"+name, os.O_RDWR, 0644)
	if err != nil {
		return nil, fmt.Errorf("open: %w", err)
	}

	devices.Lock()
	devices.m[name] = dev
	devices.Unlock()
	return dev, nil
}
