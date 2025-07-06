package vault

import (
	"golang.org/x/crypto/argon2"
)

type Argon2Params struct {
	Time    uint32
	Memory  uint32
	Threads uint8
}

const (
	SaltSize  = 16
	KeySize   = 32
	NonceSize = 12
)

func deriveKey(password, salt []byte, p Argon2Params) []byte {
	return argon2.IDKey(password, salt, p.Time, p.Memory, p.Threads, KeySize)
}
