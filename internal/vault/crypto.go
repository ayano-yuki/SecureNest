package vault

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/binary"
	"io"

	"golang.org/x/crypto/pbkdf2"
)

const (
	SaltSize = 16
	KeySize  = 32
	NonceSize = 12
)

func deriveKey(password, salt []byte, iterations int) []byte {
	return pbkdf2.Key(password, salt, iterations, KeySize, sha256.New)
}

func generateSalt() ([]byte, error) {
	salt := make([]byte, SaltSize)
	if _, err := rand.Read(salt); err != nil {
		return nil, err
	}
	return salt, nil
}

func generateNonce() ([]byte, error) {
	nonce := make([]byte, NonceSize)
	if _, err := rand.Read(nonce); err != nil {
		return nil, err
	}
	return nonce, nil
}

func writeUint32(w io.Writer, val uint32) error {
	buf := make([]byte, 4)
	binary.BigEndian.PutUint32(buf, val)
	_, err := w.Write(buf)
	return err
}

func readUint32(r io.Reader) (uint32, error) {
	buf := make([]byte, 4)
	if _, err := io.ReadFull(r, buf); err != nil {
		return 0, err
	}
	return binary.BigEndian.Uint32(buf), nil
}
