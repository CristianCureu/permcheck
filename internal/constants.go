package internal

import "os"

const (
	InsecureTag = "[!]"
	SecureTag   = "[✓]"
)

const (
	ColorReset = "\033[0m"
	ColorRed   = "\033[31m"
	ColorGreen = "\033[32m"
)

var SensitiveFiles = map[string]os.FileMode{
	".env":               0o600, // only owner can read/write
	"id_rsa":             0o600,
	"config.yaml":        0o640, // readable by group, writeable by owner
	"id_dsa":             0o600,
	"id_ecdsa":           0o600,
	"id_ed25519":         0o600,
	"config.yml":         0o640,
	"docker-compose.yml": 0o640,
	".git-credentials":   0o600,
	"authorized_keys":    0o600,
}

var SensitiveExtensions = map[string]os.FileMode{
	".pem": 0o600,
	".key": 0o600,
	".crt": 0o644,
}
