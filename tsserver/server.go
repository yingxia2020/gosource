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
	"os"
	"os/exec"

	"github.com/valyala/fasthttp"
)

func main() {
	s := &fasthttp.Server{
		Handler:     handler,
		Concurrency: fasthttp.DefaultConcurrency,
	}

	if err := s.ListenAndServe(":8078"); err != nil {
		log.Fatalf("Error in ListenAndServe tensor serving server: %s", err)
	}
}

func handler(ctx *fasthttp.RequestCtx) {
	serverAddr := os.Getenv("resnetAddr")
	var out []byte
	var err error
	// If not set, for docker environment use default value localhost:8500
	if len(serverAddr) == 0 {
		out, err = exec.Command("python",
			"/tmp/resnet/serving/tensorflow_serving/example/resnet_client_grpc.py").CombinedOutput()
	} else {
		// For Kubernete environment, the value for "resnetAddr" needs to be set
		out, err = exec.Command("python", "/tmp/resnet/serving/tensorflow_serving/example/resnet_client_grpc.py",
			"--server="+serverAddr+":8500").CombinedOutput()
	}

	if err != nil {
		fmt.Println(err.Error() + ": " + string(out))
		log.Fatal(err)
	} else {
		fmt.Fprint(ctx, string(out))
	}
}
