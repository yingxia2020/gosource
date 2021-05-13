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
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"time"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

type Record struct {
	Stock  float64 `json:"stock"`
	Strike float64 `json:"strike"`
	Year   float64 `json:"year"`
}

const LEN = 4096

func main() {
	var mc [LEN]Record

	for j := 0; j < 10; j++ {
		rand.Seed(time.Now().UTC().UnixNano())
		for i := 0; i < LEN; i++ {
			v1 := rand.Float64()*45 + 5
			v2 := rand.Float64()*15 + 10
			v3 := rand.Float64()*4 + 1
			r := Record{Stock: v1, Strike: v2, Year: v3}
			mc[i] = r
		}

		output, _ := json.MarshalIndent(mc, "", "    ")
		err := ioutil.WriteFile(fmt.Sprintf("output%d.json", j), output, 0644)
		check(err)
	}
}
