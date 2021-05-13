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
	"log"
	"net/http"
	"regexp"
	"time"
)

func main() {
	http.HandleFunc("/", handler)

	http.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello World! I'm a HTTP server!")
	})

	//fs := http.FileServer(http.Dir("static/"))
	//http.Handle("/static/", http.StripPrefix("/static/", fs))

	log.Fatal(http.ListenAndServe(":8088", nil))
}

func doMakeResponse(w http.ResponseWriter, r *http.Request, done chan<- struct{}) {
	f, _ := w.(http.Flusher)
	fmt.Println(r.URL.String())
	re := regexp.MustCompile(`[/?&]`)
	result := re.Split(r.URL.String(), -1)

	var buffer bytes.Buffer

	// result[2] is model type, result[3] is action
	if len(result) >= 4 {
		if result[2] == "2" {
			ecommerceSim(&buffer, result[3], result)
		} /*else if result[2] == "1" {
			bankSim(&buffer, result[3], result)
		} else {
			supportSim(&buffer, result[3], result)
		} */
	}
	fmt.Fprintf(w, buffer.String())
	fmt.Println(len(buffer.String()))
	f.Flush()          // Send it to client.
	done <- struct{}{} // job finish
}

func handler(w http.ResponseWriter, r *http.Request) {
	done := make(chan struct{})
	/*
		buffer.WriteString("<pre>\n")
		buffer.WriteString("0\n")
		for i := minCurr; i <= maxCurr; i++ {
			buffer.WriteString(fmt.Sprintf("%s%02d\n", Currencies[rand.Intn(currSize)], i))
		}
		buffer.WriteString("</pre>\n")
	*/
	go doMakeResponse(w, r, done)
	select {
	case <-done:
		return
	case <-time.After(time.Second * 3):
		fmt.Fprint(w, "Server is busy.")
	}
}
