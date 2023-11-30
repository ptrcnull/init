//go:build israc_embed_keys

package israc

import (
	_ "embed"

	"github.com/gliderlabs/ssh"
)

//go:embed authorized_keys
var AuthorizedKeys string

func GetAuthorizedKeys() ([]ssh.PublicKey, error) {
	return ParseAuthorizedKeys(AuthorizedKeys), nil
}
