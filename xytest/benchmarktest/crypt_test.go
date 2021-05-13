// from fib_test.go
package main

import (
	"bufio"
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"strings"
	"testing"

	ecies "github.com/ecies/go"
	"github.com/spacemonkeygo/openssl"
	"golang.org/x/crypto/scrypt"
)

const ChunkSize = 10 * 1024
var k *ecies.PrivateKey

func init() {
	k, _ = ecies.GenerateKey()
}

func Fib(n int) int {
	if n < 2 {
		return n
	}
	return Fib(n-1) + Fib(n-2)
}

func EncryptString(input string, key, iv []byte) ([]byte, error) {
	ctx, err := openssl.NewGCMEncryptionCipherCtx(len(key)*8, nil, key, iv)
	if err != nil {
		return nil, fmt.Errorf("Failed making GCM encryption ctx: %v", err)
	}

	reader := strings.NewReader(input)
	var buf bytes.Buffer
	writer := bufio.NewWriter(&buf)
	chunk := make([]byte, ChunkSize)
	for {
		chunkSize, err := reader.Read(chunk)
		if err == io.EOF || chunkSize == 0 {
			break
		} else if err != nil {
			return nil, fmt.Errorf("Failed to read a chunk: %v", err)
		}

		encData, err := ctx.EncryptUpdate(chunk[:chunkSize])
		if err != nil {
			return nil, fmt.Errorf("Failed to perform an encryption: %v", err)
		}

		if _, err := writer.Write(encData); err != nil {
			return nil, fmt.Errorf("Failed to write an encrypted data: %v", err)
		}
	}

	encData, err := ctx.EncryptFinal()
	if err != nil {
		return nil, fmt.Errorf("Failed to finalize encryption: %v", err)
	}
	if _, err := writer.Write(encData); err != nil {
		return nil, fmt.Errorf("Failed to write a final encrypted data: %v", err)
	}

	tag, err := ctx.GetTag()
	if err != nil {
		return nil, fmt.Errorf("Failed to get GCM tag: %v", err)
	}
	if _, err := writer.Write(tag); err != nil {
		return nil, fmt.Errorf("Failed to write a gcm tag: %v", err)
	}
	writer.Flush()

	return buf.Bytes(), nil
}

func cipher() {
	key, _ := hex.DecodeString("1234567890ABCDEF1234567890ABCDEF")
	iv, _ := hex.DecodeString("1234567890ABCDEF")

	input := "hello worldhello worldhello worldhello worldhello worldhello worldhello worldhello worldhello worldhello worldhello worldhello worldhello worldhello worldhello worldhello worldhello worldhello world"
	EncryptString(input, key, iv)

}

func password(p string) string {
	salt := make([]byte, 10)
	rand.Read(salt)
	key, _ := scrypt.Key([]byte(p), salt, 32768, 8, 1, 32)
	return fmt.Sprintf("%x_%x", salt, key)
}

func encryptDecrypt() {
	ciphertext, _ := ecies.Encrypt(k.PublicKey, []byte("THIS IS THE TEST"))
	_, _ = ecies.Decrypt(k, ciphertext)
}

func BenchmarkEcies(b *testing.B) {
	for n := 0; n < b.N; n++ {
		encryptDecrypt()
	}
}

/*
func BenchmarkFib10(b *testing.B) {
	// run the Fib function b.N times
	for n := 0; n < b.N; n++ {
		Fib(10)
	}
}
func BenchmarkCipher(b *testing.B) {
	for n := 0; n < b.N; n++ {
		cipher()
	}
}
func BenchmarkPassword(b *testing.B) {
	for n := 0; n < b.N; n++ {
		password("hello")
	}
}
*/
