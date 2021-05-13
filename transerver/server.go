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
	"regexp"
	"sync"

	"github.com/valyala/fasthttp"
)

var lock sync.Mutex

func main() {
	s := &fasthttp.Server{
		Handler:     handler,
		Concurrency: fasthttp.DefaultConcurrency,
	}

	if err := s.ListenAndServe(":8078"); err != nil {
		log.Fatalf("Error in ListenAndServe transcode server: %s", err)
	}
}

func handler(ctx *fasthttp.RequestCtx) {
	lock.Lock()
	defer lock.Unlock()
	out, err := exec.Command("./ffmpeg", "-y", "-i",
		"video_1920x1080_8bit_60Hz_P420.mp4", "-c:v", "libsvt_hevc", "-profile:v", "1",
		"-vf", "scale=240:180", "-f", "null", "-").CombinedOutput()
	if err != nil {
		log.Fatal(err)
	}

	re := regexp.MustCompile(` fps=(\d+) `)
	matches := re.FindStringSubmatch(string(out))

	// fmt.Fprint(ctx, string(out))
	fmt.Fprint(ctx, "FPS="+matches[1]+"\n")
}
