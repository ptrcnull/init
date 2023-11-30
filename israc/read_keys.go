//go:build !israc_embed_keys

package israc

import (
	"fmt"
	"os"

	"github.com/gliderlabs/ssh"
	"github.com/ptrcnull/init/common"
)

func GetAuthorizedKeys() ([]ssh.PublicKey, error) {
	path := common.Getopt("israc.keys_file", "/etc/israc_keys")
	file, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("cannot read %s: %w", path, err)
	}

	return ParseAuthorizedKeys(string(file)), nil
}
