package israc

import (
	"fmt"
	"strings"

	"github.com/gliderlabs/ssh"
)

func ParseAuthorizedKeys(in string) []ssh.PublicKey {
	var keys []ssh.PublicKey

	for _, rawKey := range strings.Split(in, "\n") {
		if len(rawKey) > 10 && !strings.HasPrefix(rawKey, "#") {
			key, _, _, _, err := ssh.ParseAuthorizedKey([]byte(rawKey))
			if err != nil {
				fmt.Printf("failed to parse ssh key '%s': %s\n", rawKey, err.Error())
				continue
			}
			keys = append(keys, key)
		}
	}

	return keys
}
