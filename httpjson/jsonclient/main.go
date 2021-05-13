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
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type person struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func main() {
	url := "http://localhost:8088/person"
	fmt.Println("URL:>", url)

	client := &http.Client{}
	var err error
	req, _ := http.NewRequest("GET", url, nil)

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("response Status:", resp.Status)
	fmt.Println("response Headers:", resp.Header)
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("response Body:", string(body))
	resp.Body.Close()

	var ying *person = &person{
		Name: "Ying",
		Age:  18,
	}

	jsonStr, _ := json.Marshal(ying)
	req, _ = http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	resp, err = client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("response Status:", resp.Status)
	fmt.Println("response Headers:", resp.Header)
	body, _ = ioutil.ReadAll(resp.Body)
	fmt.Println("response Body:", string(body))
	resp.Body.Close()

	req, _ = http.NewRequest("GET", url, nil)
	resp, err = client.Do(req)
	if err != nil {
		panic(err)
	}

	fmt.Println("response Status:", resp.Status)
	fmt.Println("response Headers:", resp.Header)
	body, _ = ioutil.ReadAll(resp.Body)
	fmt.Println("response Body:", string(body))
	resp.Body.Close()
}
