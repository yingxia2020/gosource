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
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

var (
	key = []byte("1234567890ABCDEF1234567890ABCDEF")
	iv  = []byte("1234567890ABCDEF")
)

var crypter *Crypter

func init() {
	var err error
	crypter, err = NewCrypter(key, iv)
	if err != nil {
		log.Fatal(err)
	}
}

func encryptData(data []byte) []byte {
	encryptedData, err := crypter.Encrypt(data)
	if err != nil {
		log.Fatal(err)
	}

	return encryptedData
}

func decryptData(data []byte) []byte {
	decryptedData, err := crypter.Decrypt(data)
	if err != nil {
		log.Fatal(err)
	}
	return decryptedData
}

func encryptFile(inputFile string, outputFile string) {
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

	bytesWritten, err := newFile.Write(encryptData(dataToEncrypt))
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

func decryptFile(encryptedFile string, decryptedFile string) {
	dataToDecrypt, err := ioutil.ReadFile(encryptedFile)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("New unencrypted file will be created: ", decryptedFile)

	newFile, err := os.Create(decryptedFile)
	if err != nil {
		log.Fatal(err)
	}

	dataDecrypted := decryptData(dataToDecrypt)
	if err != nil {
		log.Fatal(err)
	}
	bytesWritten, err := newFile.Write(dataDecrypted)
	fmt.Println(bytesWritten, "bytes written")
	newFile.Close()
}

func printUsage() {
	fmt.Println("Usage example: opensslcrypt -op \"enc\" or \"dec\" -i \"filename\" -o \"newfilename\"")
}

func main() {
	argSlice := os.Args
	if len(argSlice) <= 1 {
		fmt.Println("No arguments provided")
		printUsage()
		os.Exit(1)
	}

	var operation string
	flag.StringVar(&operation, "op", "", "Operation to perform, \"enc\" for encryption and \"dec\" for decryption")
	var input string
	flag.StringVar(&input, "i", "", "Name of input file")
	var output string
	flag.StringVar(&output, "o", "", "Name of output file")
	flag.Parse()

	fmt.Println(operation)
	if operation == "enc" {
		encryptFile(input, output)
	} else if operation == "dec" {
		decryptFile(input, output)
	} else {
		fmt.Println("Operation not specified")
		printUsage()
		os.Exit(1)
	}
}
