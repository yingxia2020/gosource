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

	"github.com/valyala/fasthttp"
)

func main() {
	requestHandler := func(ctx *fasthttp.RequestCtx) {
		switch string(ctx.Path()) {
		case "/upload":
			uploadHandler(ctx)
		default:
			handler(ctx)
		}
	}

	s := &fasthttp.Server{
		Handler:            requestHandler,
		MaxRequestBodySize: 10 * 1024 * 1024,
		Concurrency:        fasthttp.DefaultConcurrency,
	}

	if err := s.ListenAndServe(":8073"); err != nil {
		log.Fatalf("Error in ListenAndServe ocr server: %s", err.Error())
	}
}

func handler(ctx *fasthttp.RequestCtx) {
	fmt.Fprint(ctx, "nothing special")
}

func uploadHandler(ctx *fasthttp.RequestCtx) {
	header, err := ctx.FormFile("uploadfile")
	if err != nil {
		log.Fatalf("Failed to get file: %s", err.Error())
		return
	}
	err = fasthttp.SaveMultipartFile(header, fmt.Sprintf("./%s", header.Filename))
	if err != nil {
		log.Fatalf("Failed to save file: %s", err.Error())
		return
	}

	fmt.Fprintf(ctx, "File %s saved successfully", header.Filename)
}
