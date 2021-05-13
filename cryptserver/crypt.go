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
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
)

func createHash(key string) []byte {
	hasher := md5.New()
	hasher.Write([]byte(key))
	return hasher.Sum(nil)
}

// Get md5 checksum value of input data
func md5Data(data []byte) string {
	hasher := md5.New()
	hasher.Write(data)

	//Get the 16 bytes hash
	hashInBytes := hasher.Sum(nil)[:16]

	//Convert the bytes to a string
	returnMD5String := hex.EncodeToString(hashInBytes)

	return returnMD5String
}

// Get md5 checksum value of input file
func md5File(filePath string) (string, error) {
	//Initialize variable returnMD5String now in case an error has to be returned
	var returnMD5String string

	//Open the passed argument and check for any error
	file, err := os.Open(filePath)
	if err != nil {
		return returnMD5String, err
	}

	//Tell the program to call the following function when the current function returns
	defer file.Close()

	//Open a new hash interface to write to
	hash := md5.New()

	//Copy the file in the hash interface and check for any error
	if _, err := io.Copy(hash, file); err != nil {
		return returnMD5String, err
	}

	//Get the 16 bytes hash
	hashInBytes := hash.Sum(nil)[:16]

	//Convert the bytes to a string
	returnMD5String = hex.EncodeToString(hashInBytes)

	return returnMD5String, nil
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
	fmt.Println("Old unencrypted file will be removed: ", inputFile)

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

	err = os.Remove(inputFile)
	if err != nil {
		log.Fatal(err)
	}
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
