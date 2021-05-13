/* Copyright (C) XY Corporation
 *
 * All Rights Reserved
 *
 * Unauthorized copying of this file, via any medium is strictly prohibited
 * Proprietary and confidential
 *
 * Written by Ying Xia, 2021
 */

package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/valyala/fasthttp"
)

type name struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

type value struct {
	Categories []string `json:"categories"`
	ID         int      `json:"id"`
	Joke       string   `json:"joke"`
}

type joke struct {
	Type  string `json:"type"`
	Value value  `json:"value"`
}

const (
	NameURL = "https://names.mcquay.me/api/v0/"
	jobNum  = 5
)

var (
	addr     = flag.String("addr", ":8070", "TCP address to listen to")
	addrTLS  = flag.String("addrTLS", ":8443", "TCP address to listen to TLS (aka SSL or HTTPS) requests. Leave empty for disabling TLS")
	certFile = flag.String("certFile", "./xyserver.crt", "Path to TLS certificate file")
	keyFile  = flag.String("keyFile", "./xyserver.key", "Path to TLS key file")
)

func main() {
	// Parse command-line flags.
	flag.Parse()

	requestHandler := func(ctx *fasthttp.RequestCtx) {
		path := string(ctx.Path())
		jobs := make(chan struct{}, jobNum)
		jobresults := make(chan string, jobNum)
		defer close(jobs)
		for i := 1; i < jobNum; i++ {
			go jokeworker(jobs, jobresults)
		}
		switch path {
		// We could only handle joke case, the other cases considered TODO
		case "/joke":
			jokeHandler(ctx, jobs, jobresults)
		default:
			notSupportedHandler(ctx)
		}
	}

	s := &fasthttp.Server{
		Handler:     requestHandler,
		Concurrency: fasthttp.DefaultConcurrency,
	}

	// Start HTTP server.
	if len(*addr) > 0 {
		log.Printf("Starting Joke HTTP server on %q", *addr)
		go func() {
			if err := s.ListenAndServe(*addr); err != nil {
				log.Fatalf("error in Joke Server  ListenAndServe: %s", err)
			}
		}()
	}

	// Start HTTPS server.
	if len(*addrTLS) > 0 {
		log.Printf("Starting Joke HTTPS server on %q", *addrTLS)
		go func() {
			if err := fasthttp.ListenAndServeTLS(*addrTLS, *certFile, *keyFile, requestHandler); err != nil {
				log.Fatalf("error in Joke Server ListenAndServeTLS: %s", err)
			}
		}()
	}

	// Wait and serve.
	select {}
}

func notSupportedHandler(ctx *fasthttp.RequestCtx) {
	fmt.Fprintln(ctx, "This is a not supported request. Have a nice day")
}

func getJSON(url string, target interface{}) error {
	// time out after 3 seconds
	var myClient = &http.Client{Timeout: 3 * time.Second}

	resp, err := myClient.Get(url)
	if resp != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		return err
	}

	return json.NewDecoder(resp.Body).Decode(target)
}

func jokeworker(jobs <-chan struct{}, jobresults chan<- string) {
	for _ = range jobs {
		var results bytes.Buffer
		rname := name{}
		err := getJSON(NameURL, &rname)
		if err != nil {
			results.WriteString(fmt.Sprintf("Error when acquire name from name server: %s\n", err.Error()))
			jobresults <- results.String()
			continue
		}
		results.WriteString(fmt.Sprintf("First name %s, Last name %s\n", rname.FirstName, rname.LastName))

		jokeURL := fmt.Sprintf("http://api.icndb.com/jokes/random?firstName=%s&lastName=%s&limitTo=[nerdy]",
			rname.FirstName, rname.LastName)

		myjoke := joke{}
		err = getJSON(jokeURL, &myjoke)
		if err != nil || myjoke.Type != "success" {
			results.WriteString(fmt.Sprintf("Error when acquire joke from joke server: %s\n", err.Error()))
			jobresults <- results.String()
			continue
		}
		results.WriteString(fmt.Sprintf("Here is your joke for today: %s\n", myjoke.Value.Joke))
		jobresults <- results.String()
	}
}

func jokeHandler(ctx *fasthttp.RequestCtx, jobs chan<- struct{}, jobresults <-chan string) {
	jobs <- struct{}{}

	results := <-jobresults
	fmt.Fprintf(ctx, results)
}
