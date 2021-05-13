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
	"bufio"
	"bytes"
	"encoding/hex"
	"fmt"
	"io"
	//"io/ioutil"
	"os"
	"strings"

	"github.com/spacemonkeygo/openssl"
)

const (
	ChunkSize = 10 * 1024
	TagSize   = 16
)

func EncryptFile(inFile, outFile *os.File, key, iv []byte) error {
	ctx, err := openssl.NewGCMEncryptionCipherCtx(len(key)*8, nil, key, iv)
	if err != nil {
		return fmt.Errorf("Failed making GCM encryption ctx: %v", err)
	}

	reader := bufio.NewReader(inFile)
	chunk := make([]byte, ChunkSize)
	for {
		chunkSize, err := reader.Read(chunk)
		if err == io.EOF || chunkSize == 0 {
			break
		} else if err != nil {
			return fmt.Errorf("Failed to read a chunk: %v", err)
		}

		encData, err := ctx.EncryptUpdate(chunk[:chunkSize])
		if err != nil {
			return fmt.Errorf("Failed to perform an encryption: %v", err)
		}

		if _, err := outFile.Write(encData); err != nil {
			return fmt.Errorf("Failed to write an encrypted data: %v", err)
		}
	}

	encData, err := ctx.EncryptFinal()
	if err != nil {
		return fmt.Errorf("Failed to finalize encryption: %v", err)
	}
	if _, err := outFile.Write(encData); err != nil {
		return fmt.Errorf("Failed to write a final encrypted data: %v", err)
	}

	tag, err := ctx.GetTag()
	if err != nil {
		return fmt.Errorf("Failed to get GCM tag: %v", err)
	}
	if _, err := outFile.Write(tag); err != nil {
		return fmt.Errorf("Failed to write a gcm tag: %v", err)
	}

	return nil
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

func DecryptFile(inFile, outFile *os.File, key, iv []byte, fileSize int) error {
	ctx, err := openssl.NewGCMDecryptionCipherCtx(len(key)*8, nil, key, iv)
	if err != nil {
		return fmt.Errorf("Failed making GCM decryption ctx: %v", err)
	}

	reader := bufio.NewReader(inFile)
	chunk := make([]byte, ChunkSize)
	tag := &bytes.Buffer{}
	totalChunkSize := 0

	for {
		chunkSize, err := reader.Read(chunk)
		if err == io.EOF {
			break
		} else if err != nil {
			return fmt.Errorf("Failed to read an encrypted chunk: %v", err)
		}

		totalChunkSize += chunkSize

		if totalChunkSize > fileSize {
			d := totalChunkSize % fileSize

			d %= ChunkSize
			if d == 0 {
				d = chunkSize
			}

			tag.Write(chunk[chunkSize-d : chunkSize])
			chunkSize -= d
		}

		fmt.Println(totalChunkSize, chunkSize, fileSize)
		if chunkSize > 0 {
			data, err := ctx.DecryptUpdate(chunk[:chunkSize])
			if err != nil {
				return fmt.Errorf("Failed to perform a decryption: %v", err)
			}

			if _, err := outFile.Write(data); err != nil {
				return fmt.Errorf("Failed to write a decrypted data: %v", err)
			}
		}
	}

	if err := ctx.SetTag(tag.Bytes()); err != nil {
		return fmt.Errorf("Failed to set expected GCM tag: %v", err)
	}

	data, err := ctx.DecryptFinal()
	if err != nil {
		return fmt.Errorf("Failed to finalize decryption: %v", err)
	}
	if _, err := outFile.Write(data); err != nil {
		return fmt.Errorf("Failed to write a final decrypted data: %v", err)
	}

	return nil
}

func DecryptToString(input []byte, key, iv []byte) (string, error) {
	ctx, err := openssl.NewGCMDecryptionCipherCtx(len(key)*8, nil, key, iv)
	if err != nil {
		return "", fmt.Errorf("Failed making GCM decryption ctx: %v", err)
	}

	/*
		var buf bytes.Buffer
		writer := bufio.NewWriter(&buf)
		totalSize := len(input)
		var tag bytes.Buffer
		tag.Write(input[totalSize-TagSize:totalSize])
		data, err := ctx.DecryptUpdate(input[:totalSize-TagSize])
		if err != nil {
			return "", fmt.Errorf("Failed to perform a decryption: %v", err)
		}

		if _, err := writer.Write(data); err != nil {
			return "", fmt.Errorf("Failed to write a decrypted data: %v", err)
		}
	*/
	reader := bytes.NewReader(input)
	var buf bytes.Buffer
	writer := bufio.NewWriter(&buf)
	chunk := make([]byte, ChunkSize)
	tag := &bytes.Buffer{}
	totalChunkSize := 0
	inputSize := len(input) - TagSize

	for {
		chunkSize, err := reader.Read(chunk)
		if err == io.EOF {
			break
		} else if err != nil {
			return "", fmt.Errorf("Failed to read an encrypted chunk: %v", err)
		}

		totalChunkSize += chunkSize

		if totalChunkSize > inputSize {
			d := totalChunkSize % inputSize
			d %= ChunkSize
			if d == 0 {
				d = chunkSize
			}

			tag.Write(chunk[chunkSize-d : chunkSize])
			chunkSize -= d
		}

		if chunkSize > 0 {
			data, err := ctx.DecryptUpdate(chunk[:chunkSize])
			if err != nil {
				return "", fmt.Errorf("Failed to perform a decryption: %v", err)
			}

			if _, err := writer.Write(data); err != nil {
				return "", fmt.Errorf("Failed to write a decrypted data: %v", err)
			}
		}
	}

	if err := ctx.SetTag(tag.Bytes()); err != nil {
		return "", fmt.Errorf("Failed to set expected GCM tag: %v", err)
	}

	fdata, err := ctx.DecryptFinal()
	if err != nil {
		return "", fmt.Errorf("Failed to finalize decryption: %v", err)
	}
	if _, err := writer.Write(fdata); err != nil {
		return "", fmt.Errorf("Failed to write a final decrypted data: %v", err)
	}
	writer.Flush()

	return buf.String(), nil
}

func main() {
	/*
		key, _ := hex.DecodeString("1234567890ABCDEF1234567890ABCDEF")
		iv, _ := hex.DecodeString("1234567890ABCDEF")

		input := "hello worldhello worldhello worldhello worldhello worldhello worldhello worldhello worldhello worldhello worldhello worldhello worldhello worldhello worldhello worldhello worldhello worldhello world"
		output, err := EncryptString(input, key, iv)
		if err != nil {
			panic(err)
		}
		// result should be "656b50f2dbf795cd99ef704fd253a446b2a7ad2cc26e0ac705ac64"
		tmp := hex.EncodeToString(output)
		fmt.Println(tmp)

		tmp1, err := hex.DecodeString(tmp)
		if err != nil {
			panic(err)
		}
		decoded, err := DecryptToString(tmp1, key, iv)
		if err != nil {
			panic(err)
		}

		fmt.Println(decoded)

		test, err := ioutil.ReadFile("testdata7.enc")
		result, err := DecryptToString(test, key, iv)
		if err != nil {
			panic(err)
		}
		ioutil.WriteFile("xytest", []byte(result), 644)
	*/
	for i := 0; i <= 9; i++ {
		inFilePath := fmt.Sprintf("output%d.json", i)                  //"testdata"
		outEncFilePath := fmt.Sprintf("testdata%d.enc", i)             //"testdata.enc"
		outDecFilePath := fmt.Sprintf("newtestdata%d", i)              //"newtestdata"
		key, _ := hex.DecodeString("1234567890ABCDEF1234567890ABCDEF") //("81be5e09c111576103a8507658d47891")
		iv, _ := hex.DecodeString("1234567890ABCDEF")                  //("81be5e09c111576103a85076")

		inFile, err := os.Open(inFilePath)
		if err != nil {
			panic(err)
		}
		defer inFile.Close()

		fileStat, err := os.Stat(inFilePath)
		if err != nil {
			panic(err)
		}
		fileSize := int(fileStat.Size())

		outEncFile, err := os.Create(outEncFilePath)
		if err != nil {
			panic(err)
		}
		defer outEncFile.Close()

		outDecFile, err := os.Create(outDecFilePath)
		if err != nil {
			panic(err)
		}
		defer outDecFile.Close()

		if err := EncryptFile(inFile, outEncFile, key, iv); err != nil {
			panic(err)
		}
		outEncFile.Seek(0, 0)
		if err := DecryptFile(outEncFile, outDecFile, key, iv, fileSize); err != nil {
			panic(err)
		}
	}
}
