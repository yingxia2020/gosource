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
	"flag"
	"fmt"
	"log"
	"math"
	"math/rand"
	"net/http"

	//"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	addr = flag.String("addr", ":8070", "TCP address to listen to")
)

func main() {
	flag.Parse()

	http.HandleFunc("/", requestHandler)
	http.Handle("/metrics", promhttp.Handler())

	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Fatalf("Error in ListenAndServe: %s", err)
	}
}

func requestHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello from web light!\n\n")

	var x = 0.0001

	for i := 0; i <= 1000000+rand.Intn(100); i++ {
		x += math.Sqrt(x)
	}

	fmt.Fprintf(w, "result=%f OK\n\n", x)
}
