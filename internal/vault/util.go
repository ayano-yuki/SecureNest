package vault

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/term"
)

func ReadPassword() ([]byte, error) {
	return term.ReadPassword(int(os.Stdin.Fd()))
}

func ProcessFolderRecursive(mode string, rootDir string, password []byte, p Argon2Params) error {
	return filepath.WalkDir(rootDir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}

		var outputPath string
		if mode == "encrypt" {
			outputPath = path + ".vault"
		} else if mode == "decrypt" {
			if !strings.HasSuffix(path, ".vault") {
				return nil
			}
			outputPath = strings.TrimSuffix(path, ".vault")
		}

		fmt.Printf("[%s] %s -> %s\n", mode, path, outputPath)

		if mode == "encrypt" {
			err = EncryptFile(path, outputPath, password, p)
		} else {
			err = DecryptFile(path, outputPath, password)
		}
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error processing %s: %v\n", path, err)
			return nil
		}

		if removeErr := os.Remove(path); removeErr != nil {
			fmt.Fprintf(os.Stderr, "Error removing %s: %v\n", path, removeErr)
		}
		return nil
	})
}
