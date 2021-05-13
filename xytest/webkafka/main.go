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
	"os/signal"
	"strings"
	"time"

	"github.com/Shopify/sarama"
	"github.com/valyala/fasthttp"
)

const (
	BROKERS  = "localhost:9092"
	TOPIC    = "dummy"
	CONTENTS = "Something cool"
)

// messages received by consumer
var messages []string

func main() {
	// Kafka producer creation:
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 5
	config.Producer.Return.Successes = true
	producer, err := sarama.NewSyncProducer([]string{BROKERS}, config)
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := producer.Close(); err != nil {
			panic(err)
		}
	}()

	requestHandler := func(ctx *fasthttp.RequestCtx) {
		switch string(ctx.Path()) {
		case "/send":
			sendHandler(ctx, producer)
		case "/summary":
			summaryHandler(ctx)
		default:
			otherHandler(ctx)
		}
	}

	s := &fasthttp.Server{
		Handler:     requestHandler,
		Concurrency: fasthttp.DefaultConcurrency,
	}

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)
	doneCh := make(chan struct{})

	// Kafka consumer part:
	go runConsumer(doneCh)

	go func() {
		if err := s.ListenAndServe(":8079"); err != nil {
			log.Fatalf("Error in ListenAndServe web kafka server: %s", err)
		}
	}()

	// wait until user click interrupt
	<-signals

	// let kafka consumer quit first
	doneCh <- struct{}{}

	// shutdown web server after 1 second
	time.Sleep(1 * time.Second)
	if err := s.Shutdown(); err != nil {
		log.Fatalf("Error in shut down web kafka server: %s", err)
	}
}

func runConsumer(doneCh chan struct{}) {
	done := false
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true
	master, err := sarama.NewConsumer([]string{BROKERS}, config)
	if err != nil {
		panic(err)
	}

	defer func() {
		done = true
		if err := master.Close(); err != nil {
			panic(err)
		}

		fmt.Println("Processed", len(messages), "messages:\n", strings.Join(messages, "\n"))
	}()

	consumer, err := master.ConsumePartition(TOPIC, 0, sarama.OffsetOldest)
	if err != nil {
		panic(err)
	}

	go func() {
		for {
			select {
			case err := <-consumer.Errors():
				if !done {
					fmt.Println(err)
				}
			case msg := <-consumer.Messages():
				messages = append(messages, string(msg.Value))
				fmt.Println("Received messages", string(msg.Key), string(msg.Value))
			}
		}
	}()
	<-doneCh
}

func summaryHandler(ctx *fasthttp.RequestCtx) {
	results := fmt.Sprintf("Processed %d messages:\n%s\n", len(messages), strings.Join(messages, "\n"))
	fmt.Fprint(ctx, results)
}

func otherHandler(ctx *fasthttp.RequestCtx) {
	fmt.Fprint(ctx, "Not a valid command\n")
}

func sendHandler(ctx *fasthttp.RequestCtx, producer sarama.SyncProducer) {
	content := ctx.QueryArgs().Peek("content")

	// Kafka producer part:
	var msg *sarama.ProducerMessage
	if len(content) == 0 {
		msg = &sarama.ProducerMessage{
			Topic: TOPIC,
			Value: sarama.StringEncoder(CONTENTS),
		}
	} else {
		msg = &sarama.ProducerMessage{
			Topic: TOPIC,
			Value: sarama.StringEncoder(content),
		}
	}
	partition, offset, err := producer.SendMessage(msg)
	if err != nil {
		log.Fatal(err)
	}

	out := fmt.Sprintf("Message is stored in topic(%s)/partition(%d)/offset(%d)\n", TOPIC, partition, offset)
	fmt.Fprint(ctx, string(out))
}
