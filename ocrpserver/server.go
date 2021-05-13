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
	"os/exec"

	"github.com/valyala/fasthttp"
)

func main() {
	s := &fasthttp.Server{
		Handler:     handler,
		Concurrency: fasthttp.DefaultConcurrency,
	}

	if err := s.ListenAndServe(":8073"); err != nil {
		log.Fatalf("Error in ListenAndServe ocr python server: %s", err)
	}
}

func handler(ctx *fasthttp.RequestCtx) {
	out, err := exec.Command("python", "/root/app/tesseract-python/ocr_image_pdf.py").CombinedOutput()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Fprint(ctx, string(out))
}
