package vault

import (
	"crypto/aes"
	"crypto/cipher"
	"io"
	"os"
)

func EncryptFile(inputPath, outputPath string, password []byte, iterations int) error {
	// 入力ファイル
	in, err := os.Open(inputPath)
	if err != nil {
		return err
	}
	defer in.Close()

	// 出力ファイル
	out, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Salt & Nonce
	salt, err := generateSalt()
	if err != nil {
		return err
	}
	nonce, err := generateNonce()
	if err != nil {
		return err
	}

	// Key
	key := deriveKey(password, salt, iterations)

	// AES-GCM
	block, err := aes.NewCipher(key)
	if err != nil {
		return err
	}
	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return err
	}

	// 入力データ
	plaintext, err := io.ReadAll(in)
	if err != nil {
		return err
	}
	ciphertext := aesgcm.Seal(nil, nonce, plaintext, nil)

	// ヘッダ書き込み
	if _, err = out.Write(salt); err != nil {
		return err
	}
	if _, err = out.Write(nonce); err != nil {
		return err
	}
	if err = writeUint32(out, uint32(iterations)); err != nil {
		return err
	}

	// 本体
	_, err = out.Write(ciphertext)
	return err
}
