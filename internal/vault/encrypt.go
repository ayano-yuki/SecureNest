package vault

import (
	"crypto/chacha20poly1305"
	"io"
	"os"
)

func EncryptFile(inputPath, outputPath string, password []byte, p Argon2Params) error {
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

	salt, err := generateSalt()
	if err != nil {
		return err
	}
	nonce, err := generateNonce()
	if err != nil {
		return err
	}

	key := deriveKey(password, salt, p)

	aead, err := chacha20poly1305.New(key)
	if err != nil {
		return err
	}

	plaintext, err := io.ReadAll(in)
	if err != nil {
		return err
	}
	ciphertext := aead.Seal(nil, nonce, plaintext, nil)

	// ヘッダ
	if _, err = out.Write(salt); err != nil {
		return err
	}
	if _, err = out.Write(nonce); err != nil {
		return err
	}
	if err = writeUint32(out, p.Time); err != nil {
		return err
	}
	if err = writeUint32(out, p.Memory); err != nil {
		return err
	}
	if _, err = out.Write([]byte{p.Threads}); err != nil {
		return err
	}

	// 本体
	_, err = out.Write(ciphertext)
	return err
}
