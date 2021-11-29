package main

import (
	"bufio"
	"io"
	"strings"
)

type Action uint8

const (
	SysInit Action = iota
	Wait
	Once
	Respawn
	AskFirst
	Shutdown
	Restart
	CtrlAltDel
)

var ActionMap = map[string]Action{
	"sysinit":    SysInit,
	"wait":       Wait,
	"once":       Once,
	"respawn":    Respawn,
	"askfirst":   AskFirst,
	"shutdown":   Shutdown,
	"restart":    Restart,
	"ctrlaltdel": CtrlAltDel,
}

type TabEntry struct {
	Device  string
	Action  Action
	Process string
}

var DefaultInittab = []TabEntry{
	{"", SysInit, "/etc/init.d/rcS"},
	{"", AskFirst, "/bin/sh"},
	{"", CtrlAltDel, "/sbin/reboot"},
	{"", Shutdown, "/sbin/swapoff -a"},
	{"", Shutdown, "/bin/umount -a -r"},
	{"", Restart, "/sbin/init"},
	{"tty2", AskFirst, "/bin/sh"},
	{"tty3", AskFirst, "/bin/sh"},
	{"tty4", AskFirst, "/bin/sh"},
}

func ParseInittab(reader io.Reader) []TabEntry {
	var res []TabEntry
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "#") {
			continue
		}

		tokens := strings.Split(line, ":")
		if len(tokens) != 4 {
			continue
		}

		action, ok := ActionMap[tokens[2]]
		if !ok {
			continue
		}

		res = append(res, TabEntry{
			Device:  tokens[0],
			Action:  action,
			Process: tokens[3],
		})
	}

	return res
}
