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
	"gocv.io/x/gocv"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", handler)

	if err := http.ListenAndServe(":8075", nil); err != nil {
		log.Fatalf("Error in ListenAndServe face detection server: %s", err)
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	img := gocv.IMRead("/home/n/temp/face.jpg", gocv.IMReadColor)
	if img.Empty() {
		fmt.Fprint(w, "Invalid Mat for face detection")
		return
	}
	defer img.Close()

	// load HOGDescriptor to recognize people
	hog := gocv.NewHOGDescriptor()
	defer hog.Close()

	d := gocv.HOGDefaultPeopleDetector()
	defer d.Close()
	hog.SetSVMDetector(d)

	rects := hog.DetectMultiScale(img)
	if len(rects) == 0 {
		fmt.Fprintln(w, "Error in TestCascadeClassifier test")
	} else {
		fmt.Fprintf(w, "%d faces are detected\n", len(rects))
	}
}
