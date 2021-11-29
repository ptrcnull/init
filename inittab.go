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

type InitTabEntry struct {
	Device  string
	Action  Action
	Process string
}

type InitTab []InitTabEntry

func (i InitTab) Entries(action Action) InitTab {
	var res InitTab

	for _, entry := range i {
		if entry.Action == action {
			res = append(res, entry)
		}
	}

	return res
}

var DefaultInitTab = InitTab{
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

func ParseInitTab(reader io.Reader) InitTab {
	var res InitTab
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

		res = append(res, InitTabEntry{
			Device:  tokens[0],
			Action:  action,
			Process: tokens[3],
		})
	}

	return res
}
