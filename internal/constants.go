package internal

import "os"

const (
	DefaultInsecureTag = "[!]"
	DefaultSecureTag   = "[✓]"

	ColorInsecureTag = "\033[31m[!]\033[0m"
	ColorSecureTag   = "\033[32m[✓]\033[0m"
)

var SensitiveFiles = map[string]os.FileMode{
	".env":        0o600, // only owner can read/write
	"id_rsa":      0o600, // private key, very strict
	"config.yaml": 0o640, // readable by group, writeable by owner
}
