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
	"encoding/hex"
	"testing"
)

const (
	INPUT  = "hello world"
	OUTPUT = "656b50f2dbf795cd99ef704fd253a446b2a7ad2cc26e0ac705ac64"
)

var (
	tkey, _ = hex.DecodeString("1234567890ABCDEF1234567890ABCDEF")
	tiv, _  = hex.DecodeString("1234567890ABCDEF")
)

func TestGCMEncodeDecode(t *testing.T) {
	encoded, _ := EncryptString(INPUT, tkey, tiv)
	if hex.EncodeToString(encoded) != OUTPUT {
		t.Fatal("Encode test failed with unexpected result.")
	}
	decoded, _ := DecryptToString(encoded, tkey, tiv)
	if string(decoded) != INPUT {
		t.Fatal("Decode test failed with unexpected result.")
	}
}
