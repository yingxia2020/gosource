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
	"html/template"
	"log"
	"net/http"
)

type Result struct {
	Firstname string
	Lastname  string
	Birthday  string
	Title     string
}

type TodoPageData struct {
	PageTitle string
	Results   []Result
}

func main() {
	tmpl := template.Must(template.ParseFiles("layout.html"))

	http.Handle("/", http.FileServer(http.Dir("css/")))
	http.HandleFunc("/results", func(w http.ResponseWriter, r *http.Request) {
		data := TodoPageData{
			PageTitle: "CNB Test Results",
			Results: []Result{
				{Firstname: "Mary", Lastname: "Su", Birthday: "01/01/1988", Title: "Teacher"},
				{Firstname: "Candice", Lastname: "Su", Birthday: "01/01/1998", Title: "Student"},
				{Firstname: "Ying", Lastname: "Su", Birthday: "01/01/1990", Title: "Lawyer"},
			},
		}
		tmpl.Execute(w, data)
	})

	log.Fatal(http.ListenAndServe(":8088", nil))
}
