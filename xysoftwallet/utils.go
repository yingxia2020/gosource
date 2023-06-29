package main

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

func generatePublidId() string {
	randomString := generateRandomString(10)
	fmt.Println("Random String:", randomString)

	hash := sha256.Sum256([]byte(randomString))
	hashBytes := hash[:10]
	hashHex := hex.EncodeToString(hashBytes)

	fmt.Println("SHA-256 Digest (First 10 bytes in hex):", hashHex)
	return hashHex
}

func generateRandomString(length int) string {
	charSet := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	randomBytes := make([]byte, length)
	_, err := rand.Read(randomBytes)
	if err != nil {
		panic(err)
	}

	for i := 0; i < length; i++ {
		randomBytes[i] = charSet[randomBytes[i]%byte(len(charSet))]
	}

	return string(randomBytes)
}
