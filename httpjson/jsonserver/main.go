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
	"encoding/json"
	"fmt"
	"net/http"
)

// The `json:"whatever"` bit is a way to tell the JSON
// encoder and decoder to use those names instead of the
// capitalised names
type person struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

var tom *person = &person{
	Name: "Tom",
	Age:  28,
}

func personHandler(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case "GET":
		// Just send out the JSON version of 'tom'
		j, _ := json.Marshal(tom)
		w.Write(j)
	case "POST":
		// Decode the JSON in the body and overwrite 'tom' with it
		d := json.NewDecoder(r.Body)
		p := &person{}
		err := d.Decode(p)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		tom = p
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprintf(w, "I can't do that.")
	}
}

func main() {
	http.HandleFunc("/person", personHandler)

	http.ListenAndServe(":8088", nil)
}
