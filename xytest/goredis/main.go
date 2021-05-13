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

	"github.com/intel/xytest/goredis/cache"
)

type Person struct {
	Name  string
	Phone string
}

const key = "XYZ"

func main() {

	cache.InitializeRedis()

	/*
		for i := 0; i < 10000; i++ {
			println("Setting Testkey -> TestValue")
			cache.SetValue("TestKey", "TestValue")

			println("Getting TestKey")
			value, err := cache.GetValue("TestKey")

			if err == nil {
				println("Value Returned : " + value.(string))
			} else {
				println("Getting Value Failed with error : " + err.Error())
			}
		}
	*/
	p1 := Person{"Ying", "123"}
	p2 := Person{"Mary", "456"}
	p3 := Person{"Tom", "789"}
	serializedValue, _ := json.Marshal(p1)
	cache.RPush(key, []string{string(serializedValue)})
	serializedValue, _ = json.Marshal(p2)
	serializedValue1, _ := json.Marshal(p3)
	cache.RPush(key, []string{string(serializedValue), string(serializedValue1)})

	results, err := cache.LRange(key)
	if err == nil {
		fmt.Println(len(results))
		// fmt.Println(results)
		for _, result := range results {
			p := Person{}
			err := json.Unmarshal([]byte(result), &p)
			if err == nil {
				fmt.Println(p)
			} else {
				fmt.Println(err.Error())
			}
		}
	}
}
