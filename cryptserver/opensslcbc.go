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
	"github.com/spacemonkeygo/openssl"
)

const Cipher = "aes-256-cbc"

type Crypter struct {
	key    []byte
	iv     []byte
	cipher *openssl.Cipher
}

func NewCrypter(key []byte, iv []byte) (*Crypter, error) {
	cipher, err := openssl.GetCipherByName(Cipher)
	if err != nil {
		return nil, err
	}

	return &Crypter{key, iv, cipher}, nil
}

func (c *Crypter) Encrypt(input []byte) ([]byte, error) {
	ctx, err := openssl.NewEncryptionCipherCtx(c.cipher, nil, c.key, c.iv)
	if err != nil {
		return nil, err
	}

	cipherbytes, err := ctx.EncryptUpdate(input)
	if err != nil {
		return nil, err
	}

	finalbytes, err := ctx.EncryptFinal()
	if err != nil {
		return nil, err
	}

	cipherbytes = append(cipherbytes, finalbytes...)
	return cipherbytes, nil
}

func (c *Crypter) Decrypt(input []byte) ([]byte, error) {
	ctx, err := openssl.NewDecryptionCipherCtx(c.cipher, nil, c.key, c.iv)
	if err != nil {
		return nil, err
	}

	cipherbytes, err := ctx.DecryptUpdate(input)
	if err != nil {
		return nil, err
	}

	finalbytes, err := ctx.DecryptFinal()
	if err != nil {
		return nil, err
	}

	cipherbytes = append(cipherbytes, finalbytes...)
	return cipherbytes, nil
}
