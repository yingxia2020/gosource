package main

import (
	ecies "github.com/ecies/go"
	"log"
	"runtime"
)

func main() {
        runtime.GOMAXPROCS(2)
	k, err := ecies.GenerateKey()
	if err != nil {
		panic(err)
	}
	log.Println("key pair has been generated")

	for i := 0; i < 1000; i++ {
		ciphertext, err := ecies.Encrypt(k.PublicKey, []byte("THIS IS THE TEST"))
		if err != nil {
			panic(err)
		}
		if i == 0 {
			log.Printf("plaintext encrypted: %v\n", ciphertext)
		}

		plaintext, err := ecies.Decrypt(k, ciphertext)
		if err != nil {
			panic(err)
		}
		if i==0 {
			log.Printf("ciphertext decrypted: %s\n", string(plaintext))
		}
	}

}
