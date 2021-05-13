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
	"time"

	"github.com/otiai10/gosseract"
	"github.com/valyala/fasthttp"
)

func main() {
	s := &fasthttp.Server{
		Handler:     handler,
		Concurrency: fasthttp.DefaultConcurrency,
	}
	// s.DisableKeepalive = true

	if err := s.ListenAndServe(":8073"); err != nil {
		log.Fatalf("Error in ListenAndServe ocr server: %s", err)
	}
}

func handler(ctx *fasthttp.RequestCtx) {
	client := gosseract.NewClient()
	defer client.Close()
	fmt.Fprintf(ctx, "OCR version %s\n", gosseract.Version())

	startTime := time.Now()
	client.SetImage("./receipt.png")
	text, err := client.Text()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Fprint(ctx, string(text))
	fmt.Fprintf(ctx, "\n Time used %s\n", time.Since(startTime))
	/*
		startTime = time.Now()
		client.SetImage("./receipt1.png")
		text, err = client.Text()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Fprint(ctx, string(text))
		fmt.Fprintf(ctx, "\n Time used %s\n", time.Since(startTime))
	*/
}
