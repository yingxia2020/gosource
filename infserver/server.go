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
	// s.DisableKeepalive = true

	if err := s.ListenAndServe(":8071"); err != nil {
		log.Fatalf("Error in ListenAndServe openvino server: %s", err)
	}
}

func handler(ctx *fasthttp.RequestCtx) {
	// infserver:v1.0
	// out, err := exec.Command("/root/app/openvino/script_run.sh", "-i", "/root/app/openvino/images_car/car_1.bmp").CombinedOutput()

	// infserver:v2.0
	// out, err := exec.Command("/root/app/openvino/test_run_sync_ImgClass_ResNET50-int8_1Batch.sh").CombinedOutput()

	// infserver:v3.0
	out, err := exec.Command("/root/app/openvino/test_run_sync_ObjDetec_SSDMobileNET-int8_1Batch.sh").CombinedOutput()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Fprint(ctx, string(out))
}
