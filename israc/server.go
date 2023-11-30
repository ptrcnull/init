package israc

import (
	"fmt"
	"log"
	"os"

	_ "embed"

	"github.com/gliderlabs/ssh"
	"github.com/ptrcnull/init/common"
)

//go:embed hostkey_ed25519
var defaultHostKey []byte

func Start() {
	keys, err := GetAuthorizedKeys()
	if err != nil {
		fmt.Printf("get authorized keys: %s\n", err.Error())
		return
	}

	publicKeyOption := ssh.PublicKeyAuth(func(ctx ssh.Context, key ssh.PublicKey) bool {
		for _, otherkey := range keys {
			if ssh.KeysEqual(key, otherkey) {
				return true
			}
		}
		return false
	})

	ssh.Handle(Handle)

	hostKeyPath := common.Getopt("israc.host_key_file", "/etc/ssh/ssh_host_ed25519_key")

	var hostKeyOption ssh.Option
	if file, err := os.Open(hostKeyPath); err == nil {
		file.Close()
		hostKeyOption = ssh.HostKeyFile(hostKeyPath)
	} else {
		fmt.Println("using embedded host key")
		hostKeyOption = ssh.HostKeyPEM(defaultHostKey)
	}

	log.Fatal(ssh.ListenAndServe(common.Getopt("israc.bind", ":2222"), nil, publicKeyOption, hostKeyOption))
}
