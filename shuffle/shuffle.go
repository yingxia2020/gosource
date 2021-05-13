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
	"fmt"
	"math/rand"
)

const SIZE = 50

func shuffle(urls []string) {
	for i := 0; i < SIZE; i++ {
		temp := rand.Intn(SIZE)
		urls[i], urls[temp] = urls[temp], urls[i]
	}
}

func main() {
	var urls = make([]string, SIZE)

	for i := 0; i < SIZE; i++ {
		if i < 30 {
			urls[i] = "aaaaaa"
		} else {
			urls[i] = "bbbbbb"
		}
	}
	fmt.Println("Before shuffle:", urls)
	shuffle(urls)
	fmt.Println("After shuffle:", urls)
}
