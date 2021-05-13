/* Copyright (C) Intel Corporation
 *
 * All Rights Reserved
 *
 * Unauthorized copying of this file, via any medium is strictly prohibited
 * Proprietary and confidential
 *
 * Written by Ying Xia <ying.xia@intel.com>, 2020
 */

package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gocql/gocql"
)

var session *gocql.Session

func main() {
	// connect cassandra
	var err error
	// connect to the cluster
	cluster := gocql.NewCluster("cassandra-0.cassandra", "cassandra-1.cassandra", "cassandra-2.cassandra")
	cluster.Keyspace = "example"
	cluster.Consistency = gocql.Quorum
	session, err = cluster.CreateSession()
	if err != nil {
		log.Fatal(err.Error())
	}
	defer session.Close()

	fmt.Println("Connected to cassandra DB")

	http.HandleFunc("/insert", insertHandler)
	http.HandleFunc("/query", queryHandler)
	http.HandleFunc("/delete", deleteHandler)

	http.ListenAndServe(":8088", nil)
}

func insertHandler(w http.ResponseWriter, r *http.Request) {
	// insert a tweet
	err := session.Query(`INSERT INTO tweet (timeline, id, text) VALUES (?, ?, ?)`,
		"me", gocql.TimeUUID(), "hello world").Exec()

	if err != nil {
		w.Write([]byte("error insert record into cassanra\n"))
	} else {
		w.Write([]byte("record inserted into cassandra successfully\n"))
	}
}

func queryHandler(w http.ResponseWriter, r *http.Request) {
	var id gocql.UUID
	var text string
	// list all tweets
	iter := session.Query(`SELECT id, text FROM tweet WHERE timeline = ?`, "me").Iter()
	for iter.Scan(&id, &text) {
		w.Write([]byte("Tweet: " + id.String() + " " + text + "\n"))
	}
	if err := iter.Close(); err != nil {
		w.Write([]byte("error close iter of cassanra\n"))
	}
}

func deleteHandler(w http.ResponseWriter, r *http.Request) {
	err := session.Query(`DELETE FROM tweet WHERE timeline = ?`, "me").Exec()
	if err != nil {
		w.Write([]byte("error delete record from cassanra\n"))
	} else {
		w.Write([]byte("record deleted from cassandra successfully\n"))
	}
}