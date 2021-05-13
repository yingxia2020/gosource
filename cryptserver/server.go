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
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/intel/goredis"
)

const (
	Prefix  = "testdata"
	EncExt  = ".enc"
	Profile = "profile"
	DEBUG   = false
)

var (
	key, _ = hex.DecodeString("1234567890ABCDEF1234567890ABCDEF")
	iv, _  = hex.DecodeString("1234567890ABCDEF")
)

var client *goredis.MyRedis

// At initial stage, read encrypted files and put to redis
func init() {
	// connect to redis server
	client = goredis.InitializeRedis("redis-crypt-service.default.svc.cluster.local:6379")

	for i := 0; i <= 9; i++ {
		encdata, err := ioutil.ReadFile("./files/" + Prefix + fmt.Sprintf("%d", i) + EncExt)
		if err != nil {
			log.Fatal(err)
		}
		client.SetRawValue(Profile+fmt.Sprintf("%d", i), encdata)
	}
}

func main() {
	http.HandleFunc("/", handler)
	http.ListenAndServe(":8076", nil)
}

func handler(w http.ResponseWriter, r *http.Request) {

	// For the case encrypt a short string such as password
	einput, ok := r.URL.Query()["encinput"]
	if ok && len(einput[0]) > 1 {
		output, err := EncryptString(einput[0], key, iv)
		if err != nil {
			log.Fatal(err)
		}
		w.Write([]byte(hex.EncodeToString(output)))
		return
	}

	// For the case decrypt a short string such as password
	dinput, ok := r.URL.Query()["decinput"]
	if ok && len(dinput[0]) > 1 {
		temp, err := hex.DecodeString(dinput[0])
		if err != nil {
			log.Fatal(err)
		}
		output, err := DecryptToString(temp, key, iv)
		if err != nil {
			log.Fatal(err)
		}
		w.Write([]byte(output))
		return
	}

	// For the case decrypt user profile
	var name string
	names, ok := r.URL.Query()["name"]

	if !ok || len(names[0]) < 1 {
		name = generateKey()
	} else {
		name = names[0]
	}

	rand.Seed(time.Now().UTC().UnixNano())
	idx := rand.Intn(10)
	if DEBUG {
		fmt.Println("Profile index = ", idx)
	}

	// Random pick up a profile to decrypt
	encdata, err := client.GetRawValue(Profile + fmt.Sprintf("%d", idx))
	if err != nil {
		log.Fatal(err)
	}
	if len(encdata) == 0 {
		w.Write([]byte("Could not find user profile data for " + name + "\n"))
		return
	}

	decdata, err := DecryptToString([]byte(encdata), key, iv)
	if err != nil {
		log.Fatal(err)
	}
	if DEBUG {
		fmt.Println(name)
	}
	client.SetRawValue(name, []byte(decdata))

	w.Write([]byte(name + " profile decryption finished successfully!\n"))
}

func generateKey() string {
	prefix := "user_"
	id := rand.Intn(1000)

	return fmt.Sprintf("%s%d", prefix, id)
}
