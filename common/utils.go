package common

import (
	"os"
	"strings"
)

func Getopt(name string, def string) string {
	for _, arg := range os.Args {
		if strings.HasPrefix(arg, name) {
			if arg == name {
				return "1"
			}
			splat := strings.SplitN(arg, "=", 2)
			if splat[0] == name {
				return splat[1]
			}
		}
	}
	return def
}
