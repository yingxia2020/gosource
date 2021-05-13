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
	"bytes"
	"io/ioutil"
	"os"
	"testing"
)

const (
	passPhase = "xy123"
	input     = "I am a secret string"
	origFile  = "crypt.go"
	encFile   = "encFile"
	decFile   = "decFile"
)

var hash = []byte{252, 55, 206, 133, 137, 30, 121, 129, 116, 231, 25, 204, 47, 50, 149, 198, 54, 73, 139, 42, 149, 48, 28, 211, 248, 224, 119, 177, 28, 80, 87, 2}

func TestEncodeDecode(t *testing.T) {
	result := encryptData([]byte(input), passPhase)

	if string(decryptData(result, passPhase)) != input {
		t.Fatal("Encode/Decode test failed with unexpected result.")
	}
}

func TestCreateHash(t *testing.T) {
	result := createHash(passPhase)

	if len(result) != 32 {
		t.Fatalf("Invalid sha256 hash value with length = %d", len(result))
	}

	if bytes.Compare(result, hash) != 0 {
		t.Fatalf("Invalid sha256 hash value")
	}
}

func TestEncodeDecodeFile(t *testing.T) {
	encryptFile(origFile, encFile, passPhase)
	if _, err := os.Stat(encFile); os.IsNotExist(err) {
		t.Fatal("Encode file test failed without output file")
	}

	decryptFile(encFile, decFile, passPhase)
	if _, err := os.Stat(decFile); os.IsNotExist(err) {
		t.Fatal("Decode file test failed without output file")
	}

	origData, err := ioutil.ReadFile(origFile)
	if err != nil {
		t.Fatal("Failed to read original file to encode")
	}

	decData, err := ioutil.ReadFile(decFile)
	if err != nil {
		t.Fatal("Failed to read decoded file")
	}

	if bytes.Compare(origData, decData) != 0 {
		t.Fatal("Encode/Decode file test failed with unexpected result.")
	}

	// Clean up after tests
	os.Remove(encFile)
	os.Remove(decFile)
}
