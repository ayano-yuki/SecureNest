package vault

import (
	"crypto/aes"
	"crypto/cipher"
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

	// ヘッダ
	salt := make([]byte, SaltSize)
	nonce := make([]byte, NonceSize)

	if _, err = io.ReadFull(in, salt); err != nil {
		return err
	}
	if _, err = io.ReadFull(in, nonce); err != nil {
		return err
	}
	iterations, err := readUint32(in)
	if err != nil {
		return err
	}

	// Key
	key := deriveKey(password, salt, int(iterations))

	// AES-GCM
	block, err := aes.NewCipher(key)
	if err != nil {
		return err
	}
	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return err
	}

	// 暗号データ
	ciphertext, err := io.ReadAll(in)
	if err != nil {
		return err
	}

	// 復号
	plaintext, err := aesgcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return errors.New("decryption failed (wrong password?)")
	}

	_, err = out.Write(plaintext)
	return err
}
