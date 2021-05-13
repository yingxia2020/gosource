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
	"io/ioutil"
	"log"
	"os"
	"time"
)

var (
	filesize  int64
	filelife  int
	filedir   string
	tryremove bool
)

func init() {
	flag.StringVar(&filedir, "d", "output", "File directory to be processed")
	flag.Int64Var(&filesize, "s", 0, "Minimum file size to keep")
	flag.IntVar(&filelife, "l", 1000, "File older than specified days will be removed")
	flag.BoolVar(&tryremove, "t", true, "Only show what files to be removed")
}

func main() {
	flag.Parse()

	var cutoff = time.Duration(filelife*24) * time.Hour
	var dir = "./" + filedir + "/"

	fileInfo, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Fatal(err.Error())
	}
	now := time.Now()
	for _, info := range fileInfo {
		if diff := now.Sub(info.ModTime()); diff > cutoff {
			fmt.Printf("Deleting %s which is %d days old\n", info.Name(), int(diff.Hours()/24))
			if !tryremove {
				err := os.Remove(dir + info.Name())
				if err != nil {
					log.Fatal(err.Error())
				}
			}
		}
		if info.Size() < filesize {
			fmt.Printf("Deleting %s which has size of %d \n", info.Name(), info.Size())
			if !tryremove {
				err := os.Remove(dir + info.Name())
				if err != nil {
					log.Fatal(err.Error())
				}
			}
		}
	}
}
