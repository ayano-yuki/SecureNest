package vault

import (
	"crypto/chacha20poly1305"
	"errors"
	"io"
	"os"
)

func DecryptFile(inputPath, outputPath string, password []byte) error {
	in, err := os.Open(inputPath)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer out.Close()

	salt := make([]byte, SaltSize)
	nonce := make([]byte, NonceSize)
	if _, err = io.ReadFull(in, salt); err != nil {
		return err
	}
	if _, err = io.ReadFull(in, nonce); err != nil {
		return err
	}
	time, err := readUint32(in)
	if err != nil {
		return err
	}
	memory, err := readUint32(in)
	if err != nil {
		return err
	}
	threadBuf := make([]byte, 1)
	if _, err = io.ReadFull(in, threadBuf); err != nil {
		return err
	}
	threads := threadBuf[0]

	p := Argon2Params{
		Time:    time,
		Memory:  memory,
		Threads: threads,
	}

	key := deriveKey(password, salt, p)

	aead, err := chacha20poly1305.New(key)
	if err != nil {
		return err
	}

	ciphertext, err := io.ReadAll(in)
	if err != nil {
		return err
	}

	plaintext, err := aead.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return errors.New("decryption failed (wrong password or corrupted data)")
	}

	_, err = out.Write(plaintext)
	return err
}
