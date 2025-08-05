package main

import (
	"crypto/rand"
	"encoding/base64"
	"flag"
	"fmt"
	"github.com/qrv1t9/webdav/internal/config"
	"golang.org/x/crypto/argon2"
	"log"
)

func main() {

	path := flag.String("config", "", "path to config file")
	password := flag.String("pass", "", "your password")
	flag.Parse()

	cfg := config.MustLoad(*path)

	if *password == "" {
		log.Fatal("password is required")
	}

	salt, err := generateCryptographicSalt(cfg.Argon.ArgonSaltLength)
	if err != nil {
		log.Fatal(err)
	}

	hash := argon2.IDKey([]byte(*password), salt, cfg.Argon.ArgonIterations, cfg.Argon.ArgonMemory, cfg.Argon.ArgonParallelism, cfg.Argon.ArgonKeyLength)

	encodedSalt := base64.RawStdEncoding.EncodeToString(salt)
	encodedHash := base64.RawStdEncoding.EncodeToString(hash)

	fmt.Printf("$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s\n", argon2.Version, cfg.Argon.ArgonMemory, cfg.Argon.ArgonIterations, cfg.Argon.ArgonParallelism, encodedSalt, encodedHash)
}

func generateCryptographicSalt(saltSize uint32) ([]byte, error) {
	salt := make([]byte, saltSize)
	_, err := rand.Read(salt)
	if err != nil {
		return nil, fmt.Errorf("salt generation failed: %w", err)
	}
	return salt, nil
}
