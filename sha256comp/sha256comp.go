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
	"crypto/sha256"
	"fmt"
	"io/ioutil"
	"log"
	"time"

	sha256simd "github.com/minio/sha256-simd"
)

const (
	Prefix = "./files/"
	Normal = false
)

var Files = []string{"world192.txt", "bible.txt", "E.coli"}

func main() {

	for _, file := range Files {
		start := time.Now()
		dataToEncrypt, err := ioutil.ReadFile(Prefix + file)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("loading time ", time.Since(start))

		if Normal {
			start := time.Now()
			sum := sha256.Sum256(dataToEncrypt)
			fmt.Println(time.Since(start))
			fmt.Printf("%x\n", sum)

		} else {
			start := time.Now()
			shaWriter := sha256simd.New()
			shaWriter.Write(dataToEncrypt)
			sumNew := shaWriter.Sum([]byte{})
			fmt.Println(time.Since(start))
			fmt.Printf("%x\n", sumNew)
		}
	}
}
