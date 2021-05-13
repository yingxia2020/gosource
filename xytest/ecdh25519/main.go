package main

import (
	"bytes"
	"crypto/rand"
	"fmt"

	"github.com/aead/ecdh"
)

func main() {
	for i:=0; i<10000; i++ {
		exampleX25519()
	}
}

func exampleX25519() {
	c25519 := ecdh.X25519()

	privateAlice, publicAlice, err := c25519.GenerateKey(rand.Reader)
	if err != nil {
		fmt.Printf("Failed to generate Alice's private/public key pair: %s\n", err)
	}
	privateBob, publicBob, err := c25519.GenerateKey(rand.Reader)
	if err != nil {
		fmt.Printf("Failed to generate Bob's private/public key pair: %s\n", err)
	}

	if err := c25519.Check(publicBob); err != nil {
		fmt.Printf("Bob's public key is not on the curve: %s\n", err)
	}
	secretAlice := c25519.ComputeSecret(privateAlice, publicBob)

	if err := c25519.Check(publicAlice); err != nil {
		fmt.Printf("Alice's public key is not on the curve: %s\n", err)
	}
	secretBob := c25519.ComputeSecret(privateBob, publicAlice)

	if !bytes.Equal(secretAlice, secretBob) {
		fmt.Printf("key exchange failed - secret X coordinates not equal\n")
	}
}

