/*******************************************************************************
* Copyright 2021 Intel Corporation
*
* Licensed under the Apache License, Version 2.0 (the "License");
* you may not use this file except in compliance with the License.
* You may obtain a copy of the License at
*
*     http://www.apache.org/licenses/LICENSE-2.0
*
* Unless required by applicable law or agreed to in writing, software
* distributed under the License is distributed on an "AS IS" BASIS,
* WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
* See the License for the specific language governing permissions and
* limitations under the License.
*******************************************************************************/

package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"

	"github.com/valyala/fasthttp"
)

var buffer_size = os.Getenv("GCM_BUFFER_SIZE")

func main() {
	s := &fasthttp.Server{
		Handler:     handler,
		Concurrency: fasthttp.DefaultConcurrency,
	}

	if err := s.ListenAndServe(":8074"); err != nil {
		log.Fatalf("Error in ListenAndServe ocr server: %s", err)
	}
}

func handler(ctx *fasthttp.RequestCtx) {

	tmp, err := strconv.Atoi(buffer_size)
	if len(buffer_size) == 0 || err != nil || tmp < 1 {
		log.Println("buffer size used is 262144")
		buffer_size = "262144"
	}

	out, err := exec.Command("./gcmbench", "-b", buffer_size).CombinedOutput()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Fprint(ctx, string(out))

	// return POD info
	podName, _ := exec.Command("printenv", "MY_POD_NAME").CombinedOutput()
	fmt.Fprint(ctx, "\nFrom POD: "+string(podName))
}
