package main

import (
	"flag"
	"fmt"
	"log"
	"runtime"

	"SecureNest/internal/vault"
)

func main() {
	mode := flag.String("mode", "", "encrypt or decrypt")
	dir := flag.String("dir", "", "target directory (recursive)")
	flag.Parse()

	if *mode != "encrypt" && *mode != "decrypt" {
		log.Fatal("Specify -mode encrypt or decrypt")
	}
	if *dir == "" {
		log.Fatal("Specify -dir")
	}

	fmt.Print("Password: ")
	password, err := vault.ReadPassword()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println()

	// Argon2パラメータ
	params := vault.Argon2Params{
		Time:    3,
		Memory:  64 * 1024, // 64MB
		Threads: uint8(runtime.NumCPU()),
	}

	if err := vault.ProcessFolderRecursive(*mode, *dir, password, params); err != nil {
		log.Fatal(err)
	}

	zeroMemory(password)
	runtime.KeepAlive(password)

	fmt.Println("Done.")
}

//go:noinline
func zeroMemory(data []byte) {
	for i := range data {
		data[i] = 0
	}
}
