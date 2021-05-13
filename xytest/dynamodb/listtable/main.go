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
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

func main() {
	start := time.Now()
	// Initialize a session in us-west-2 that the SDK will use to load
	// credentials from the shared credentials file ~/.aws/credentials.
	sess, err := session.NewSession(&aws.Config{
		Region:   aws.String("us-west-2"),
		Endpoint: aws.String("http://localhost:8000"),
	},
	)

	// Create DynamoDB client
	svc := dynamodb.New(sess)
	fmt.Println(time.Since(start))
	start = time.Now()
	// Get the list of tables
	result, err := svc.ListTables(&dynamodb.ListTablesInput{})

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println(time.Since(start))
	start = time.Now()

	fmt.Println("Tables:")

	for _, n := range result.TableNames {
		fmt.Println(*n)
	}

	fmt.Println("")
	fmt.Println(time.Since(start))
}
