/* Copyright (C) Intel Corporation
 *
 * All Rights Reserved
 *
 * Unauthorized copying of this file, via any medium is strictly prohibited
 * Proprietary and confidential
 *
 * Written by Ying Xia <ying.xia@intel.com>, 2019
 */

package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
)

func createHash(key string) []byte {
	hash := sha256.Sum256([]byte(key))
	return hash[:]
}

func encryptData(data []byte, passphrase string) []byte {
	block, err := aes.NewCipher(createHash(passphrase))
	if err != nil {
		log.Fatal(err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		log.Fatal(err)
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		log.Fatal(err)
	}

	encryptedData := gcm.Seal(nonce, nonce, data, nil)
	return encryptedData
}

func decryptData(data []byte, passphrase string) []byte {
	key := createHash(passphrase)
	block, err := aes.NewCipher(key)
	if err != nil {
		log.Fatal(err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		log.Fatal(err)
	}

	nonceSize := gcm.NonceSize()
	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	decryptedData, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		log.Fatal(err)
	}
	return decryptedData
}

func encryptFile(inputFile string, outputFile string, passphrase string) {
	dataToEncrypt, err := ioutil.ReadFile(inputFile)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("New encrypted file will be created: ", outputFile)

	newFile, err := os.Create(outputFile)
	if err != nil {
		log.Fatal(err)
	}

	bytesWritten, err := newFile.Write(encryptData(dataToEncrypt, passphrase))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(bytesWritten, "bytes Written")

	newFile.Close()
}

func decryptFile(encryptedFile string, decryptedFile string, passphrase string) {
	dataToDecrypt, err := ioutil.ReadFile(encryptedFile)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("New unencrypted file will be created: ", decryptedFile)

	newFile, err := os.Create(decryptedFile)
	if err != nil {
		log.Fatal(err)
	}

	dataDecrypted := decryptData(dataToDecrypt, passphrase)
	if err != nil {
		log.Fatal(err)
	}
	bytesWritten, err := newFile.Write(dataDecrypted)
	fmt.Println(bytesWritten, "bytes written")
	newFile.Close()
}
