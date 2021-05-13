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
	"flag"
	"fmt"
	"log"
	"regexp"
	"strings"

	"github.com/valyala/fasthttp"
)

var (
	addr     = flag.String("addr", ":8088", "TCP address to listen to")
	compress = flag.Bool("compress", false, "Whether to enable transparent response compression")
)

func main() {
	flag.Parse()

	h := requestHandler
	if *compress {
		h = fasthttp.CompressHandler(h)
	}

	s := &fasthttp.Server{
		Handler:     h,
		Concurrency: fasthttp.DefaultConcurrency,
	}
	//	s.DisableKeepalive = false

	if err := s.ListenAndServe(*addr); err != nil {
		log.Fatalf("Error in ListenAndServe: %s", err)
	}
}

func requestHandler(ctx *fasthttp.RequestCtx) {
	var buffer bytes.Buffer

	temp := strings.TrimPrefix(ctx.Request.URI().String(), "http://localhost")
	re := regexp.MustCompile(`[/?&]`)
	result := re.Split(temp, -1)

	if len(result) >= 4 {
		if result[2] == "2" {
			ecommerceSim(&buffer, result[3], result)
		} else if result[2] == "3" {
			supportSim(&buffer, result[3], result)
		} /*else { // 1 is bank
			bankSim(&buffer, result[3], result)
		} */
	}
	fmt.Fprintf(ctx, buffer.String())
}
