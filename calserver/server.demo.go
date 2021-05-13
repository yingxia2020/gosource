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
	"log"
	"math"

	"github.com/valyala/fasthttp"
)

func main() {
	s := &fasthttp.Server{
		Handler:     handler,
		Concurrency: fasthttp.DefaultConcurrency,
	}

	if err := s.ListenAndServe(":8072"); err != nil {
		log.Fatalf("Error in ListenAndServe calc server: %s", err)
	}
}

/* Use 5000000 times original SQRT calculations to scale up on powerful machines */
func handler(ctx *fasthttp.RequestCtx) {
	var x = 0.0001

	for i := 0; i <= 5000000; i++ {
		x += math.Sqrt(x)
	}

	fmt.Fprintf(ctx, "result=%f OK\n\n", x)
}
