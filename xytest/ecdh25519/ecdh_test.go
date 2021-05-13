package main

import (
	"bytes"
	"crypto/rand"
	"fmt"
	"testing"

	"github.com/aead/ecdh"
)

func ExampleX25519() {
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

func BenchmarkX25519(b *testing.B) {
	curve := ecdh.X25519()
	privateAlice, _, err := curve.GenerateKey(rand.Reader)
	if err != nil {
		b.Fatalf("Failed to generate Alice's private/public key pair: %s", err)
	}
	_, publicBob, err := curve.GenerateKey(rand.Reader)
	if err != nil {
		b.Fatalf("Failed to generate Bob's private/public key pair: %s", err)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		curve.ComputeSecret(privateAlice, publicBob)
	}
}
