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
	"flag"
	"fmt"
	"github.com/valyala/fasthttp"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
	"time"
)

const (
	JOBSIZE    = 100
	WORKERSIZE = 20
)

var (
	addr    = flag.String("addr", ":8081", "TCP address to listen to")
	jobs    = make(chan int, JOBSIZE)
	results = make(chan string, JOBSIZE)
	jobID   = 0
)

func worker(id int, jobs <-chan int, results chan<- string) {
	for j := range jobs {
		fmt.Println("worker", id, "started  job", j)
		out, err := exec.Command("./hellogo").Output()
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("worker", id, "finished job", j)
		results <- fmt.Sprintf("job id %d, result %s", id, out)
	}
}

func main() {
	flag.Parse()

	s := &fasthttp.Server{
		Handler:     requestHandler,
		Concurrency: fasthttp.DefaultConcurrency,
	}

	var gracefulStop = make(chan os.Signal)
	signal.Notify(gracefulStop, syscall.SIGTERM)
	signal.Notify(gracefulStop, syscall.SIGINT)
	go func() {
		sig := <-gracefulStop
		close(jobs)
		close(results)
		fmt.Printf("caught sig: %+v\n", sig)
		fmt.Println("Wait for 1 second to finish processing")
		time.Sleep(1 * time.Second)
		os.Exit(0)
	}()

	// start workers first and run some jobs
	for w := 1; w <= WORKERSIZE; w++ {
		go worker(w, jobs, results)
		jobs <- jobID
		jobID++
	}

	if err := s.ListenAndServe(*addr); err != nil {
		log.Fatalf("Error in ListenAndServe: %s", err)
	}
}

func requestHandler(ctx *fasthttp.RequestCtx) {
	host, err := os.Hostname()
	if err != nil {
		panic(err)
	}
	fmt.Fprintf(ctx, "Hello from %v!\n\n", host)

	// Add one job to be used by later request
	jobs <- jobID
	jobID++

	fmt.Fprintf(ctx, "Run result %s!\n\n", <-results)
}
