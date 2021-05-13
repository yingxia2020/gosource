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
	"encoding/base64"
	"testing"
)

const (
	input  = "hello world"
	output = "RanFyUZSP9u/HLZjyI5zXQ=="
)

var (
	testkey = []byte("1234567890ABCDEF1234567890ABCDEF")
	testiv  = []byte("1234567890ABCDEF")
)

func TestEncodeDecode(t *testing.T) {
	crypter, _ := NewCrypter(testkey, testiv)

	encoded, _ := crypter.Encrypt([]byte("hello world"))
	if base64.StdEncoding.EncodeToString(encoded) != output {
		t.Fatal("Encode test failed with unexpected result.")
	}
	decoded, _ := crypter.Decrypt(encoded)
	if string(decoded) != input {
		t.Fatal("Decode test failed with unexpected result.")
	}
}
