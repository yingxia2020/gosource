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
	requestHandler := func(ctx *fasthttp.RequestCtx) {
		switch string(ctx.Path()) {
		case "/comp":
			compHandler(ctx)
		case "/decomp":
			decompHandler(ctx)
		default:
			otherHandler(ctx)
		}
	}

	s := &fasthttp.Server{
		Handler:     requestHandler,
		Concurrency: fasthttp.DefaultConcurrency,
	}

	if err := s.ListenAndServe(":8079"); err != nil {
		log.Fatalf("Error in ListenAndServe (de)compression server: %s", err)
	}
}

func otherHandler(ctx *fasthttp.RequestCtx) {
	out, err := exec.Command("./ex1", "--zlib-levels", "0", "--semidyn-config",
		"none", "--file", "mozilla").CombinedOutput()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Fprint(ctx, string(out))
}

func compHandler(ctx *fasthttp.RequestCtx) {
	out, err := exec.Command("./excomp", "--zlib-levels", "0", "--semidyn-config",
		"none", "--file", "ooffice").CombinedOutput()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Fprint(ctx, string(out))
}

func decompHandler(ctx *fasthttp.RequestCtx) {
	out, err := exec.Command("./exdecomp", "--zlib-levels", "0", "--semidyn-config",
		"none", "--file", "oofficedef").CombinedOutput()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Fprint(ctx, string(out))
}
