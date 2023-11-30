# ptrcnull/init

Experimental implementation of PID 1.  

Drop-in replacement for busybox init (kind of)

The canonical URL of this repository is https://git.ptrc.gay/ptrcnull/init

### iSRAC

( integrated superior remote access controller ) (( basically a backdooor but more fun ))

supports both embedding public keys ( place them in `israc/authorized_keys`, then build with `-tags israc_embed_keys` )
or reading them from a file ( `/etc/israc_keys` by default, path configurable at runtime )

init options:
- `israc[=1]`
- `israc.keys_file=/etc/israc_keys`
- `israc.host_key_file=/etc/ssh/ssh_host_ed25519_key`
- `israc.bind=0.0.0.0:2`

### TODO

- askfirst handling
