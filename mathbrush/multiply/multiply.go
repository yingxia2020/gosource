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
	"fmt"
	"io/ioutil"
	"math/rand"
	"time"
)

var (
	x = 230
	y = 4
)

func main() {
	var buf, buf1 bytes.Buffer
	var a, b int
	rand.Seed(time.Now().UnixNano())
	for j := 0; j < x; j++ {
		for i := 0; i < y; i++ {
			a = rand.Intn(89) + 10
			b = rand.Intn(89) + 10
			buf.WriteString(fmt.Sprintf("%d X %d =           ", a, b))
			buf1.WriteString(fmt.Sprintf("%d X %d = %-10d", a, b, a*b))
		}
		buf.WriteString("\n\n")
		buf1.WriteString("\n\n")
	}

	//fmt.Println(buf.String())
	//fmt.Println("\n")
	//fmt.Println(buf1.String())

	ioutil.WriteFile("quest.txt", buf.Bytes(), 0644)
	ioutil.WriteFile("answer.txt", buf1.Bytes(), 0644)
}
